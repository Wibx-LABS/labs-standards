# syntax=docker/dockerfile:1
# Template Labs — Go. Imagem toolbox: sem CMD; o compose escolhe o papel (command: orchestrator | watchdog | migrate).
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app

FROM alpine:3
RUN adduser -D -H -u 10001 engine
COPY --link --from=builder /out/app /usr/local/bin/app
USER engine
ENTRYPOINT ["app"]
