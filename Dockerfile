FROM golang:1.10-alpine as build

RUN apk update && apk add git make \
  && go get -d github.com/jmcfarlane/markdown \
  && cd $GOPATH/src/github.com/jmcfarlane/markdown \
  && make get \
  && make \
  && mv markdown /usr/local/bin/markdown

FROM golang:1.10-alpine

COPY --from=build /usr/local/bin/markdown /usr/local/bin/markdown

# could be made configurable via arg
WORKDIR /usr/src/app

ENTRYPOINT ["/usr/local/bin/markdown"]
