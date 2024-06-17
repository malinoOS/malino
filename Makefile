parentFolder := $(shell pwd)
SHELL := /bin/bash

.PHONY: all toolkit deb libmsb libmalino-cs cleanallsubs install deps

all: help

help:
	@printf "malino Makefile\n\n"
	@printf "make stable - Builds malino, and creates a debian package file\n"
	@printf "make dev    - Builds malino, and installs it onto the current system\n"
	@printf "make deps   - !THIS ONLY WORKS ON DEBIAN-BASED DISTROS! Installs dependencies needed for malino to build.\n"

# stable target builds the toolkit and creates a debian package file
stable: toolkit libmsb libmalino-cs deb cleanallsubs
# dev target is for during development, it makes the toolkit and installs it every time
dev: toolkit libmsb libmalino-cs install cleanallsubs
# deps target installs dependencies needed for malino to build, only works on debian
deps:
	wget https://dot.net/v1/dotnet-install.sh -O dotnet-install.sh
	chmod +x ./dotnet-install.sh
	./dotnet-install.sh --channel 8.0
	rm ./dotnet-install.sh
	sudo apt install golang-go qemu-system-x86 p7zip

toolkit:
	@echo " GO malino"
	@cd $(parentFolder)/toolkit; \
	go mod tidy; \
	go build -o $(parentFolder)/malino -ldflags "-X main.Version=$(shell date +%y%m%d)"

libmsb:
	@echo " MK libmsb"
	@make -C libmsb/

libmalino-cs:
	@echo " MK libmalino-cs"
	@make -C libmalino-cs/

cleanallsubs:
	@rm -f libmsb/*.o libmsb/libmsb.so

	@rm -f libmalino-cs/libmalino-cs.dll libmalino-cs/libmalino-cs.deps.json libmalino-cs/libmalino-cs.pdb
	@rm -rf libmalino-cs/bin
	@rm -rf libmalino-cs/obj

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
	@echo " MK malino-deb/usr/bin"
	@mkdir -p malino-deb/usr/bin
	@echo " MK malino-deb/opt/malino"
	@mkdir -p malino-deb/opt/malino

	@echo "  W malino-deb/DEBIAN/control"
	@sudo printf 'Package: malino\nVersion: $(shell date +%y%m%d)\nArchitecture: all\nDepends: golang-go, qemu-system-x86, qemu-utils, 7zip | p7zip\nMaintainer: Winksplorer <winksplorer@gordae.com>\nDescription: The Malino Linux-based OS development tookit\n' | sudo tee malino-deb/DEBIAN/control > /dev/null
	
	@echo " MV malino malino-deb/usr/bin/malino"
	@mv $(parentFolder)/malino malino-deb/usr/bin/malino

	@echo " MV libmsb/libmsb.so malino-deb/opt/malino/libmsb.so"
	@sudo mv libmsb/libmsb.so malino-deb/opt/malino/libmsb.so

	@echo " MV libmalino-cs/libmalino-cs.dll malino-deb/opt/malino/libmalino-cs.dll"
	@sudo mv libmalino-cs/libmalino-cs.dll malino-deb/opt/malino/libmalino-cs.dll

	@echo " CH malino-deb TO root:root"
	@sudo chown -R root:root malino-deb
	@echo " CH malino-deb TO 0755"
	@sudo chmod -R 0755 malino-deb/DEBIAN

	@echo "DEB malino-stable-$(shell date +%y%m%d)-any.deb"
	@sudo dpkg-deb --build malino-deb > /dev/null

	@mv malino-deb.deb malino-stable-$(shell date +%y%m%d)-any.deb
	@echo " MV malino-stable-$(shell date +%y%m%d)-any.deb TO $(shell whoami)"
	@sudo chown $(shell whoami) malino-stable-$(shell date +%y%m%d)-any.deb

	@echo " RM malino-deb"
	@sudo rm -rf malino-deb

install:
	@echo " MV malino /usr/bin/malino"
	@sudo mv malino /usr/bin/malino
	@echo " MK /opt/malino"
	@sudo mkdir -p /opt/malino
	@echo " MV libmsb/libmsb.so /opt/malino/libmsb.so"
	@sudo mv libmsb/libmsb.so /opt/malino/libmsb.so
	@echo " MV libmalino-cs/libmalino-cs.dll /opt/malino/libmalino-cs.dll"
	@sudo mv libmalino-cs/libmalino-cs.dll /opt/malino/libmalino-cs.dll