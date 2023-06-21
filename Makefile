WORKDIR := $(PWD)

LINTER_CONFIG_LOCATION=.golangci.yml
REPORT_ARTIFACTS=./reports
BOOTSTRAP_PROJECT_DIR=son

export UNIT_TEST_PRECOMPILE=true

# If the first argument is "migrate-new"...
ifneq (,$(findstring migrate-new,$(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  MIGRATIONNAME := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(MIGRATIONNAME):;@:)
endif

generate:
	@echo "> Generating mock files..."
	go generate -x ./...
.PHONY: generate

local-lint:
	@echo "> Running lint..."
	golangci-lint run --config=.golangci.yml
.PHONY: local-lint

local-test:
	@echo "> Running tests..."
	go test -v -race -cover -coverprofile "./unit.cov" -covermode=atomic ./...
.PHONY: local-test

local-test-cover: local-test
	@echo "> Running test cover..."
	go tool cover -html="./unit.cov"
.PHONY: local-test-cover

local-build:
	@echo "> Building binary..."
	go build -v --race -o $(WORKDIR)/bin/ $(WORKDIR)/cmd/${BOOTSTRAP_PROJECT_DIR}
.PHONY: local-build

local-run:
	@echo "> Building binary..."
	go build -v --race -o $(WORKDIR)/bin/ $(WORKDIR)/cmd/${BOOTSTRAP_PROJECT_DIR}
	$(WORKDIR)/bin/${BOOTSTRAP_PROJECT_DIR}
.PHONY: local-run

migrate-new:
	@echo "> Creating a new migrate file..."
# use https://github.com/golang-migrate/migrate library
	migrate create -ext sql -dir ./migrations $(MIGRATIONNAME)
.PHONY: new-migrate

download-deps:
	@echo "> Download dependencies"
	go mod tidy && go mod vendor
.PHONY: download-deps

compose-up:
	docker-compose -f docker-compose.yml up -d

compose-down:
	docker-compose -f docker-compose.yml down -v

compose-restart:
	docker-compose -f docker-compose.yml down -v
	docker-compose -f docker-compose.yml up -d

compose-rebuild:
	docker-compose -f docker-compose.yml down -v
	docker-compose -f docker-compose.yml up --build -d

test-hw-2:
	wrk -t1 -c1 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
	wrk -t10 -c10 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'
	wrk -t100 -c100 -d30s --timeout 30s -s 'resources/wrk/search.lua' 'http://localhost:8086'


