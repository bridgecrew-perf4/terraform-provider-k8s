#!make

TEST?=$$(${GOBIN} list ./... | grep -v 'vendor')
HOSTNAME=supply.com
NAMESPACE=devops
NAME=k8s
BINARY=terraform-provider-${NAME}
GOBIN=~/sdk/go1.15.2/bin/go
VERSION=1.7.1
OS_ARCH=darwin_amd64
SHELL=/bin/sh

default: install

build:
	${GOBIN} build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 ${GOBIN} build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	${GOBIN} test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 ${GOBIN} test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 ${GOBIN} test $(TEST) -v $(TESTARGS) -timeout 120m

init: ${GOBIN}.mod ${GOBIN}.sum

${GOBIN}.mod:
	-${GOBIN} mod init terraform-provider-k8s

${GOBIN}.sum:
	-${GOBIN} mod vendor

clean:
	${GOBIN} clean

scrub:
	${GOBIN} clean -cache -modcache

.PHONY: init