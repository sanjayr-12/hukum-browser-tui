package ui

import (
	"github.com/charmbracelet/glamour/ansi"
)

func stringPtr(s string) *string { return &s }
func boolPtr(b bool) *bool       { return &b }
func uintPtr(u uint) *uint       { return &u }

// OneDarkStyleConfig provides authentic Atom One Dark terminal typography
// and Chroma syntax highlighting without intrusive background bars or raw markdown prefixes.
var OneDarkStyleConfig = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockPrefix: "\n",
			BlockSuffix: "\n",
			Color:       stringPtr("#ABB2BF"),
		},
		Margin: uintPtr(1),
	},
	BlockQuote: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:  stringPtr("#5C6370"),
			Italic: boolPtr(true),
		},
		Indent:      uintPtr(1),
		IndentToken: stringPtr("│ "),
	},
	Paragraph: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: stringPtr("#ABB2BF"),
		},
	},
	List: ansi.StyleList{
		LevelIndent: 2,
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockPrefix: "\n",
			BlockSuffix: "\n",
			Bold:        boolPtr(true),
		},
	},
	H1: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: stringPtr("#ECEFF4"),
			Bold:  boolPtr(true),
		},
	},
	H2: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: stringPtr("#61AFEF"),
			Bold:  boolPtr(true),
		},
	},
	H3: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: stringPtr("#56B6C2"),
			Bold:  boolPtr(true),
		},
	},
	H4: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: stringPtr("#E5C07B"),
			Bold:  boolPtr(true),
		},
	},
	H5: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: stringPtr("#C678DD"),
			Bold:  boolPtr(true),
		},
	},
	H6: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: stringPtr("#5C6370"),
			Bold:  boolPtr(false),
		},
	},
	Text: ansi.StylePrimitive{
		Color: stringPtr("#ABB2BF"),
	},
	Strikethrough: ansi.StylePrimitive{
		CrossedOut: boolPtr(true),
	},
	Emph: ansi.StylePrimitive{
		Italic: boolPtr(true),
	},
	Strong: ansi.StylePrimitive{
		Bold:  boolPtr(true),
		Color: stringPtr("#ECEFF4"),
	},
	HorizontalRule: ansi.StylePrimitive{
		Color:  stringPtr("#3E4451"),
		Format: "\n────────────────────────────────────────\n",
	},
	Item: ansi.StylePrimitive{
		BlockPrefix: "• ",
		Color:       stringPtr("#61AFEF"),
	},
	Enumeration: ansi.StylePrimitive{
		BlockPrefix: ". ",
		Color:       stringPtr("#61AFEF"),
	},
	Task: ansi.StyleTask{
		Ticked:   "[✓] ",
		Unticked: "[ ] ",
	},
	Link: ansi.StylePrimitive{
		Color:     stringPtr("#61AFEF"),
		Underline: boolPtr(true),
	},
	LinkText: ansi.StylePrimitive{
		Color: stringPtr("#56B6C2"),
		Bold:  boolPtr(true),
	},
	Image: ansi.StylePrimitive{
		Color:     stringPtr("#61AFEF"),
		Underline: boolPtr(true),
	},
	ImageText: ansi.StylePrimitive{
		Color:  stringPtr("#56B6C2"),
		Format: "Image: {{.text}}",
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:           stringPtr("#E5C07B"),
			BackgroundColor: stringPtr("#21252B"),
			Prefix:          " ",
			Suffix:          " ",
		},
	},
	CodeBlock: ansi.StyleCodeBlock{
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: stringPtr("#ABB2BF"),
			},
			Margin: uintPtr(2),
		},
		Chroma: &ansi.Chroma{
			Text: ansi.StylePrimitive{
				Color: stringPtr("#ABB2BF"),
			},
			Error: ansi.StylePrimitive{
				Color: stringPtr("#E06C75"),
			},
			Comment: ansi.StylePrimitive{
				Color:  stringPtr("#5C6370"),
				Italic: boolPtr(true),
			},
			CommentPreproc: ansi.StylePrimitive{
				Color: stringPtr("#E5C07B"),
			},
			Keyword: ansi.StylePrimitive{
				Color: stringPtr("#C678DD"),
			},
			KeywordReserved: ansi.StylePrimitive{
				Color: stringPtr("#C678DD"),
			},
			KeywordNamespace: ansi.StylePrimitive{
				Color: stringPtr("#C678DD"),
			},
			KeywordType: ansi.StylePrimitive{
				Color: stringPtr("#E5C07B"),
			},
			Operator: ansi.StylePrimitive{
				Color: stringPtr("#56B6C2"),
			},
			Punctuation: ansi.StylePrimitive{
				Color: stringPtr("#ABB2BF"),
			},
			Name: ansi.StylePrimitive{
				Color: stringPtr("#61AFEF"),
			},
			NameConstant: ansi.StylePrimitive{
				Color: stringPtr("#D19A66"),
			},
			NameBuiltin: ansi.StylePrimitive{
				Color: stringPtr("#E5C07B"),
			},
			NameTag: ansi.StylePrimitive{
				Color: stringPtr("#E06C75"),
			},
			NameAttribute: ansi.StylePrimitive{
				Color: stringPtr("#D19A66"),
			},
			NameClass: ansi.StylePrimitive{
				Color: stringPtr("#E5C07B"),
				Bold:  boolPtr(true),
			},
			NameDecorator: ansi.StylePrimitive{
				Color: stringPtr("#61AFEF"),
			},
			NameFunction: ansi.StylePrimitive{
				Color: stringPtr("#61AFEF"),
			},
			LiteralNumber: ansi.StylePrimitive{
				Color: stringPtr("#D19A66"),
			},
			LiteralString: ansi.StylePrimitive{
				Color: stringPtr("#98C379"),
			},
			LiteralStringEscape: ansi.StylePrimitive{
				Color: stringPtr("#56B6C2"),
			},
			GenericDeleted: ansi.StylePrimitive{
				Color: stringPtr("#E06C75"),
			},
			GenericEmph: ansi.StylePrimitive{
				Italic: boolPtr(true),
			},
			GenericInserted: ansi.StylePrimitive{
				Color: stringPtr("#98C379"),
			},
			GenericStrong: ansi.StylePrimitive{
				Bold: boolPtr(true),
			},
			GenericSubheading: ansi.StylePrimitive{
				Color: stringPtr("#56B6C2"),
			},
			Background: ansi.StylePrimitive{
				BackgroundColor: stringPtr("#21252B"),
			},
		},
	},
	Table: ansi.StyleTable{
		CenterSeparator: stringPtr("┼"),
		ColumnSeparator: stringPtr("│"),
		RowSeparator:    stringPtr("─"),
	},
	DefinitionDescription: ansi.StylePrimitive{
		BlockPrefix: "\n• ",
	},
}
