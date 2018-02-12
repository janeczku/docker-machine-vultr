.PHONY : build dist dist-clean release tag-release deps vet test lint

NAME=docker-machine-driver-vultr
VERSION := $(shell cat VERSION)

ifneq ($(CIRCLE_BUILD_NUM),)
	BUILD:=$(VERSION)-$(CIRCLE_BUILD_NUM)
else
	BUILD:=$(VERSION)
endif

LDFLAGS:=-X main.Version=$(VERSION)

all: build

build:
	mkdir -p build
	go build -a -ldflags "$(LDFLAGS)" -o build/$(NAME)-$(BUILD) ./cmd/docker-machine-vultr

tag-release:
	git tag -f $(VERSION)
	git push -f origin master --tags

deps:
	go get -v github.com/tcnksm/ghr
	go get -v github.com/golang/lint/golint

lint:
	@golint $$(go list ./... 2> /dev/null | grep -v /vendor/)

vet:
	@go vet $$(go list ./... 2> /dev/null | grep -v /vendor/)

test: lint vet
	@go test $$(go list ./... 2> /dev/null | grep -v /vendor/)

dist-clean:
	rm -rf release

dist: dist-clean
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -a -ldflags "$(LDFLAGS)" -o release/$(NAME)-Linux-x86_64 ./bin
	GOOS=linux GOARCH=arm GOARM=6 go build -a -ldflags "$(LDFLAGS)" -o release/$(NAME)-Linux-armhf ./bin
	GOOS=darwin GOARCH=amd64 go build -a -ldflags "$(LDFLAGS)" -o release/$(NAME)-Darwin-x86_64 ./bin
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags "$(LDFLAGS)" -o release/$(NAME)-Windows-x86_64.exe ./bin
	for file in release/$(NAME)-*; do openssl dgst -md5 < $${file} > $${file}.md5; done

release: dist
	ghr -u janeczku -r docker-machine-vultr --replace $(VERSION) release/
