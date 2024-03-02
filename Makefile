PREFIX=/opt/idss
VERSION=1.0

# go command for linux and windows.
GO=CGO_ENABLED=0 go
PARAMS=-ldflags '-s -w -extldflags "-static"'

# upx is a tool to compress executable program.
UPX=upx

PRGS=data-matcher


all:	$(PRGS)

data-matcher:
	$(GO) build $(PARAMS) -o $@ ./bin/data-matcher

clean:
	rm -f $(PRGS)

install:
	$(UPX) $(PRGS) || echo $?
	mkdir -p $(PREFIX)-$(VERSION)/etc
	ln -snf $(PREFIX)-$(VERSION) $(PREFIX)
	cp -a etc/*.tpl $(PREFIX)/etc
	cp -a  Changelog.md $(PRGS) $(PREFIX)
