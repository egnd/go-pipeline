#!make

MAKEFLAGS += --always-make

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

%:
	@:

########################################################################################################################

owner: ## Reset folder owner
	sudo chown --changes -R $$(whoami) ./
	@echo "Success"

check-conflicts: ## Find git conflicts
	@if grep -rn '^<<<\<<<< ' .; then exit 1; fi
	@if grep -rn '^===\====$$' .; then exit 1; fi
	@if grep -rn '^>>>\>>>> ' .; then exit 1; fi
	@echo "All is OK"

check-todos: ## Find TODO's
	@if grep -rn '@TO\DO:' .; then exit 1; fi
	@echo "All is OK"

check-master: ## Check for latest master in current branch
	@git remote update
	@if ! git log --pretty=format:'%H' | grep $$(git log --pretty=format:'%H' -n 1 origin/master) > /dev/null; then exit 1; fi
	@echo "All is OK"

mocks: ## Generate mocks
	@clear && rm -rf mocks && mkdir -p mocks
	mockery --all --case=underscore --recursive --outpkg=mocks --output=mocks --dir=./pkg

tests: ## Run unit tests
	@clear && rm -rf coverage && mkdir -p coverage
	CGO_ENABLED=1 go test -mod=readonly -race -cover -covermode=atomic -coverprofile=coverage/profile.out ./pkg

benchmarks: ## Run benchmarks
	@clear
	go test -mod=readonly -benchmem -bench . ./pkg/benchmarks_test.go

coverage: tests ## Check code coveragem
	go tool cover -func=coverage/profile.out
	go tool cover -html=coverage/profile.out -o coverage/report.html

profiling: ## Run unit tests
	@clear && rm -rf coverage && mkdir -p coverage
	go test -mod=readonly -cpuprofile=coverage/cpu.prof -memprofile=coverage/mem.prof ./pkg
	go tool pprof -svg coverage/cpu.prof > coverage/cpu.svg
	go tool pprof -svg coverage/mem.prof > coverage/mem.svg

lint: ## Lint source code
	@clear
	golangci-lint run --color=always --config=.golangci.yml ./pkg

########################################################################################################################

docker-lint:
	docker run --rm -it -v $$(pwd):/src -w /src --entrypoint make golangci/golangci-lint:v1.41 lint

docker-mocks:
	docker run --rm -it -v $$(pwd):/src -w /src --entrypoint sh vektra/mockery:v2 -c "apk add -q make && make mocks"

docker-tests:
	docker run --rm -it -v $$(pwd):/src -w /src --entrypoint make golang:1.16 tests

docker-coverage:
	docker run --rm -it -v $$(pwd):/src -w /src --env-file=.env --entrypoint make golang:1.16 coverage
	@echo "Read report at file://$$(pwd)/coverage/report.html"

docker-benchmarks:
	docker run --rm -it -v $$(pwd):/src -w /src --entrypoint make golang:1.16 benchmarks

docker-profiling:
	docker run --rm -it -v $$(pwd):/src -w /src --entrypoint sh golang:1.16 -c "apt-get update && apt-get install -y graphviz && make profiling"
	@echo "Look at file://$$(pwd)/coverage/cpu.svg"
	@echo "Look at file://$$(pwd)/coverage/mem.svg"

docker-vendor:
	docker run --rm -it -v $$(pwd):/src -w /src --env-file=.env --entrypoint go golang:1.16 mod vendor
