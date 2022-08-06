# =========================================================================== #
#            MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>      #
#                                                                             #
#                 ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                 #
#                 ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                 #
#                 ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                 #
#                 ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                 #
#                 ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                 #
#                 ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                 #
#                 ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                 #
#                 ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                 #
#                                                                             #
#                        This machine kills fascists.                         #
#                                                                             #
# =========================================================================== #

all: compile

version       ?=  0.1.2
toplevel      ?=  aurae
auraetarget   ?=  aurae
auraedtarget  ?=  auraed
auraefstarget ?=  auraefs
org           ?=  kris-nova
authorname    ?=  Kris Nóva
authoremail   ?=  kris@nivenly.com
license       ?=  MIT
year          ?=  2022
copyright     ?=  Copyright (c) $(year)

firecrackerversion  =  1.1.1
firecrackertarball  =  firecracker-v$(firecrackerversion)-x86_64.tgz

compile: gen aurae auraed auraefs ## Compile for the local architecture ⚙

.PHONY: aurae
aurae: ## Compile aurae (local arch)
	@echo "Compiling [aurae] ..."
	go build -v -ldflags "\
	-X 'github.com/$(org)/$(toplevel).Version=$(version)' \
	-X 'github.com/$(org)/$(toplevel).AuthorName=$(authorname)' \
	-X 'github.com/$(org)/$(toplevel).AuthorEmail=$(authoremail)' \
	-X 'github.com/$(org)/$(toplevel).Copyright=$(copyright)' \
	-X 'github.com/$(org)/$(toplevel).License=$(license)'" \
	-o bin/$(auraetarget) cmd/aurae/*.go

.PHONY: auraed
auraed: ## Compile auraed (local arch)
	@echo "Compiling [auraed] ..."
	go build -v -ldflags "\
	-X 'github.com/$(org)/$(toplevel).Version=$(version)' \
	-X 'github.com/$(org)/$(toplevel).AuthorName=$(authorname)' \
	-X 'github.com/$(org)/$(toplevel).AuthorEmail=$(authoremail)' \
	-X 'github.com/$(org)/$(toplevel).Copyright=$(copyright)' \
	-X 'github.com/$(org)/$(toplevel).License=$(license)'" \
	-o bin/$(auraedtarget) cmd/auraed/*.go

.PHONY: auraefs
auraefs: ## Compile auraefs (local arch)
	@echo "Compiling [auraefs] ..."
	go build -v -ldflags "\
	-X 'github.com/$(org)/$(toplevel).Version=$(version)' \
	-X 'github.com/$(org)/$(toplevel).AuthorName=$(authorname)' \
	-X 'github.com/$(org)/$(toplevel).AuthorEmail=$(authoremail)' \
	-X 'github.com/$(org)/$(toplevel).Copyright=$(copyright)' \
	-X 'github.com/$(org)/$(toplevel).License=$(license)'" \
	-o bin/$(auraefstarget) cmd/auraefs/*.go

install: ## Install the program to /usr/bin 🎉
	@echo "Installing..."
	sudo cp -v bin/$(auraetarget) /usr/bin/$(auraetarget)
	sudo cp -v bin/$(auraefstarget) /usr/bin/$(auraefstarget)
	sudo cp -v bin/$(auraedtarget) /usr/bin/$(auraedtarget)


gen: generate ## Alias for generate
generate: ## Will generate Go code from auraefs .proto files
	@echo "Generating..."
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	rpc/*.proto

.PHONY: test
test: compile ## 🤓 Run go tests
	@echo "Testing..."
	go test -v $(shell go list ./... | grep -v auraefs)

.PHONY: test-auraefs
test-auraefs:
	@echo "Testing auraefs..."
	go test -v ./auraefs

clean: ## Clean your artifacts 🧼
	@echo "Cleaning..."
	rm -rvf release/*
	rm -rvf rpc/*.pb.*
	rm -rvf bin/*

# auare-firecracker
# auare-firecracker-jailer
# auare-firecracker-rebase-snap
# auare-firecracker-seccompiler-bin
firecracker: ## Install firecracker for Aurae (x86)
	@echo "Installing Firecracker $(firecrackerversion)..."
	@if [ ! -f $(firecrackertarball) ]; then wget https://github.com/firecracker-microvm/firecracker/releases/download/v$(firecrackerversion)/firecracker-v$(firecrackerversion)-x86_64.tgz; fi
	@tar -xzf $(firecrackertarball)
	cd release-v$(firecrackerversion)-x86_64 && cp -v firecracker-v* /bin/auare-firecracker
	cd release-v$(firecrackerversion)-x86_64 && cp -v jailer* /bin/aurae-firecracker-jailer
	cd release-v$(firecrackerversion)-x86_64 && cp -v rebase-snap* /bin/aurae-firecracker-rebase-snap
	cd release-v$(firecrackerversion)-x86_64 && cp -v seccompiler-bin* /bin/aurae-firecracker-seccompiler-bin

.PHONY: help
help:  ## 🤔 Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'