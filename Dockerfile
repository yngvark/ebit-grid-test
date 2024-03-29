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

# Use a non root, unprivileged nginx
# https://hub.docker.com/r/nginxinc/nginx-unprivileged
FROM nginxinc/nginx-unprivileged:1.25.4-alpine

COPY web/mime.types etc/nginx/mime.types
COPY web/index.html web/wasm_exec.js    /usr/share/nginx/html/

COPY --from=builder /build/app.wasm     /usr/share/nginx/html/app.wasm
