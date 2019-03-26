OUT_DIR = _output
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=sidecar
PLATFORM=GOARCH=amd64 GOOS=linux

#export $(GOARCH)
#export $(GOOS)

all: build
.PHONY: all

build:
	mkdir $(OUT_DIR); \
	cd cmd;		\
   	$(PLATFORM) $(GOBUILD) -o ../$(OUT_DIR)/$(BINARY_NAME) -v
.PHONY: build

update:
	glide up -v
.PHONY: update

fmt:
	gofmt -s -l -w ./cmd/
	gofmt -s -l -w ./pkg/

clean:
	$(GOCLEAN); \
	rm -rf $(OUT_DIR)
.PHONY: clean

