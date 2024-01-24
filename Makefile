bootstrap ?= 5.3.2

# build: Compile locally (no dependencies are fetched)
build: tidy test static/bootstrap.min.css
	go build

clean:
	rm -f static/bootstrap.min.css
	rm -f markdown

# run: Compile and run
run: tidy test build
	./markdown

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

static/bootstrap.min.css:
	curl -s -o static/bootstrap.min.css \
		https://cdn.jsdelivr.net/npm/bootstrap@$(bootstrap)/dist/css/bootstrap.min.css
