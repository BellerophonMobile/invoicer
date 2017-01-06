.PHONY: all install test clean

SCRIPT_DIR = $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

export GOPATH = $(SCRIPT_DIR)
export WORKSPACE = $(SCRIPT_DIR)

all: install

install:
	go install invoicer

fmt:
	go fmt invoicer

test:
	go vet invoicer
	@if [ "$$(gofmt -l src/invoicer | wc -l)" -gt 0 ]; then echo "run go fmt"; exit 1; fi

clean:
	-rm -r bin pkg
