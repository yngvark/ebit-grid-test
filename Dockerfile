FROM golang:1.22-alpine as builder

WORKDIR /build

RUN apk add alsa-lib-dev libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev mesa-dev pkgconf \
        git

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" GOOS=js GOARCH=wasm go build -o app.wasm ./main.go

FROM caddy:2-alpine

COPY web/index.html web/wasm_exec.js    /usr/share/caddy/
COPY --from=builder /build/app.wasm     /usr/share/caddy/

VOLUME /config /data
