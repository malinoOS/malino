parentFolder := $(shell pwd)
SHELL := /bin/bash

.PHONY: all toolkit

all: toolkit deb install

toolkit:
	cd $(parentFolder)/toolkit; \
	go mod tidy; \
	go build -o $(parentFolder)/malino -ldflags "-X main.Version=$(shell date +%y%m%d)"

deb:
	if mkdir malino-deb; then \
		echo Building deb...; \
	else \
		sudo rm -rf malino-deb; \
		mkdir malino-deb; \
	fi
	-rm *.deb
	mkdir malino-deb/DEBIAN
	mkdir malino-deb/usr
	mkdir malino-deb/usr/bin
	sudo printf 'Package: malino\nVersion: $(shell date +%y%m%d)\nArchitecture: all\nDepends: golang-go, qemu-system-x86, qemu-utils\nMaintainer: Winksplorer <winksplorer@gordae.com>\nDescription: The Malino Linux-based OS development tookit\n' | sudo tee malino-deb/DEBIAN/control
	cp $(parentFolder)/malino malino-deb/usr/bin/malino
	sudo chown -R root:root malino-deb
	sudo chmod 0755 malino-deb/DEBIAN
	sudo chmod 0755 malino-deb/usr/bin/*
	sudo dpkg-deb --build malino-deb
	mv malino-deb.deb malino-stable-$(shell date +%y%m%d)-any.deb
	sudo chown $(shell whoami) malino-stable-$(shell date +%y%m%d)-any.deb
	sudo rm -rf malino-deb

install:
	sudo cp malino /usr/bin/malino