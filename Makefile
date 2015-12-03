
HARDWARE=$(shell uname -m)
VERSION=$(shell awk '/const Version/ { print $$4 }' version.go | sed 's/"//g')
DEPS=$(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
PACKAGES=$(shell go list ./...)
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr

.PHONY: test build static deps vet cover docker clean authors

default: build

build:
	mkdir -p bin
	go build -o bin/prometheus-k8s

static:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o bin/prometheus-k8s

deps:
	@echo "--> Installing build dependencies"
	@go get -d -v ./... $(DEPS)

vet:
	@echo "--> Running go tool vet"
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@go tool vet $(VETARGS) .

lint:
	@echo "--> Running golint"
	@which golint 2>/dev/null ; if [ $$? -eq 1 ]; then \
		go get -u github.com/golang/lint/golint; \
	fi
	@golint .

format:
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

cover:
	@echo "--> Running go coverage"
	go list github.com/${AUTHOR}/${NAME} | xargs -n1 go test --cover

docker: build
	@echo "--> Running a docker build"
	sudo docker build -t ${AUTHOR}/${NAME} .

test: deps
	@$(MAKE) vet
	@$(MAKE) format
	@$(MAKE) lint
	@echo "--> Running go tests"
	go test -v

clean:
	rm -rf ./bin

authors:
	git log --format='%aN <%aE>' | sort -u > AUTHORS

changelog: release
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B > changelog
