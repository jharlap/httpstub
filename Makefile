SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

TOOLS=github.com/wadey/gocovmerge honnef.co/go/staticcheck/cmd/staticcheck honnef.co/go/simple/cmd/gosimple golang.org/x/tools/cmd/cover honnef.co/go/unused/cmd/unused

.DEFAULT_GOAL: test

.PHONY: test
test:
	@# vet or staticcheck errors are unforgivable. gosimple produces warnings
	@go vet $(NON_VENDOR_PKGS)
	@staticcheck $(NON_VENDOR_PKGS)
	@unused -exported $(NON_VENDOR_PKGS)
	@-gosimple $(NON_VENDOR_PKGS)
	@go test $(testargs) ./...

.PHONY: install
install:
	go install ./...

.PHONY: bootstrap
bootstrap:
	$(foreach tool,$(TOOLS),$(call goget, $(tool)))

define goget
	go get -u $(1)
	
endef

