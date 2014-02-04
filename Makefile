ifeq ($(TMPDIR),)
	TMPDIR := /tmp
endif

TMP := $(TMPDIR)/go
SELF := $(TMP)/src/github.com/jingweno/gh
export GOPATH := $(TMP):$(PWD)/Godeps/_workspace

objects = main.go cmd/*.go commands/*.go git/*.go github/*.go utils/*.go

test: $(SELF)
	go test -cover -v ./...

hub: $(SELF) $(objects)
	go build -o $@

fmt:
	go fmt ./...

$(SELF):
	mkdir -p $(dir $@)
	ln -snf "$(PWD)" $@

clean:
	rm -f hub
	rm -rf $(TMP)
