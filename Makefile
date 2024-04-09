SCSS := $(shell find . -name '*.scss')
CSS  := $(SCSS:.scss=.css)

SRC := $(shell find . -name '*.go')
BIN := zoe

.PHONY: all clean test run build upgrade help

all: 			# default action
	@pre-commit install --install-hooks
	@git config commit.template .git-commit-template

clean:			# clean-up environment
	@find . -name '*.sw[po]' -delete
	rm -f $(BIN)

test:			# run test
	go test -v ./...

run:			# run in the local environment
	go run

build: $(SRC) $(CSS)	# build the binary/library
	go build -ldflags "-s -w" -o $(BIN) cmd/$(BIN)/main.go

upgrade:		# upgrade all the necessary packages
	pre-commit autoupdate

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

test run build: linter
linter:
	@go mod tidy
	@gofmt -s -w $(SRC)

%.css: %.scss
	sass $< $@
