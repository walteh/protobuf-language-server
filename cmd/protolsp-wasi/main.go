//go:build wasip1

package main

import (
	"github.com/walteh/protobuf-language-server/components"
	"github.com/walteh/protobuf-language-server/go-lsp/logs"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp/defines"
	"github.com/walteh/protobuf-language-server/proto/view"
)

func main() {
	config := &lsp.Options{
		CompletionProvider: &defines.CompletionOptions{
			TriggerCharacters: &[]string{"."},
		},
	}

	server := lsp.NewServer(config)

	server.SetWasi(true)

	logs.Init(nil)
	view.Init(server)
	server.OnDocumentSymbolWithSliceDocumentSymbol(components.ProvideDocumentSymbol)
	server.OnDefinition(components.JumpDefine)
	server.OnDocumentFormatting(components.Format)
	server.OnCompletion(components.Completion)
	server.OnHover(components.Hover)
	server.OnDocumentRangeFormatting(components.FormatRange)
	server.Run()
}
