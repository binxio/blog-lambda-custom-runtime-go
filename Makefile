.PHONY: help init clean
.DEFAULT_GOAL := help
environment = "example"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

create: dist ## create env
	@sceptre launch-env $(environment)

delete: ## delete env
	@sceptre delete-env $(environment)

deploy: delete create info ## deploy the example

info: ## describe resources
	@sceptre describe-stack-outputs $(environment) lambda

build: ## build the lambda
	$(shell GOOS=linux go build -o bootstrap bootstrap.go)
	file bootstrap

dist: clean build ## build the lambda
	zip lambda.zip bootstrap

clean: ## remove artifacts
	-rm bootstrap lambda.zip

artillery: ## load test
	artillery quick --count 100 -n 100 $(shell sceptre --output json describe-stack-outputs example lambda | jq -r ' .[] | select(.OutputKey=="CustomServiceEndpoint") | .OutputValue')