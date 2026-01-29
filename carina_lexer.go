package nebel

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

func init() {
	lexers.Register(chroma.MustNewLexer(
		&chroma.Config{
			Name:      "Carina",
			Aliases:   []string{"carina", "crn"},
			Filenames: []string{"*.crn"},
			MimeTypes: []string{"text/x-carina"},
		},
		func() chroma.Rules {
			return chroma.Rules{
				"root": {
					// Comments
					{Pattern: `#.*$`, Type: chroma.Comment, Mutator: nil},

					// Keywords
					{Pattern: `\b(provider|let|backend|import|as|input|output)\b`, Type: chroma.Keyword, Mutator: nil},

					// Built-in types/functions
					{Pattern: `\b(ref|list|bool|cidr|string|number)\b`, Type: chroma.KeywordType, Mutator: nil},

					// Booleans
					{Pattern: `\b(true|false)\b`, Type: chroma.KeywordConstant, Mutator: nil},

					// Resource types (aws.s3.bucket, etc.)
					{Pattern: `(aws|gcp|azure)\.[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)*`, Type: chroma.NameClass, Mutator: nil},

					// Region constants (aws.Region.xxx)
					{Pattern: `aws\.Region\.[a-z][a-z0-9_]*`, Type: chroma.NameConstant, Mutator: nil},

					// Strings
					{Pattern: `"`, Type: chroma.StringDouble, Mutator: chroma.Push("string")},

					// Numbers
					{Pattern: `\b[0-9]+\b`, Type: chroma.NumberInteger, Mutator: nil},

					// Operators
					{Pattern: `=`, Type: chroma.Operator, Mutator: nil},

					// Punctuation
					{Pattern: `[{}()\[\],.]`, Type: chroma.Punctuation, Mutator: nil},

					// Property names (at the start of a line, before =)
					{Pattern: `^\s*([a-z][a-z0-9_]*)\s*(?==)`, Type: chroma.NameAttribute, Mutator: nil},

					// Identifiers
					{Pattern: `[a-z][a-z0-9_]*`, Type: chroma.NameVariable, Mutator: nil},

					// Whitespace
					{Pattern: `\s+`, Type: chroma.Text, Mutator: nil},
				},
				"string": {
					{Pattern: `\\.`, Type: chroma.StringEscape, Mutator: nil},
					{Pattern: `"`, Type: chroma.StringDouble, Mutator: chroma.Pop(1)},
					{Pattern: `[^"\\]+`, Type: chroma.StringDouble, Mutator: nil},
				},
			}
		},
	))
}
