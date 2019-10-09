BUILD_DIR := build

build: build-elementary build-gameoflife build-gameoflifedx
.PHONY: build

build-%:
	go build -o $(BUILD_DIR)/$* ./cmd/$*

distclean:
	rm -rf build

lint:
	golangci-lint run
.PHONY: lint

test:
	go test  ./...
.PHONY: test

coverage:
	go test -cover ./...
.PHONY: coverage
