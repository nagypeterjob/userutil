UNAME := $(shell uname)

ifeq ($(UNAME), Linux)
target := linux
endif
ifeq ($(UNAME), Darwin)
target := darwin
endif


COMMIT_HASH=$$(git rev-list -1 HEAD)
TAG_VERSION=$$(git tag --sort=committerdate | tail -1)

test:
	go test -count=1 -race -cover -v $(shell go list ./... | grep -v -e /vendor/)

build:
	GOOS=$(target) go build -o "bin/userutil" 

lint:
	golint -set_exit_status `go list ./...`

golangci:
	golangci-lint \
	 --enable bodyclose \
	 --enable golint \
	 --enable gosec \
	 --enable unconvert \
	 --enable dupl \
	 --enable goconst \
	 --enable gocyclo \
	 --enable gocognit \
	 --enable maligned \
	 --enable unparam  \
	 --enable dogsled \
	 --enable nakedret \
	 --enable prealloc \
	 --enable gocritic \
	 --enable gomnd \
	 run ./...

vendor:
	go mod vendor

tidy:
	go mod tidy