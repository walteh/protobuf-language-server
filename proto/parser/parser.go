package parser

import (
	"io"

	protobuf "github.com/emicklei/proto"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp/defines"
)

// ParseProtos parses protobuf files from filenames and return parser.ProtoSet.
func ParseProto(document_uri defines.DocumentUri, r io.Reader) (Proto, error) {
	parser := protobuf.NewParser(r)
	p, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return NewProto(document_uri, p), nil
}
