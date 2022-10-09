FROM golang:1.10-alpine

ARG VERSION="development"

ENV GO_PATH="/go"

ADD . $GO_PATH/src/github.com/Thijn/drone-sentry
WORKDIR $GO_PATH/src/github.com/Thijn/drone-sentry

RUN apk add --update ca-certificates git
RUN go get -t ./...
RUN go test -v ./...

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o bin/drone-sentry -a -installsuffix cgo -ldflags "-s -X main.version=$VERSION"

FROM alpine:latest
RUN apk add --update ca-certificates
LABEL maintainer="Sierra Softworks <admin@sierrasoftworks.com>"

COPY --from=0 /go/src/github.com/Thijn/drone-sentry/bin/drone-sentry /bin/drone-sentry

LABEL VERSION=$VERSION

WORKDIR /bin
ENTRYPOINT ["/bin/drone-sentry"]
