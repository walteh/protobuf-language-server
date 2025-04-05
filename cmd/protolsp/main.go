package main

import (
	"flag"

	"github.com/walteh/protobuf-language-server/components"
	"github.com/walteh/protobuf-language-server/proto/view"

	"github.com/walteh/protobuf-language-server/go-lsp/logs"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp/defines"
)

var (
	address *string
	logPath *string
	stdio   *bool
)

func init() {
	logPath = flag.String("logs", logs.DefaultLogFilePath(), "logs file path")
	address = flag.String("listen", "", "address on which to listen for remote connections")
	stdio = flag.Bool("stdio", false, "")
}

func main() {
	flag.Parse()
	logs.Init(logPath)

	config := &lsp.Options{
		CompletionProvider: &defines.CompletionOptions{
			TriggerCharacters: &[]string{"."},
		},
	}
	if *address != "" {
		config.Address = *address
		config.Network = "tcp"
	}

	server := lsp.NewServer(config)

	view.Init(server)
	server.OnDocumentSymbolWithSliceDocumentSymbol(components.ProvideDocumentSymbol)
	server.OnDefinition(components.JumpDefine)
	server.OnDocumentFormatting(components.Format)
	server.OnCompletion(components.Completion)
	server.OnHover(components.Hover)
	server.OnDocumentRangeFormatting(components.FormatRange)
	server.Run()
}
