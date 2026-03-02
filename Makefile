.PHONY: generate test lint check all

generate: ## Run templ generate inside Docker
	docker compose run --rm dev templ generate

test: ## Run go test
	docker compose run --rm dev go test ./...

lint: ## Run go vet
	docker compose run --rm dev go vet ./...

check: ## Verify generated files are fresh (templ generate && git diff --exit-code)
	docker compose run --rm dev sh -c "templ generate && git diff --exit-code"

all: generate test lint ## generate + test + lint
