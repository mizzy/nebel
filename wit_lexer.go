package nebel

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

func init() {
	lexers.Register(chroma.MustNewLexer(
		&chroma.Config{
			Name:      "WIT",
			Aliases:   []string{"wit"},
			Filenames: []string{"*.wit"},
			MimeTypes: []string{"text/x-wit"},
		},
		func() chroma.Rules {
			return chroma.Rules{
				"root": {
					// Comments
					{Pattern: `//.*$`, Type: chroma.Comment, Mutator: nil},
					{Pattern: `/\*`, Type: chroma.CommentMultiline, Mutator: chroma.Push("comment")},

					// Keywords
					{Pattern: `\b(interface|world|import|export|use|func|type|record|variant|enum|flags|resource|package|include)\b`, Type: chroma.Keyword, Mutator: nil},

					// Built-in types
					{Pattern: `\b(bool|s8|s16|s32|s64|u8|u16|u32|u64|f32|f64|char|string|list|option|result|tuple)\b`, Type: chroma.KeywordType, Mutator: nil},

					// Booleans
					{Pattern: `\b(true|false)\b`, Type: chroma.KeywordConstant, Mutator: nil},

					// Strings
					{Pattern: `"`, Type: chroma.StringDouble, Mutator: chroma.Push("string")},

					// Numbers
					{Pattern: `\b[0-9]+\b`, Type: chroma.NumberInteger, Mutator: nil},

					// Operators and punctuation
					{Pattern: `->`, Type: chroma.Operator, Mutator: nil},
					{Pattern: `[{}()\[\],;:=<>.]`, Type: chroma.Punctuation, Mutator: nil},

					// Package paths (e.g., wasi:http/outgoing-handler@0.2.0)
					{Pattern: `[a-z][a-z0-9-]*:[a-z][a-z0-9-]*/[a-z][a-z0-9-]*(@[0-9.]+)?`, Type: chroma.NameNamespace, Mutator: nil},

					// Field/param names before colon
					{Pattern: `[a-z][a-z0-9-]*(?=\s*:)`, Type: chroma.NameAttribute, Mutator: nil},

					// Type names (PascalCase or kebab-case starting with lowercase)
					{Pattern: `_\b`, Type: chroma.NameBuiltin, Mutator: nil},
					{Pattern: `[a-z][a-z0-9-]*`, Type: chroma.NameVariable, Mutator: nil},

					// Whitespace
					{Pattern: `\s+`, Type: chroma.Text, Mutator: nil},
				},
				"string": {
					{Pattern: `\\.`, Type: chroma.StringEscape, Mutator: nil},
					{Pattern: `"`, Type: chroma.StringDouble, Mutator: chroma.Pop(1)},
					{Pattern: `[^"\\]+`, Type: chroma.StringDouble, Mutator: nil},
				},
				"comment": {
					{Pattern: `\*/`, Type: chroma.CommentMultiline, Mutator: chroma.Pop(1)},
					{Pattern: `[^*]+`, Type: chroma.CommentMultiline, Mutator: nil},
					{Pattern: `\*`, Type: chroma.CommentMultiline, Mutator: nil},
				},
			}
		},
	))
}
