# protobuf-language-server

A language server implementation for Google Protocol Buffers

> [!CAUTION]
> This project is an experiment in building a golang language server in wasi (i.e. pure go in wasm without syscall/js).
>
> Please use the the upstream [vscode-extension](https://github.com/lasorda/protobuf-language-server) by [lasorda](https://github.com/lasorda) for general use.

## features

1. documentSymbol
2. jump to defines
3. format files with retab
4. code completion

