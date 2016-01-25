SOURCES = $(shell script/build files)

HELP_CMD = \
	man/hub-alias.1 \
	man/hub-browse.1 \
	man/hub-ci-status.1 \
	man/hub-compare.1 \
	man/hub-create.1 \
	man/hub-fork.1 \
	man/hub-pull-request.1 \

HELP_EXT = \
	man/hub-am.1 \
	man/hub-apply.1 \
	man/hub-checkout.1 \
	man/hub-cherry-pick.1 \
	man/hub-clone.1 \
	man/hub-fetch.1 \
	man/hub-init.1 \
	man/hub-merge.1 \
	man/hub-push.1 \
	man/hub-remote.1 \
	man/hub-submodule.1 \

HELP_ALL = man/hub.1 $(HELP_CMD) $(HELP_EXT)

TEXT_WIDTH = 87

.PHONY: clean test test-all man-pages

all: bin/hub

bin/hub: $(SOURCES)
	script/build -o $@

test:
	script/build test

test-all: bin/cucumber
	script/test

bin/ronn bin/cucumber:
	script/bootstrap

man-pages: $(HELP_ALL:=.ronn) $(HELP_ALL) $(HELP_ALL:=.txt)

%.txt: %.ronn
	groff -Wall -mtty-char -mandoc -Tutf8 -rLL=$(TEXT_WIDTH)n $< | col -b >$@

%.1: %.1.ronn bin/ronn
	bin/ronn --organization=GITHUB --manual="Hub Manual" man/*.ronn

%.1.ronn: bin/hub
	bin/hub help $(*F) --plain-text | script/format-ronn $(*F) $@

man/hub.1.ronn:
	true

clean:
	rm -rf bin/hub
	git clean -fdx man/
