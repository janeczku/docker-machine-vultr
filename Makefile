.PHONY : build dist dist-clean fmt release tag-release deps vet test

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
	godep go build -a -ldflags "$(LDFLAGS)" -o build/$(NAME)-$(BUILD) ./bin

dist-clean:
	rm -rf dist
	rm -rf release

dist: dist-clean
	mkdir -p release
	mkdir -p dist
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 godep go build -a -ldflags "$(LDFLAGS)" -o dist/linux/amd64/$(NAME) ./bin
	mkdir -p dist/linux/armhf && GOOS=linux GOARCH=arm GOARM=6 godep go build -a -ldflags "$(LDFLAGS)" -o dist/linux/armhf/$(NAME) ./bin
	mkdir -p dist/darwin/amd64 && GOOS=darwin GOARCH=amd64 godep go build -a -ldflags "$(LDFLAGS)" -o dist/darwin/amd64/$(NAME) ./bin
	mkdir -p dist/windows/amd64 && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 godep go build -a -ldflags "$(LDFLAGS)" -o dist/windows/amd64/$(NAME).exe ./bin
	tar -cvzf release/$(NAME)-$(VERSION)-linux-amd64.tar.gz -C dist/linux/amd64 $(NAME)
	cd $(shell pwd)/release && md5sum $(NAME)-$(VERSION)-linux-amd64.tar.gz > $(NAME)-$(VERSION)-linux-amd64.tar.gz.md5
	tar -cvzf release/$(NAME)-$(VERSION)-linux-armhf.tar.gz -C dist/linux/armhf $(NAME)
	cd $(shell pwd)/release && md5sum $(NAME)-$(VERSION)-linux-armhf.tar.gz > $(NAME)-$(VERSION)-linux-armhf.tar.gz.md5
	tar -cvzf release/$(NAME)-$(VERSION)-darwin-amd64.tar.gz -C dist/darwin/amd64 $(NAME)
	cd $(shell pwd)/release && md5sum $(NAME)-$(VERSION)-darwin-amd64.tar.gz > $(NAME)-$(VERSION)-darwin-amd64.tar.gz.md5
	tar -cvzf release/$(NAME)-$(VERSION)-windows-amd64.tar.gz -C dist/windows/amd64 $(NAME).exe
	cd $(shell pwd)/release && md5sum $(NAME)-$(VERSION)-windows-amd64.tar.gz > $(NAME)-$(VERSION)-windows-amd64.tar.gz.md5

release: dist
	ghr -u janeczku -r docker-machine-vultr --replace $(VERSION) release/

tag-release:
	git tag -f `cat VERSION`
	git push -f origin master --tags

deps:
	go get -u github.com/tools/godep
	go get -u github.com/tcnksm/ghr

vet:
	@if [ -n "$(shell gofmt -l .)" ]; then \
		echo 1>&2 'The following files need to be formatted:'; \
		gofmt -l .; \
		exit 1; \
	fi
	godep go vet .

test:
	godep go test -race ./...