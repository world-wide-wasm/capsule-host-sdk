# Capsule host-sdk
> 🚧 This is a work in progress

Capsule Host SDK (HDK) is a WebAssembly SDK, based on Wazero, to build WebAssembly host applications with Golang. The WebAssembly modules used by the host are built with TinyGo.


## Docker 
> - Use this project with DevContainer
> - ✋ this image works only on a arm architecture
> - 🚧 TODO: create a multiarch image

Build the workspace image:
```bash
docker compose build
```

Push the workspace image to the Docker hub:
```bash
docker compose push
```
