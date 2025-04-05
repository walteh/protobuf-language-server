# protobuf-language-server

A language server implementation for Google Protocol Buffers

## how to use

1. Get code from [https://github.com/walteh/protobuf-language-server](https://github.com/walteh/protobuf-language-server)
2. Build the target `protobuf-language-server`, add `protobuf-language-server` to `PATH`

## features

1. documentSymbol
2. jump to defines
3. format file with clang-format
4. code completion

## build vscode extension(optional for deveplop)

```shell
npm install -g vsce
npm install -g yarn
npm install
vsce package --no-yarn
```
