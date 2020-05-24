#
# Copyright (c) 2015, Matt Jones <urlugal@gmail.com>
# All rights reserved.
# MIT License
# For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
#
# version 0.1.23
#
 SHELL = /bin/bash

# TODO: document Makefile


.PHONY: all build clean coverage help install package pretty test

# The name of the binary to build
#
ifndef pkg
pkg := $(shell pwd | awk -F/ '{print $$NF}')
endif

# Set the target OS
# Ex: windows, darwin, linux
#
ifndef target_os
	#target_os = linux
	ifeq ($(OS),Windows_NT)
		target_os = windows
		ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
			target_arch = amd64
		else
			ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
				target_arch = amd64
			endif

			#ifeq ($(PROCESSOR_ARCHITECTURE),x86)
			#    target_arch = 386
			#endif
		endif
	else
		UNAME_S := $(shell uname -s)
		ifeq ($(UNAME_S),Linux)
			target_os = linux
		endif

		ifeq ($(UNAME_S),Darwin)
			target_os = darwin
		endif
		UNAME_P := $(shell uname -p)

		ifeq ($(UNAME_P),x86_64)
			target_arch = amd64
		endif

		#ifneq ($(filter %86,$(UNAME_P)),)
		#    target_arch = 386
		#endif
	endif
endif

ifeq ($(target_os),windows)
	target_ext = .exe
endif

# Set the target arch
# Ex: amd64, x86_64
#
ifndef target_arch
target_arch = amd64
endif


all: pretty clean build

# TODO: need to add pretty back in when I figure out how
build: prep
	@GOOS=$(target_os) GOARCH=$(target_arch) go build -o ./bin/$(pkg)-$(target_os)

release: prep
	@GOOS=$(target_os) GOARCH=$(target_arch) go build -ldflags="-s -w" -o ./bin/$(pkg)$(target_ext)

clean:
	@rm -rf ./bin ./rules

# TODO: write help command for Makefile
# TODO: documentation
help:

install: pretty
	@GOOS=$(target_os) GOARCH=$(target_arch) go install

package: test clean build

prep:
	@go get -u

pretty:
	@golint *.go
	@golint core/*.go
	@gofmt -w *.go	
	@gofmt -w core/*.go
	@go vet *.go
	@go vet core/*.go

test: pretty
	@cd ./$(pkg) && go test -cover

