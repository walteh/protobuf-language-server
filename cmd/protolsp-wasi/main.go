//go:build wasip1

package main

import (
	"github.com/walteh/protobuf-language-server/components"
	"github.com/walteh/protobuf-language-server/go-lsp/logs"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp/defines"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp/tokens"
	"github.com/walteh/protobuf-language-server/proto/view"
)

func main() {

	config := &lsp.Options{
		CompletionProvider: &defines.CompletionOptions{
			TriggerCharacters: &[]string{"."},
		},
		SemanticTokensProvider: tokens.GetSemanticTokensOptions(),
	}

	server := lsp.NewServer(config)

	logs.Init(nil)
	view.Init(server)
	server.OnDocumentSymbolWithSliceDocumentSymbol(components.ProvideDocumentSymbol)
	server.OnDefinition(components.JumpDefine)
	server.OnDocumentFormatting(components.FormatWithRetab)
	server.OnCompletion(components.Completion)
	server.OnHover(components.Hover)
	server.OnDocumentRangeFormatting(components.FormatRange)
	server.OnSemanticTokens(components.ProvideSemanticTokens)

	server.Run()
}
