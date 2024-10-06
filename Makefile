include deploy/.env
export

# Change these variables as necessary.
main_package_path = ./cmd/
binary_name = follow_service
temp := $(CURDIR)/tmp
temp_bin := ${temp}/bin
project_pkg = github.com/Karzoug/meower-post-service
api_file_name=post

version ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
commit_hash ?= $(shell git rev-parse --short HEAD 2>/dev/null)
build_date ?= $(shell date +%FT%T%z)

# remove debug info from the binary & make it smaller
ldflags += -s -w
ldflags += -X ${project_pkg}/pkg/buildinfo.version=${version} -X ${project_pkg}/pkg/buildinfo.commitHash=${commit_hash} -X ${project_pkg}/pkg/buildinfo.buildDate=${build_date}

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test fmt lint
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

fmt:
	go run golang.org/x/tools/cmd/goimports@latest -local=${project_pkg} -l -w  .
	go run mvdan.cc/gofumpt@latest -l -w  .

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## lint: run linters
.PHONY: lint
lint:
	$(temp_bin)/golangci-lint run ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build:
	go build -ldflags "${ldflags}" -o ${temp_bin}/${binary_name} ${main_package_path}

## run-host: run the  application on host machine
.PHONY: run-host
run-host: build	
	${temp_bin}/${binary_name}

## run: run the  application in docker
.PHONY: run
run: 
	VERSION=${version} BUILD_DATE=${build_date} COMMIT_HASH=${commit_hash} docker compose -f "deploy/docker-compose.yaml" up -d --build

## stop: stop the  application in docker
.PHONY: stop
stop: 
	docker compose -f "deploy/docker-compose.yaml" -v down

## install-deps: install dependencies to local binary directory
.PHONY: install-deps
install-deps:
	GOBIN=$(temp_bin) go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(temp_bin) v1.61.0

## generate: generate all necessary code
.PHONY: generate
generate:
	$(temp_bin)/oapi-codegen --config=oapi_server.config.yaml api/openapi/${api_file_name}.yaml
	$(temp_bin)/oapi-codegen --config=oapi_models.config.yaml api/openapi/${api_file_name}.yaml
	protoc --go_out=. --go_opt=paths=import   --go-grpc_out=. --go-grpc_opt=paths=import   api/proto/${api_file_name}.proto

## clean: clean all temporary files
.PHONY: clean
clean:
	rm -rf $(temp)

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push