VERSION=3.0.1
PREFIX=target/data-matcher-$(VERSION)

# go command for linux and windows.
GO=CGO_ENABLED=0 go
PARAMS=-ldflags '-s -w -extldflags "-static"'

# upx is a tool to compress executable program.
UPX=upx

PRGS=data-matcher kafka2nats


all:    $(PRGS)

data-matcher:
	$(GO) build $(PARAMS) -o $@ ./bin/data-matcher

kafka2nats:
	$(GO) build $(PARAMS) -o $@ ./bin/kafka2nats

clean:
	rm -f $(PRGS)

.PHONY: ./test
test:
	$(GO) test ./engine ./matcher ./model ./test

install:
	$(UPX) $(PRGS) || echo $?
	mkdir -p $(PREFIX)/etc
	cp -a etc/*.tpl $(PREFIX)/etc
	cp -a  Changelog.md $(PRGS) $(PREFIX)

	cd `dirname $(PREFIX)` && tar cvfz `basename $(PREFIX)`.tar.gz `basename $(PREFIX)`

