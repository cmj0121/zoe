.PHONY: all clean test run stop build help

all: 			# default action

clean:			# clean-up environment
	docker compose down -v

test:			# run test

run:			# run in the local environment
	docker compose up -d

stop:
	docker compose down

build:			# build the binary/library
	docker compose build

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'
