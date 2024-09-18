SRC := $(shell find . -name '*.go')
BIN := bin/zoe

.PHONY: all clean test run build upgrade help prologue

all: 					# default action
	@pre-commit install --install-hooks
	@git config commit.template .git-commit-template

clean:					# clean-up environment
	@find . -name '*.sw[po]' -delete
	rm -f $(BIN)

test: prologue			# run test
	go test -v ./...

build: prologue $(SRC) 	# build the binary/library
	go build -ldflags "-s -w" -o $(BIN) cmd/zoe/main.go

upgrade:				# upgrade all the necessary packages
	pre-commit autoupdate

help:					# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

prologue: $(SRC)
	@go mod tidy
	@gofmt -s -w $(SRC)
