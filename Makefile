.DEFAULT_GOAL := default

clean:
	@rm ./bin/*

default: install-deps build 

build: bin_dir
	@go build -o bin/snykctl ./cmd/main.go

install: build
	@cp bin/snykctl /usr/local/bin/snykctl

bin_dir: 
	@mkdir -p ./bin

install-deps: install-goimports

install-goimports:
	@if [ ! -f ./goimports ]; then \
		cd ~ && go get -u golang.org/x/tools/cmd/goimports; \
	fi

.PHONY: clean build 
