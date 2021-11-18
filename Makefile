.PHONY: all dep build clean check test test-coverage lint

BINARY="arctl"
PKG_LIST := $(shell go list ./...)

all: dep check test build

dep:
	@go mod download

build:
	@go build -o ${BINARY}

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

test:
	@go test -short ${PKG_LIST}

test-coverage:
	@go test -short -coverprofile coverage.out -covermode=atomic ${PKG_LIST}
	@cat coverage.out >> coverage.all

check:
	@go fmt ${PKG_LIST}
	@go vet ${PKG_LIST}

lint:
	@golint -set_exit_status ${PKG_LIST}
