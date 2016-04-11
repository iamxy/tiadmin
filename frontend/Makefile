# NOTE: for some reason bash is necessary for updating the PATH to work
# See http://stackoverflow.com/questions/8941110/how-i-could-add-dir-to-path-in-makefile
SHELL            := /bin/bash

# Update the path to prefer binstubs over globals
PATH                 := node_modules/.bin:$(PATH)

REPO_ROOT            := $(realpath ..)
ORG_ROOT             := $(REPO_ROOT)/..
GITHUB_ROOT          := $(ORG_ROOT)/..
BOWER_COMPONENTS     := bower_components
NODE_MODULES         := node_modules
DIST                 := dist
APP                  := app

# Ensure we only have one entry in GOPATH
GOPATH := $(GITHUB_ROOT)/../..
#                        ^  ^~ GOPATH
#                        |~ GOPATH/src

GOPATH_BIN := $(GOPATH)/bin
GO_BINDATA := $(GOPATH_BIN)/go-bindata

REMOTE_DEPS             := npm.installed
INDEX                   := index.html
GOBINDATA_TARGET        := web-ui.go


.PHONY: all
all: $(GOBINDATA_TARGET)

npm.installed:
	npm install --progress=false
	touch $@

$(DIST): Gruntfile.js npm.installed $(APP)
	grunt build

$(GO_BINDATA):
	go get -u github.com/jteeuwen/go-bindata/...

$(GOBINDATA_TARGET): $(GO_BINDATA) $(DIST)
	$(GO_BINDATA) -nometadata -pkg frontend -prefix "$(DIST)/" -o $@ $(DIST)/...
	gofmt -s -w $@
	goimports -w $@

.PHONY: clean
clean:
	rm -rf $(DIST)
	rm -rf $(GOBINDATA_TARGET)
