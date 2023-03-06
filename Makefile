rsa: ## Generate RSA private and public keys if it not exist
ifeq (,$(wildcard ./private.pem))
	@printf "\033[36m%s\033[0m\n" "Generating RSA keys..."
	openssl genrsa -out private.pem 2048
	openssl rsa -in private.pem -outform PEM -pubout -out public.pem
endif

lint: ## Run linter
	golangci-lint run

test: ## Run tests
	go test -v ./...
