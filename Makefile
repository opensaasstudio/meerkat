init:
	mkdir -p ./gitignored/deployconfig
	cp -n -v ./deploytool/**/*.yaml ./gitignored/deployconfig/ || true

.PHONY: deploy_slackapigatewayhandlerbylambda
deploy_slackapigatewayhandlerbylambda:
	GOOS=linux GOARCH=amd64 go build -o build/meerkat-slackapigatewayhandlerbylambda ./adapter/slack/cmd/meerkat-slackapigatewayhandlerbylambda
	cat build/meerkat-slackapigatewayhandlerbylambda | go run ./deploytool/lambda/deploy.go --configFile ./gitignored/deployconfig/meerkat-slackapigatewayhandlerbylambda.yaml
