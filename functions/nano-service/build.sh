#!/bin/bash
tinygo build -scheduler=none --no-debug \
  -o nano-service.wasm \
  -target wasi main.go

ls -lh *.wasm
