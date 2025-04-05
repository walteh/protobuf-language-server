package components

import (
	"context"

	protobuf "github.com/emicklei/proto"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp/defines"
	"github.com/walteh/protobuf-language-server/proto/view"
)

// SemanticTokenTypes defines the token types we support
var SemanticTokenTypes = []string{
	"namespace",     // 0: package names
	"type",          // 1: built-in types (string, int32, etc)
	"class",         // 2: message types
	"enum",          // 3: enum types
	"interface",     // 4: service types
	"struct",        // 5: not used
	"typeParameter", // 6: not used
	"parameter",     // 7: field names
	"variable",      // 8: not used
	"property",      // 9: field options
	"enumMember",    // 10: enum values
	"function",      // 11: not used
	"method",        // 12: RPC methods
	"keyword",       // 13: protobuf keywords (message, service, rpc, etc)
	"string",        // 14: string literals
	"number",        // 15: numeric literals
	"operator",      // 16: = ( ) [ ] { }
	"comment",       // 17: comments
}

// SemanticTokenModifiers defines the token modifiers we support
var SemanticTokenModifiers = []string{
	"declaration",   // 0: type declarations
	"definition",    // 1: field definitions
	"readonly",      // 2: not used
	"static",        // 3: not used
	"deprecated",    // 4: deprecated fields/methods
	"documentation", // 5: documentation comments
}

// ProvideSemanticTokensLegend provides the legend for semantic tokens
func ProvideSemanticTokensLegend() *defines.SemanticTokensLegend {
	return &defines.SemanticTokensLegend{
		TokenTypes:     SemanticTokenTypes,
		TokenModifiers: SemanticTokenModifiers,
	}
}

// GetSemanticTokensOptions returns the semantic tokens options for server initialization
func GetSemanticTokensOptions() *defines.SemanticTokensOptions {
	return &defines.SemanticTokensOptions{
		Legend: *ProvideSemanticTokensLegend(),
		Full:   &[]bool{true}[0],
		Range:  &[]bool{false}[0],
	}
}

// ProvideSemanticTokens provides semantic tokens for a Protobuf file
func ProvideSemanticTokens(ctx context.Context, params *defines.SemanticTokensParams) (*defines.SemanticTokens, error) {
	doc, err := view.ViewManager.GetFile(params.TextDocument.Uri)
	if err != nil {
		return nil, nil
	}

	proto := doc.Proto()
	if proto == nil {
		return nil, nil
	}

	// Initialize the semantic tokens data array
	var data []uint

	// Helper function to add a token to the data array
	addToken := func(line, startChar, length uint32, tokenType, tokenModifiers uint32) {
		data = append(data, uint(line), uint(startChar), uint(length), uint(tokenType), uint(tokenModifiers))
	}

	// Process packages
	for _, pkg := range proto.Packages() {
		// Package keyword
		pos := pkg.ProtoPackage.Position
		addToken(uint32(pos.Line-1), uint32(pos.Column-1), 7, 13, 0) // "package" keyword
		// Package name
		addToken(uint32(pos.Line-1), uint32(pos.Column+7), uint32(len(pkg.ProtoPackage.Name)), 0, 0) // namespace
	}

	// Process messages
	for _, msg := range proto.Messages() {
		pos := msg.Protobuf().Position
		// Message keyword
		addToken(uint32(pos.Line-1), uint32(pos.Column-1), 7, 13, 0) // "message" keyword
		// Message name
		addToken(uint32(pos.Line-1), uint32(pos.Column+7), uint32(len(msg.Protobuf().Name)), 2, 0) // class

		// Process message fields
		for _, field := range msg.Fields() {
			fpos := field.ProtoField.Position
			// Field type
			addToken(uint32(fpos.Line-1), uint32(fpos.Column-1), uint32(len(field.ProtoField.Type)), 1, 0) // type
			// Field name
			addToken(uint32(fpos.Line-1), uint32(fpos.Column+len(field.ProtoField.Type)+1), uint32(len(field.ProtoField.Name)), 7, 1) // parameter

			// Field options
			for _, opt := range field.ProtoField.Options {
				optPos := opt.Position
				addToken(uint32(optPos.Line-1), uint32(optPos.Column-1), uint32(len(opt.Name)), 9, 0) // property
			}
		}
	}

	// Process enums
	for _, enum := range proto.Enums() {
		pos := enum.Protobuf().Position
		// Enum keyword
		addToken(uint32(pos.Line-1), uint32(pos.Column-1), 4, 13, 0) // "enum" keyword
		// Enum name
		addToken(uint32(pos.Line-1), uint32(pos.Column+4), uint32(len(enum.Protobuf().Name)), 3, 0) // enum

		// Process enum fields
		for _, element := range enum.Protobuf().Elements {
			if field, ok := element.(*protobuf.EnumField); ok {
				// Enum value name
				addToken(uint32(field.Position.Line-1), uint32(field.Position.Column-1), uint32(len(field.Name)), 10, 0) // enumMember
			}
		}
	}

	// Process services
	for _, svc := range proto.Services() {
		pos := svc.Protobuf().Position
		// Service keyword
		addToken(uint32(pos.Line-1), uint32(pos.Column-1), 7, 13, 0) // "service" keyword
		// Service name
		addToken(uint32(pos.Line-1), uint32(pos.Column+7), uint32(len(svc.Protobuf().Name)), 4, 0) // interface

		// Process RPC methods
		for _, rpc := range svc.RPCs() {
			rpos := rpc.ProtoRPC.Position
			// RPC keyword
			addToken(uint32(rpos.Line-1), uint32(rpos.Column-1), 3, 13, 0) // "rpc" keyword
			// Method name
			addToken(uint32(rpos.Line-1), uint32(rpos.Column+3), uint32(len(rpc.ProtoRPC.Name)), 12, 0) // method
		}
	}

	return &defines.SemanticTokens{
		Data: data,
	}, nil
}
