.PHONY: deps build test help
APP=ivwcli

any: help

deps: ## install dependencies
	@go mod tidy

build: deps ## build the Go binary
	go build -o $(APP)

test: deps ## run tests (integration+unit)
	go test --tags=integration ./... -cover

docker-build: ## build the docker image
	docker build -t $(APP) .

MATCH_RECIPES ?= "Potato,Veggie,Mushroom"
POSTCODE ?= "10120"
DELIVERY_FROM ?= "10AM"
DELIVERY_TO ?= "3PM"
LOGGING ?= 0
ifeq ($(LOGGING),1)
	LOG_FLAG="-l"
endif

run-docker: docker-build ## run the application via the docker container
	@if [ -z "$(FILE)" ]; then echo "missing file param"; exit 1; fi

	docker run --rm --name $(APP) \
			--mount type=bind,source="$(FILE)",target=/input_file.json  $(APP) \
			--file /input_file.json --match-recipes $(MATCH_RECIPES) -p $(POSTCODE) --from $(DELIVERY_FROM) --to $(DELIVERY_TO) $(LOG_FLAG)

help:
	@echo -- Usage --
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


