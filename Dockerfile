FROM golang:1.22-alpine as builder
# https://ebitengine.org/en/documents/install.html
RUN apk add alsa-lib-dev libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev mesa-dev pkgconf \
        git

RUN mkdir /build
COPY . /build/

RUN cd /build && \
    go mod tidy && \
    env GOOS=js GOARCH=wasm go build -o app.wasm ./main.go

# Use a non root, unprivileged nginx
# https://hub.docker.com/r/nginxinc/nginx-unprivileged
FROM nginxinc/nginx-unprivileged:1.25.4-alpine

COPY web/mime.types etc/nginx/mime.types
COPY web/index.html web/wasm_exec.js    /usr/share/nginx/html/

COPY --from=builder /build/app.wasm     /usr/share/nginx/html/app.wasm
