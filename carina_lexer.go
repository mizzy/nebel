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

					// Keywords (storage, declaration, control, other) — see carina-core keywords.rs
					{Pattern: `\b(fn|let|arguments|attributes|backend|exports|moved|provider|removed|upstream_state|validation|else|for|if|in|import|read|require|use)\b`, Type: chroma.Keyword, Mutator: nil},

					// Null literal
					{Pattern: `\bnull\b`, Type: chroma.KeywordConstant, Mutator: nil},

					// Built-in types/functions
					{Pattern: `\b(ref|list|bool|cidr|string|number)\b`, Type: chroma.KeywordType, Mutator: nil},

					// Booleans
					{Pattern: `\b(true|false)\b`, Type: chroma.KeywordConstant, Mutator: nil},

					// Resource types (aws.s3.bucket, awscc.ec2.SecurityGroup, etc.)
					// Last segment can be PascalCase (new casing) or snake_case (legacy)
					{Pattern: `(aws|awscc|gcp|azure)\.[a-z][a-zA-Z0-9_]*(\.[a-zA-Z][a-zA-Z0-9_]*)*`, Type: chroma.NameClass, Mutator: nil},

					// Region constants (aws.Region.xxx)
					{Pattern: `aws\.Region\.[a-z][a-z0-9_]*`, Type: chroma.NameConstant, Mutator: nil},

					// Strings
					{Pattern: `"`, Type: chroma.StringDouble, Mutator: chroma.Push("string")},
					{Pattern: `'`, Type: chroma.StringSingle, Mutator: chroma.Push("sstring")},

					// Numbers
					{Pattern: `\b[0-9]+\b`, Type: chroma.NumberInteger, Mutator: nil},

					// Operators
					{Pattern: `=`, Type: chroma.Operator, Mutator: nil},

					// Punctuation
					{Pattern: `[{}()\[\],.]`, Type: chroma.Punctuation, Mutator: nil},

					// Property names (identifier directly followed by `=`)
					{Pattern: `[a-zA-Z][a-zA-Z0-9_]*(?=\s*=[^=])`, Type: chroma.NameAttribute, Mutator: nil},

					// Identifiers
					{Pattern: `[a-zA-Z][a-zA-Z0-9_]*`, Type: chroma.NameVariable, Mutator: nil},

					// Whitespace
					{Pattern: `\s+`, Type: chroma.Text, Mutator: nil},
				},
				"string": {
					{Pattern: `\\.`, Type: chroma.StringEscape, Mutator: nil},
					{Pattern: `"`, Type: chroma.StringDouble, Mutator: chroma.Pop(1)},
					{Pattern: `[^"\\]+`, Type: chroma.StringDouble, Mutator: nil},
				},
				"sstring": {
					{Pattern: `\\.`, Type: chroma.StringEscape, Mutator: nil},
					{Pattern: `'`, Type: chroma.StringSingle, Mutator: chroma.Pop(1)},
					{Pattern: `[^'\\]+`, Type: chroma.StringSingle, Mutator: nil},
				},
			}
		},
	))
}
