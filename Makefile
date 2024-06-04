parentFolder := $(shell pwd)
SHELL := /bin/bash

.PHONY: all toolkit

all: toolkit deb install

toolkit:
	@echo " GO malino"
	@cd $(parentFolder)/toolkit; \
	go mod tidy; \
	go build -o $(parentFolder)/malino -ldflags "-X main.Version=$(shell date +%y%m%d)"

deb:
	@if mkdir malino-deb; then \
		echo "deb" > /dev/null; \
	else \
		sudo rm -rf malino-deb; \
		mkdir malino-deb; \
	fi
	@-rm *.deb
	@echo " MK malino-deb/DEBIAN"
	@mkdir malino-deb/DEBIAN
	@echo " MK malino-deb/usr"
	@mkdir malino-deb/usr
	@echo " MK malino-deb/usr/bin"
	@mkdir malino-deb/usr/bin
	@echo "  W malino-deb/DEBIAN/control"
	@sudo printf 'Package: malino\nVersion: $(shell date +%y%m%d)\nArchitecture: all\nDepends: golang-go, qemu-system-x86, qemu-utils, 7zip | p7zip\nMaintainer: Winksplorer <winksplorer@gordae.com>\nDescription: The Malino Linux-based OS development tookit\n' | sudo tee malino-deb/DEBIAN/control > /dev/null
	@echo " CP malino malino-deb/usr/bin/malino"
	@cp $(parentFolder)/malino malino-deb/usr/bin/malino
	@echo " CH malino-deb TO root:root"
	@sudo chown -R root:root malino-deb
	@echo " CH malino-deb/DEBIAN TO 0755"
	@sudo chmod 0755 malino-deb/DEBIAN
	@echo " CH malino-deb/usr/bin TO 0755"
	@sudo chmod 0755 malino-deb/usr/bin/*
	@echo "DEB malino-stable-$(shell date +%y%m%d)-any.deb"
	@sudo dpkg-deb --build malino-deb > /dev/null
	@mv malino-deb.deb malino-stable-$(shell date +%y%m%d)-any.deb
	@echo " CH malino-stable-$(shell date +%y%m%d)-any.deb TO $(shell whoami)"
	@sudo chown $(shell whoami) malino-stable-$(shell date +%y%m%d)-any.deb
	@echo " RM malino-deb"
	@sudo rm -rf malino-deb

install:
	@echo " CP malino /usr/bin/malino"
	@sudo cp malino /usr/bin/malino