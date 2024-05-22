parentFolder := $(shell pwd)
SHELL := /bin/bash

all: buildTk install

buildTk:
	cd $(parentFolder)/malinoTk; \
	go mod tidy; \
	go build -o $(parentFolder)/malino

install:
	sudo cp malino /usr/bin/malino