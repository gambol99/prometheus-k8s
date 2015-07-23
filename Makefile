#
#   Author: Rohith
#   Date: 2015-07-20 19:40:22 +0100 (Mon, 20 Jul 2015)
#
#  vim:ts=2:sw=2:et
#

NAME=prometheus-k8s
AUTHOR=gambol99
VERSION=$(shell awk '/const version/ { print $$4 }' version.go | sed 's/"//g')

.PHONY: test

default: build

build:
	mkdir -p bin
	go build -o bin/prometheus-k8s

static:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o bin/prometheus-k8s

docker: build
	sudo docker build -t ${AUTHOR}/${NAME} .

test: build
	bin/prometheus-k8s -api=10.250.1.201 -port=8080 -logtostderr=true -v=3 -dry-run

clean:
	rm -rf ./bin

authors:
	git log --format='%aN <%aE>' | sort -u > AUTHORS

changelog: release
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B > changelog
