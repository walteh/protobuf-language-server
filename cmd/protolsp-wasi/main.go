//go:build wasip1

package main

import (
	"github.com/walteh/protobuf-language-server/components"
	"github.com/walteh/protobuf-language-server/go-lsp/logs"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp/defines"
	"github.com/walteh/protobuf-language-server/proto/view"
)

// note about logging: the logging will go to the extenstion host process, so if we want to actually return logs to the extension,
// we need to do some extra work
// - like intercepting the logs and returning them to the extension in the response

func main() {
	config := &lsp.Options{
		CompletionProvider: &defines.CompletionOptions{
			TriggerCharacters: &[]string{"."},
		},
	}
	// if *address != "" {
	// 	config.Address = *address
	// 	config.Network = "tcp"
	// }

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

	// Keep the program running
	select {}
}
