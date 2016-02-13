.PHONY : build clean fmt release

NAME=docker-machine-driver-vultr
VERSION:=$(shell git describe --abbrev=0 --tags)

#ifdef $CIRCLE_BUILD_NUM
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

dist-clean:
	rm -rf dist
	rm -rf release

dist: dist-clean
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -a -ldflags "$(LDFLAGS)" -o dist/linux/amd64/$(NAME) ./bin 
	mkdir -p dist/linux/armhf && GOOS=linux GOARCH=arm GOARM=6 go build -a -ldflags "$(LDFLAGS)" -o dist/linux/armhf/$(NAME) ./bin 
	mkdir -p dist/darwin/amd64 && GOOS=darwin GOARCH=amd64 go build -a -ldflags "$(LDFLAGS)" -o dist/darwin/amd64/$(NAME) ./bin 
	mkdir -p dist/windows/amd64 && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags "$(LDFLAGS)" -o dist/windows/amd64/$(NAME).exe ./bin 

release: dist
	mkdir -p release
	tar -cvzf release/$(NAME)-$(VERSION)-linux-amd64.tar.gz -C dist/linux/amd64 $(NAME)
	tar -cvzf release/$(NAME)-$(VERSION)-linux-armhf.tar.gz -C dist/linux/armhf $(NAME)
	tar -cvzf release/$(NAME)-$(VERSION)-darwin-amd64.tar.gz -C dist/darwin/amd64 $(NAME)
	tar -cvzf release/$(NAME)-$(VERSION)-windows-amd64.tar.gz -C dist/windows/amd64 $(NAME).exe
	ghr -u janeczku -r docker-machine-vultr --replace $(VERSION) release/

get-deps:
	go get -u github.com/tcnksm/ghr
	go get -u github.com/tools/godep
	go get -u github.com/ChimeraCoder/tokenbucket
	go get -u github.com/JamesClonk/vultr
	go get -u github.com/docker/machine
	go get -u github.com/docker/docker/pkg/term
	go get -u golang.org/x/crypto/ssh
	go get -u golang.org/x/crypto/ssh/terminal
	go get -u github.com/Azure/go-ansiterm
	go get -u github.com/Sirupsen/logrus

check-gofmt:
	if [ -n "$(shell gofmt -l .)" ]; then \
		echo 1>&2 'The following files need to be formatted:'; \
		gofmt -l .; \
		exit 1; \
	fi

test:
	go vet .
	go test -race ./...