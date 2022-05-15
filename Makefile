EXECUTABLES = git go find pwd wget
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

VERSION ?= $(shell git describe --tags `git rev-list --tags --max-count=1`)
BINARY = itool
MAIN = main.go

BUILDDIR = dist
GITREV = $(shell git rev-parse --short HEAD)
BUILDTIME = $(shell date +'%FT%TZ%z')
GO_BUILDER_VERSION=latest

deps:
ifeq ($(wildcard frida/linux/libfrida-core.a),)
	wget -O frida_linux.tar.xz https://github.com/frida/frida/releases/download/15.1.22/frida-core-devkit-15.1.22-linux-x86_64.tar.xz
	tar -xf frida_linux.tar.xz -C frida/linux
	rm -f frida_linux.tar.xz
endif

ifeq ($(wildcard frida/macos/libfrida-core.a),)
	wget -O frida_macos.tar.xz https://github.com/frida/frida/releases/download/15.1.22/frida-core-devkit-15.1.22-macos-x86_64.tar.xz
	tar -xf frida_macos.tar.xz -C frida/macos
	rm -f frida_macos.tar.xz
endif

snapshot: deps
	docker run --rm --privileged \
	-v $(CURDIR):/itool \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v $(GOPATH)/src:/go/src \
	-w /itool \
	ghcr.io/gythialy/golang-cross:$(GO_BUILDER_VERSION) --snapshot --rm-dist

release: deps
	docker run --rm --privileged \
	-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
	-v $(CURDIR):/itool \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v $(GOPATH)/src:/go/src \
	-w /itool \
	ghcr.io/gythialy/golang-cross:$(GO_BUILDER_VERSION) --rm-dist