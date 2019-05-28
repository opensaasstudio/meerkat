package main

import (
	"archive/zip"
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hori-ryota/zaperr"
	"gopkg.in/yaml.v2"
)

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}

func Main() error {
	ctx := context.Background()
	config := Config{}

	configFilePath := flag.String("configFile", "", "function config file")
	binFilePath := flag.String("binFile", "", "binFile path")
	flag.Parse()

	if configFilePath == nil || *configFilePath == "" {
		return zaperr.New("configFile flag is required")
	}
	configText, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(configText, &config); err != nil {
		return err
	}

	binFile, err := readBinFile(binFilePath)
	if err != nil {
		return err
	}
	zipFile := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipFile)
	h := &zip.FileHeader{
		Name:   config.Name,
		Method: zip.Deflate,

		Modified: time.Now(),
	}
	h.SetMode(0755)

	w, err := zipWriter.CreateHeader(h)
	if err != nil {
		return err
	}
	if _, err := w.Write(binFile); err != nil {
		return err
	}
	zipWriter.Close()

	lambdaClient := lambda.New(session.New(&aws.Config{
		Region: aws.String(config.Region),
	}))

	exists, err := existsFunction(ctx, config.Name, lambdaClient)
	if err != nil {
		return err
	}

	if !exists {
		// create
		lambdaInput := lambda.CreateFunctionInput{
			FunctionName: aws.String(config.Name),
			Handler:      aws.String(config.Name),
			Role:         aws.String(config.Role),
			Code: &lambda.FunctionCode{
				ZipFile: zipFile.Bytes(),
			},
			MemorySize: aws.Int64(int64(config.Memory)),
			Timeout:    aws.Int64(int64(config.Timeout)),

			Runtime: aws.String("go1.x"),
		}

		lambdaInput.Environment = &lambda.Environment{
			Variables: make(map[string]*string, len(config.ENV)),
		}
		for k, v := range config.ENV {
			lambdaInput.Environment.Variables[k] = aws.String(v)
		}

		lambdaOutput, err := lambdaClient.CreateFunctionWithContext(ctx, &lambdaInput)
		if err != nil {
			return err
		}
		log.Printf("%+v\n", lambdaOutput)
		return nil
	}

	// update
	lambdaConfigurationInput := lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(config.Name),
		Handler:      aws.String(config.Name),
		Role:         aws.String(config.Role),
		MemorySize:   aws.Int64(int64(config.Memory)),
		Timeout:      aws.Int64(int64(config.Timeout)),

		Runtime: aws.String("go1.x"),
	}

	lambdaConfigurationInput.Environment = &lambda.Environment{
		Variables: make(map[string]*string, len(config.ENV)),
	}
	for k, v := range config.ENV {
		lambdaConfigurationInput.Environment.Variables[k] = aws.String(v)
	}
	lambdaConfigurationOutput, err := lambdaClient.UpdateFunctionConfigurationWithContext(ctx, &lambdaConfigurationInput)
	if err != nil {
		return err
	}
	log.Printf("%+v\n", lambdaConfigurationOutput)

	lambdaCodeInput := lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(config.Name),
		ZipFile:      zipFile.Bytes(),
	}
	lambdaCodeOutput, err := lambdaClient.UpdateFunctionCodeWithContext(ctx, &lambdaCodeInput)
	if err != nil {
		return err
	}
	log.Printf("%+v\n", lambdaCodeOutput)

	return nil
}

type Config struct {
	Name    string
	Region  string
	Role    string
	Memory  uint32
	Timeout uint32

	ENV map[string]string
}

func readBinFile(filePath *string) ([]byte, error) {
	if filePath == nil || *filePath == "" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(*filePath)
}

func existsFunction(ctx context.Context, functionName string, lambdaClient *lambda.Lambda) (bool, error) {
	_, err := lambdaClient.GetFunctionWithContext(ctx, &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			if err.Code() == lambda.ErrCodeResourceNotFoundException {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
