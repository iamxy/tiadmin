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

# Ensure we only have one entry in GOPATH
GOPATH := $(GITHUB_ROOT)/../..
#                        ^  ^~ GOPATH
#                        |~ GOPATH/src

GOPATH_BIN := $(GOPATH)/bin
GO_BINDATA := $(GOPATH_BIN)/go-bindata

REMOTE_DEPS             := npm.installed grunt.built
INDEX                   := index.html
GOBINDATA_TARGET        := web-ui.go


.PHONY: all
all: $(GOBINDATA_TARGET)

npm.installed:
	npm install --progress=false
	touch $@

grunt.built: Gruntfile.js npm.installed
	grunt build
	touch $@

$(GO_BINDATA):
	go get -u github.com/jteeuwen/go-bindata/...

$(GOBINDATA_TARGET): $(GO_BINDATA) grunt.built
	$(GO_BINDATA) -nometadata -pkg frontend -prefix "$(DIST)/" -o $@ $(DIST)/...
	gofmt -s -w $@
	goimports -w $@

.PHONY: clean
clean:
	rm -f $(REMOTE_DEPS)
	rm -rf $(NODE_MODULES) $(BOWER_COMPONENTS) $(DIST)
	rm -rf $(GOBINDATA_TARGET)
