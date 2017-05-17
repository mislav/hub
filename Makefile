SOURCES = $(shell script/build files)
SOURCES_FMT = $(shell script/build files | cut -d/ -f1-2 | sort -u)

HELP_CMD = \
	share/man/man1/hub-alias.1 \
	share/man/man1/hub-browse.1 \
	share/man/man1/hub-ci-status.1 \
	share/man/man1/hub-compare.1 \
	share/man/man1/hub-create.1 \
	share/man/man1/hub-fork.1 \
	share/man/man1/hub-pr.1 \
	share/man/man1/hub-pull-request.1 \
	share/man/man1/hub-release.1 \
	share/man/man1/hub-issue.1 \
	share/man/man1/hub-sync.1 \

HELP_EXT = \
	share/man/man1/hub-am.1 \
	share/man/man1/hub-apply.1 \
	share/man/man1/hub-checkout.1 \
	share/man/man1/hub-cherry-pick.1 \
	share/man/man1/hub-clone.1 \
	share/man/man1/hub-fetch.1 \
	share/man/man1/hub-help.1 \
	share/man/man1/hub-init.1 \
	share/man/man1/hub-merge.1 \
	share/man/man1/hub-push.1 \
	share/man/man1/hub-remote.1 \
	share/man/man1/hub-submodule.1 \

HELP_ALL = share/man/man1/hub.1 $(HELP_CMD) $(HELP_EXT)

TEXT_WIDTH = 87

.PHONY: clean test test-all man-pages fmt install

all: bin/hub

bin/hub: $(SOURCES)
	script/build -o $@

test:
	script/build test

test-all: bin/cucumber
	script/test

bin/ronn bin/cucumber:
	script/bootstrap

fmt:
	go fmt $(filter %.go,$(SOURCES_FMT))
	go fmt $(filter-out %.go,$(SOURCES_FMT))

man-pages: $(HELP_ALL:=.ronn) $(HELP_ALL) $(HELP_ALL:=.txt)

%.txt: %.ronn
	groff -Wall -mtty-char -mandoc -Tutf8 -rLL=$(TEXT_WIDTH)n $< | col -b >$@

%.1: %.1.ronn bin/ronn
	bin/ronn --organization=GITHUB --manual="Hub Manual" share/man/man1/*.ronn

%.1.ronn: bin/hub
	bin/hub help $(*F) --plain-text | script/format-ronn $(*F) $@

share/man/man1/hub.1.ronn:
	true

install: bin/hub man-pages
	bash < script/install.sh

clean:
	git clean -fdx bin share/man
