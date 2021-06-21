default: build

test:
	sh -x $(CURDIR)/make.sh test

build:
	sh -x $(CURDIR)/make.sh build

changelog:
	sh  $(CURDIR)/changelog.sh

stable:
	sh  $(CURDIR)/release.sh stable

beta:
	sh  $(CURDIR)/release.sh beta


.PHONY: clean test build release install changelog
