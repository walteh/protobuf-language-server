package tokens

import (
	"github.com/walteh/protobuf-language-server/go-lsp/lsp/defines"
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
