#!make

.EXPORT_ALL_VARIABLES:
PARAM=$(filter-out $@,$(MAKECMDGOALS))

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

%:
	@:

########################################################################################################################

test: ## test package inside docker container
	@clear
	docker build --force-rm --tag go-workers --file Dockerfile .
	docker run --rm --interactive --workdir "/src" --volume "$$(pwd)":"/src":rw go-workers sh -c "make test"
	@$(MAKE) owner

owner:  ## reset project owner
	sudo chown --changes -R $$(whoami) ./
