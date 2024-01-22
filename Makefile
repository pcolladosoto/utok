# The Go compiler to use.
GOC = go

# Configure this release's commit hash.
COMMIT = $(shell git rev-parse --short HEAD)

# Directory to dump built binaries to.
BIN_DIR = ./bin

# Linux Binary name.
BIN_NAME = utok

# Markdown-formatted manpage to parse with pandoc.
DOC_FILE = $(BIN_NAME).1.md

# Path of the buildroot created with rpmdev-setuptree.
RPM_BUILDROOT = $(shell echo ${HOME})/.rpmbuild

# Go sources to consider.
SOURCES = $(wildcard *.go)

# Files to delete after building stuff.
TRASH = $(BIN_DIR)/* *.rpm $(addprefix $(BIN_NAME), .1 -go .1.gz)

help:
		@echo "usage: make <target>"
		@echo "targets:"
		@echo "  linux  build for linux/amd64"
		@echo "  mac    build for macOS/amd64"
		@echo ""
		@echo "  doc    build and compress the manpage"
		@echo ""
		@echo "  rpm    build the RPM. Make sure the machine this runs on has a"
		@echo "         RPM buildroot configured through rpmdev-setuptree."
		@echo ""
		@echo "  clean  delete every built executable under $(BIN_DIR)"

linux: $(SOURCES)
		@echo "Building commit $(COMMIT) targeting Linux/amd64"
		@GOOS=linux GOARCH=amd64 $(GOC) build -o $(BIN_DIR)/$(BIN_NAME) \
				-ldflags "-X main.commit=$(COMMIT)"

mac: $(SOURCES)
		@echo "Building commit $(COMMIT) targeting macOS/amd64"
		@GOOS=darwin GOARCH=amd64 $(GOC) build -o $(BIN_DIR)/$(BIN_NAME).darwin \
				-ldflags "-X main.commit=$(COMMIT)"

doc: $(DOC_FILE)
	@echo "Building documentation"
	@pandoc --standalone --to man $(DOC_FILE) | gzip > $(basename $(DOC_FILE)).gz

.PHONY: rpm clean

rpm: linux doc
	@echo "Building RPM with buildroot $(RPM_BUILDROOT)"
	@echo "Copying artifacts to the RPM buildroot..."
	@cp $(BIN_DIR)/$(BIN_NAME)     $(RPM_BUILDROOT)/SOURCES/$(BIN_NAME)
	@cp $(basename $(DOC_FILE)).gz $(RPM_BUILDROOT)/SOURCES/$(basename $(DOC_FILE)).gz
	@cp $(BIN_NAME).spec           $(RPM_BUILDROOT)/SPECS/$(BIN_NAME).spec
	@echo "Building the RPM..."
	@rpmbuild -bb $(RPM_BUILDROOT)/SPECS/$(BIN_NAME).spec

clean:
		@rm -f $(TRASH)
