BUILD_DIR := build

build: build-elementary build-gameoflife build-gameoflifedx
.PHONY: build

build-%:
	go build -o $(BUILD_DIR)/$* ./cmd/$*

distclean:
	rm -rf build
