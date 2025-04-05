package components

import (
	"context"
	"reflect"
	"sort"
	"strconv"

	protobuf "github.com/emicklei/proto"
	"github.com/walteh/protobuf-language-server/go-lsp/lsp/defines"
	"github.com/walteh/protobuf-language-server/proto/parser"
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
		Range:  &[]bool{true}[0],
	}
}

// Token represents a semantic token with absolute position
type Token struct {
	Line          int
	StartChar     int
	Length        int
	TokenType     int
	TokenModifier int
}

// Helper function to check if a value is nil
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	return (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && v.IsNil()
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

	// Collect tokens with absolute positions first
	var tokens []Token

	// Process syntax (if available in the proto.Protobuf())
	processSyntax(proto.Protobuf(), &tokens)

	// Process package declarations
	processPackages(proto, &tokens)

	// Process imports
	processImports(proto, &tokens)

	// Process messages
	processMessages(proto, &tokens)

	// Process enums
	processEnums(proto, &tokens)

	// Process services
	processServices(proto, &tokens)

	// Sort tokens by position
	sort.Slice(tokens, func(i, j int) bool {
		if tokens[i].Line != tokens[j].Line {
			return tokens[i].Line < tokens[j].Line
		}
		return tokens[i].StartChar < tokens[j].StartChar
	})

	// Convert to relative positions as required by LSP
	return encodeTokens(tokens), nil
}

// processSyntax handles the syntax statement
func processSyntax(proto *protobuf.Proto, tokens *[]Token) {
	for _, element := range proto.Elements {
		if syntax, ok := element.(*protobuf.Syntax); ok {
			pos := syntax.Position

			// Syntax keyword
			*tokens = append(*tokens, Token{
				Line:          pos.Line - 1,
				StartChar:     pos.Column - 1,
				Length:        6,  // "syntax"
				TokenType:     13, // keyword
				TokenModifier: 0,
			})

			// Equal sign
			*tokens = append(*tokens, Token{
				Line:          pos.Line - 1,
				StartChar:     pos.Column + 6, // After "syntax"
				Length:        1,              // "="
				TokenType:     16,             // operator
				TokenModifier: 0,
			})

			// Quote and value
			valuePos := pos.Column + 8 // After 'syntax = '
			*tokens = append(*tokens, Token{
				Line:          pos.Line - 1,
				StartChar:     valuePos,
				Length:        len(syntax.Value) + 2, // Include quotes
				TokenType:     14,                    // string
				TokenModifier: 0,
			})
		}
	}
}

// processPackages handles package declarations
func processPackages(proto parser.Proto, tokens *[]Token) {
	for _, pkg := range proto.Packages() {
		pos := pkg.ProtoPackage.Position

		// Package keyword
		*tokens = append(*tokens, Token{
			Line:          pos.Line - 1,
			StartChar:     pos.Column - 1,
			Length:        7,  // "package"
			TokenType:     13, // keyword
			TokenModifier: 0,
		})

		// Package name
		nameStartChar := pos.Column + 7 // After "package "
		*tokens = append(*tokens, Token{
			Line:          pos.Line - 1,
			StartChar:     nameStartChar,
			Length:        len(pkg.ProtoPackage.Name),
			TokenType:     0, // namespace
			TokenModifier: 0,
		})
	}
}

// processImports handles import statements
func processImports(proto parser.Proto, tokens *[]Token) {
	for _, imp := range proto.Imports() {
		pos := imp.ProtoImport.Position

		// Import keyword
		*tokens = append(*tokens, Token{
			Line:          pos.Line - 1,
			StartChar:     pos.Column - 1,
			Length:        6,  // "import"
			TokenType:     13, // keyword
			TokenModifier: 0,
		})

		// Import path (string)
		pathStartChar := pos.Column + 6 // After "import "
		*tokens = append(*tokens, Token{
			Line:          pos.Line - 1,
			StartChar:     pathStartChar,
			Length:        len(imp.ProtoImport.Filename) + 2, // Include quotes
			TokenType:     14,                                // string
			TokenModifier: 0,
		})
	}
}

// processMessages handles message declarations and their fields
func processMessages(proto parser.Proto, tokens *[]Token) {
	for _, msg := range proto.Messages() {
		pos := msg.Protobuf().Position

		// Message keyword
		*tokens = append(*tokens, Token{
			Line:          pos.Line - 1,
			StartChar:     pos.Column - 1,
			Length:        7,  // "message"
			TokenType:     13, // keyword
			TokenModifier: 0,
		})

		// Message name
		nameStartChar := pos.Column + 7 // After "message "
		*tokens = append(*tokens, Token{
			Line:          pos.Line - 1,
			StartChar:     nameStartChar,
			Length:        len(msg.Protobuf().Name),
			TokenType:     2, // class
			TokenModifier: 0,
		})

		// Process message fields
		for _, field := range msg.Fields() {
			fpos := field.ProtoField.Position

			// Field type
			*tokens = append(*tokens, Token{
				Line:          fpos.Line - 1,
				StartChar:     fpos.Column - 1,
				Length:        len(field.ProtoField.Type),
				TokenType:     1, // type
				TokenModifier: 0,
			})

			// Field name
			nameStartChar := fpos.Column - 1 + len(field.ProtoField.Type) + 1 // After type and space
			*tokens = append(*tokens, Token{
				Line:          fpos.Line - 1,
				StartChar:     nameStartChar,
				Length:        len(field.ProtoField.Name),
				TokenType:     7, // parameter
				TokenModifier: 1, // definition
			})

			// Field options
			for _, opt := range field.ProtoField.Options {
				optPos := opt.Position
				*tokens = append(*tokens, Token{
					Line:          optPos.Line - 1,
					StartChar:     optPos.Column - 1,
					Length:        len(opt.Name),
					TokenType:     9, // property
					TokenModifier: 0,
				})

				// Option value
				valuePos := optPos.Column - 1 + len(opt.Name) + 1 // After name and '='
				valueLength := 1                                  // Minimum length
				if !isNil(opt.Constant) && opt.Constant.Source != "" {
					valueLength = len(opt.Constant.Source)
				}
				*tokens = append(*tokens, Token{
					Line:          optPos.Line - 1,
					StartChar:     valuePos,
					Length:        valueLength,
					TokenType:     14, // string or number depending on type
					TokenModifier: 0,
				})
			}
		}
	}
}

// processEnums handles enum declarations and their values
func processEnums(proto parser.Proto, tokens *[]Token) {
	for _, enum := range proto.Enums() {
		pos := enum.Protobuf().Position

		// Enum keyword
		*tokens = append(*tokens, Token{
			Line:          pos.Line - 1,
			StartChar:     pos.Column - 1,
			Length:        4,  // "enum"
			TokenType:     13, // keyword
			TokenModifier: 0,
		})

		// Enum name
		nameStartChar := pos.Column + 4 // After "enum "
		*tokens = append(*tokens, Token{
			Line:          pos.Line - 1,
			StartChar:     nameStartChar,
			Length:        len(enum.Protobuf().Name),
			TokenType:     3, // enum
			TokenModifier: 0,
		})

		// Process enum fields
		for _, element := range enum.Protobuf().Elements {
			if field, ok := element.(*protobuf.EnumField); ok {
				// Enum value name
				*tokens = append(*tokens, Token{
					Line:          field.Position.Line - 1,
					StartChar:     field.Position.Column - 1,
					Length:        len(field.Name),
					TokenType:     10, // enumMember
					TokenModifier: 0,
				})

				// Enum value
				valueStartChar := field.Position.Column - 1 + len(field.Name) + 2 // After name and " = "
				valueStr := strconv.Itoa(field.Integer)
				*tokens = append(*tokens, Token{
					Line:          field.Position.Line - 1,
					StartChar:     valueStartChar,
					Length:        len(valueStr),
					TokenType:     15, // number
					TokenModifier: 0,
				})
			}
		}
	}
}

// processServices handles service declarations and their RPCs
func processServices(proto parser.Proto, tokens *[]Token) {
	for _, svc := range proto.Services() {
		pos := svc.Protobuf().Position

		// Service keyword
		*tokens = append(*tokens, Token{
			Line:          pos.Line - 1,
			StartChar:     pos.Column - 1,
			Length:        7,  // "service"
			TokenType:     13, // keyword
			TokenModifier: 0,
		})

		// Service name
		nameStartChar := pos.Column + 7 // After "service "
		*tokens = append(*tokens, Token{
			Line:          pos.Line - 1,
			StartChar:     nameStartChar,
			Length:        len(svc.Protobuf().Name),
			TokenType:     4, // interface
			TokenModifier: 0,
		})

		// Process RPC methods
		for _, rpc := range svc.RPCs() {
			rpos := rpc.ProtoRPC.Position

			// RPC keyword
			*tokens = append(*tokens, Token{
				Line:          rpos.Line - 1,
				StartChar:     rpos.Column - 1,
				Length:        3,  // "rpc"
				TokenType:     13, // keyword
				TokenModifier: 0,
			})

			// Method name
			nameStartChar := rpos.Column + 3 // After "rpc "
			*tokens = append(*tokens, Token{
				Line:          rpos.Line - 1,
				StartChar:     nameStartChar,
				Length:        len(rpc.ProtoRPC.Name),
				TokenType:     12, // method
				TokenModifier: 0,
			})

			// Request type (if possible to determine)
			if rpc.ProtoRPC.RequestType != "" {
				reqPos := nameStartChar + len(rpc.ProtoRPC.Name) + 1 // After name and '('
				reqTypeLength := len(rpc.ProtoRPC.RequestType)
				if reqTypeLength > 0 {
					*tokens = append(*tokens, Token{
						Line:          rpos.Line - 1,
						StartChar:     reqPos,
						Length:        reqTypeLength,
						TokenType:     1, // type
						TokenModifier: 0,
					})
				}
			}

			// Response type (if possible to determine)
			if rpc.ProtoRPC.ReturnsType != "" {
				// This is an approximation as exact position isn't available
				returnTypeLength := len(rpc.ProtoRPC.ReturnsType)
				if returnTypeLength > 0 {
					*tokens = append(*tokens, Token{
						Line:          rpos.Line - 1,
						StartChar:     rpos.Column + 20, // Rough estimate after "rpc Name (Request) returns ("
						Length:        returnTypeLength,
						TokenType:     1, // type
						TokenModifier: 0,
					})
				}
			}
		}
	}
}

// encodeTokens converts absolute token positions to relative positions
// as required by the LSP spec
func encodeTokens(tokens []Token) *defines.SemanticTokens {
	if len(tokens) == 0 {
		return &defines.SemanticTokens{
			Data: []uint{},
		}
	}

	data := make([]uint, 0, len(tokens)*5)

	prevLine := tokens[0].Line
	prevChar := tokens[0].StartChar

	for _, token := range tokens {
		// Calculate delta line and delta start
		deltaLine := token.Line - prevLine
		deltaStart := token.StartChar
		if deltaLine == 0 {
			deltaStart = token.StartChar - prevChar
		}

		// The 5 values per token: deltaLine, deltaStart, length, tokenType, tokenModifiers
		data = append(data,
			uint(deltaLine),
			uint(deltaStart),
			uint(token.Length),
			uint(token.TokenType),
			uint(token.TokenModifier),
		)

		// Update previous position
		prevLine = token.Line
		prevChar = token.StartChar
	}

	return &defines.SemanticTokens{
		Data: data,
	}
}
