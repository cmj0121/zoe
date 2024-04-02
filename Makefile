.PHONY: all clean test run build upgrade help

all: 			# default action
	@pre-commit install --install-hooks
	@git config commit.template .git-commit-template

clean:			# clean-up environment
	@find . -name '*.sw[po]' -delete

test:			# run test

run:			# run in the local environment

build:			# build the binary/library

upgrade:		# upgrade all the necessary packages
	pre-commit autoupdate

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'
