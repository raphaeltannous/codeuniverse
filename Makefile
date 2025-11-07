default: help

PROJECTNAME=$(shell basename "$(PWD)")

CLI_MAIN_FOLDER=./cmd/server
BIN_FOLDER=bin
BIN_NAME=${PROJECTNAME}

MIGRATIONS_DIR=./migrations
MIGRATIONS_ENV=./migrations/.env

## setup: install all build dependencies for ci
setup: mod-download

## compile: compiles project in current system
compile: clean generate fmt vet test build

## watch: format, test, build and run the project
watch:
	@echo "  >  Watching go files..."
	@if !type "entr" > /dev/null 2>&1; then \
		echo "Please install entr: http://eradman.com/entrproject/"; \
	else \
		make migrate; \
		find . -path ./.database -prune -or -type f -name "*.go" -print \
			| entr -nr make clean generate fmt vet build run; \
	fi

## migrate: migrate the database to the lastest version
migrate:
	@echo " > Migrating database..."
	@if !type "goose" > /dev/null 2>&1; then \
		echo " > goose not found, installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
	fi
	@echo " > Running migrations"
	@goose -env ${MIGRATIONS_ENV} -dir ${MIGRATIONS_DIR} up

## run: build project and run it
run: compile
	@echo " > Running binary"
	./${BIN_FOLDER}/${BIN_NAME}

clean:
	@echo "  >  Cleaning build cache"
	@-rm -rf ${BIN_FOLDER}/amd64 ${BIN_FOLDER}/${BIN_NAME} \
		&& go clean ./...

build:
	@echo "  >  Building binary"
	@mkdir -p ${BIN_FOLDER}
	@go build \
		-o ${BIN_FOLDER}/${BIN_NAME} \
		"${CLI_MAIN_FOLDER}"

fmt:
	@echo "  >  Formatting code"
	@go fmt ./...

generate:
	@echo "  >  Go generate"
	@if !type "stringer" > /dev/null 2>&1; then \
		go install golang.org/x/tools/cmd/stringer@latest; \
	fi
	@go generate ./...

mod-download:
	@echo "  >  Download dependencies..."
	@go mod download && go mod tidy

test:
	@echo "  >  Executing unit tests"
	@go test -v -timeout 60s -race ./...

vet:
	@echo "  >  Checking code with vet"
	@go vet ./...

.PHONY: help
all: help
help: Makefile
	@echo "Choose a command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
