# base ci/cd image
FROM docker

RUN apk add --no-cache \
		ca-certificates \
        alpine-sdk

RUN apk add --no-cache --virtual .build-deps \
		bash \
		gcc \
		musl-dev \
		openssl \
		go \
        ;

RUN export \
        GOROOT_BOOTSTRAP="$(go env GOROOT)" \
        GOOS="$(go env GOOS)" \
        GOARCH="$(go env GOARCH)" \
        GOHOSTOS="$(go env GOHOSTOS)" \
        GOHOSTARCH="$(go env GOHOSTARCH)" \
        PATH="/usr/local/go/bin:$PATH";

RUN mkdir -p "go/src" "go/bin" && chmod -R 777 "go"
RUN export GOPATH="go"
