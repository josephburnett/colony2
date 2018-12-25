#!/usr/bin/env bash
WASM_FILES="$(go env GOROOT)/misc/wasm"
GOOS=js GOARCH=wasm go build -o resources/build/client.wasm ./cmd/client.go
cp ${WASM_FILES}/wasm_exec.js ./resources/build
go run ./cmd/server.go
