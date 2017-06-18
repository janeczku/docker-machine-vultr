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
	go build -a -ldflags "$(LDFLAGS)" -o build/$(NAME)-$(BUILD) ./bin

tag-release:
	git tag -f `cat VERSION`
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
	rm -rf dist
	rm -rf release

dist: dist-clean
	mkdir -p release
	mkdir -p dist
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -a -ldflags "$(LDFLAGS)" -o dist/linux/amd64/$(NAME) ./bin
	mkdir -p dist/linux/armhf && GOOS=linux GOARCH=arm GOARM=6 go build -a -ldflags "$(LDFLAGS)" -o dist/linux/armhf/$(NAME) ./bin
	mkdir -p dist/darwin/amd64 && GOOS=darwin GOARCH=amd64 go build -a -ldflags "$(LDFLAGS)" -o dist/darwin/amd64/$(NAME) ./bin
	mkdir -p dist/windows/amd64 && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags "$(LDFLAGS)" -o dist/windows/amd64/$(NAME).exe ./bin
	tar -cvzf release/$(NAME)-$(VERSION)-Linux-x86_64.tar.gz -C dist/linux/amd64 $(NAME)
	cd $(shell pwd)/release && md5sum $(NAME)-$(VERSION)-Linux-x86_64.tar.gz > $(NAME)-$(VERSION)-Linux-x86_64.tar.gz.md5
	tar -cvzf release/$(NAME)-$(VERSION)-Linux-armv7l.tar.gz -C dist/linux/armhf $(NAME)
	cd $(shell pwd)/release && md5sum $(NAME)-$(VERSION)-Linux-armv7l.tar.gz > $(NAME)-$(VERSION)-Linux-armv7l.tar.gz.md5
	tar -cvzf release/$(NAME)-$(VERSION)-Darwin-x86_64.tar.gz -C dist/darwin/amd64 $(NAME)
	cd $(shell pwd)/release && md5sum $(NAME)-$(VERSION)-Darwin-x86_64.tar.gz > $(NAME)-$(VERSION)-Darwin-x86_64.tar.gz.md5
	tar -cvzf release/$(NAME)-$(VERSION)-Windows-x86_64.tar.gz -C dist/windows/amd64 $(NAME).exe
	cd $(shell pwd)/release && md5sum $(NAME)-$(VERSION)-Windows-x86_64.tar.gz > $(NAME)-$(VERSION)-Windows-x86_64.tar.gz.md5

release: dist
	ghr -u janeczku -r docker-machine-vultr --replace $(VERSION) release/
