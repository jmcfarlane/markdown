# build: Compile locally (no dependencies are fetched)
build: test
	go build

# help: This message ;)
help:
	@echo "Features this build supports:"
	@grep -E '^# [-a-z./]+:' Makefile | sed -e 's|#| > make|g' | sort

# install: Install into $GOPATH
install: get test
	go install

# test: Run tests (currently there are none)
test:
	go test -v
