# syntax = docker/dockerfile:1.2

FROM golang:1.20.0-alpine AS build
WORKDIR /src
RUN apk add --no-cache file git
ENV GOMODCACHE /root/.cache/gocache
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
    CGO_ENABLED=0 go build -o /out/mockingio -ldflags '-s -d -w' ./; \
    file /out/mockingio | grep "statically linked"

FROM scratch
COPY --from=build /out/mockingio /bin/mockingio
ENTRYPOINT ["/bin/mockingio"]
