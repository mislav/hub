SOURCES = $(shell script/build files)

MIN_COVERAGE = 89.4

HELP_CMD = \
	share/man/man1/hub-alias.1 \
	share/man/man1/hub-browse.1 \
	share/man/man1/hub-ci-status.1 \
	share/man/man1/hub-compare.1 \
	share/man/man1/hub-create.1 \
	share/man/man1/hub-delete.1 \
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

bin/hub: $(SOURCES)
	script/build -o $@

test:
	go test ./...

test-all: bin/cucumber
ifdef CI
	script/test --coverage $(MIN_COVERAGE)
else
	script/test
endif

bin/cucumber:
	script/bootstrap

fmt:
	go fmt ./...

man-pages: $(HELP_ALL:=.md) $(HELP_ALL) $(HELP_ALL:=.txt)

%.txt: %
	groff -Wall -mtty-char -mandoc -Tutf8 -rLL=$(TEXT_WIDTH)n $< | col -b >$@

$(HELP_ALL): share/man/.man-pages.stamp
share/man/.man-pages.stamp: $(HELP_ALL:=.md)
	go run md2roff-bin/cmd.go --organization=GITHUB --manual="Hub Manual" share/man/man1/*.md
	touch $@

%.1.md: bin/hub
	bin/hub help $(*F) --plain-text >$@

share/man/man1/hub.1.md:
	true

install: bin/hub man-pages
	bash < script/install.sh

clean:
	git clean -fdx bin share/man

.PHONY: clean test test-all man-pages fmt install
