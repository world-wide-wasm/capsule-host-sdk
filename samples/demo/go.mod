module demo

go 1.22.1

require github.com/world-wide-wasm/capsule-host-sdk v0.0.0-20240331063815-0662baa6b24a

require github.com/tetratelabs/wazero v1.7.0 // indirect

replace github.com/world-wide-wasm/capsule-host-sdk => ../..
