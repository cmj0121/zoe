.PHONY: all clean test run build upgrade run stop help

all: 			# default action

clean:			# clean-up environment

test:			# run test

build:			# build the binary/library
	docker compose build

upgrade:		# upgrade all the necessary packages

run:			# run the service in local
	docker compose up -d

stop:			# stop the service in local
	docker compose stop

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'
