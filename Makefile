# build: Compile locally (no dependencies are fetched)
build: tidy test
	go build

# tidy: Run go mod tidy
tidy:
	go mod tidy

# help: This message ;)
help:
	@echo "Features this build supports:"
	@grep -E '^# [-a-z./]+:' Makefile | sed -e 's|#| > make|g' | sort

# install: Install into $GOPATH
install: build
	go install

# test: Run tests (currently there are none)
test:
	go test -v
