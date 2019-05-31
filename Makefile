.PHONY: init
init: go_dep
	mkdir -p ./gitignored/deployconfig
	cp -n -v ./deploytool/**/*.yaml ./gitignored/deployconfig/ || true

.PHONY: install_tools
install_tools:
	go get github.com/golang/dep/cmd/dep

.PHONY: go_dep
go_dep:
	dep ensure

.PHONY: deploy
deploy: deploy_slackapigatewayhandlerbylambda deploy_askallifneeded

.PHONY: deploy_slackapigatewayhandlerbylambda
deploy_slackapigatewayhandlerbylambda:
	GOOS=linux GOARCH=amd64 go build -o build/meerkat-slackapigatewayhandlerbylambda ./adapter/slack/cmd/meerkat-slackapigatewayhandlerbylambda
	cat build/meerkat-slackapigatewayhandlerbylambda | go run ./deploytool/lambda/deploy.go --configFile ./gitignored/deployconfig/meerkat-slackapigatewayhandlerbylambda.yaml

.PHONY: deploy_askallifneeded
deploy_askallifneeded:
	GOOS=linux GOARCH=amd64 go build -o build/meerkat-askallifneeded ./adapter/cmd/ask_all_if_needed
	cat build/meerkat-askallifneeded | go run ./deploytool/lambda/deploy.go --configFile ./gitignored/deployconfig/meerkat-askallifneeded.yaml

.PHONY:
test:
	go test -v ./...
