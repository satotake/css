package css_test

import (
	"bytes"
	"testing"

	"github.com/benbjohnson/css"
)

// Ensure than the printer prints nodes correctly.
func TestPrinter_Print(t *testing.T) {
	var tests = []struct {
		in css.Node
		s  string
	}{
		// 0. Full stylesheet with multiple rules.
		{in: &css.StyleSheet{
			Rules: []css.Rule{
				&css.QualifiedRule{
					Prelude: []css.ComponentValue{
						&css.Token{Tok: css.IdentToken, Value: "foo"},
						&css.Token{Tok: css.WhitespaceToken, Value: " "},
						&css.Token{Tok: css.IdentToken, Value: "bar"},
					},
					Block: &css.SimpleBlock{
						Token: &css.Token{Tok: css.LBraceToken},
						Values: []css.ComponentValue{
							&css.Token{Tok: css.IdentToken, Value: "font-size"},
							&css.Token{Tok: css.ColonToken},
							&css.Token{Tok: css.IdentToken, Value: "10px"},
						},
					},
				},
				&css.AtRule{
					Name: "baz",
					Prelude: []css.ComponentValue{
						&css.Token{Tok: css.WhitespaceToken, Value: " "},
						&css.Token{Tok: css.IdentToken, Value: "my-rule"},
					},
				},
				&css.AtRule{
					Name: "qux",
					Prelude: []css.ComponentValue{
						&css.Token{Tok: css.WhitespaceToken, Value: " "},
						&css.Token{Tok: css.IdentToken, Value: "your-rule"},
					},
					Block: &css.DeclarationBlock{
						Token: &css.Token{Tok: css.LBraceToken},
						Declarations: css.Declarations{
							&css.Declaration{
								Name: "font-size",
								Values: []css.ComponentValue{
									&css.Token{Tok: css.WhitespaceToken},
									&css.Token{Tok: css.DimensionToken, Value: "10px"},
								},
								Important: true,
							},
						},
					},
				},
			},
		}, s: `foo bar{font-size:10px} @baz my-rule; @qux your-rule{font-size:10px!important;}`},

		// Test that nil values are safe to print.
		{in: (*css.StyleSheet)(nil), s: ``},       // 1
		{in: (css.Rules)(nil), s: ``},             // 2
		{in: (*css.AtRule)(nil), s: ``},           // 3
		{in: (*css.QualifiedRule)(nil), s: ``},    // 4
		{in: (css.Declarations)(nil), s: ``},      // 5
		{in: (*css.Declaration)(nil), s: ``},      // 6
		{in: (css.ComponentValues)(nil), s: ``},   // 7
		{in: (*css.SimpleBlock)(nil), s: ``},      // 8
		{in: (*css.DeclarationBlock)(nil), s: ``}, // 9
		{in: (*css.Function)(nil), s: ``},         // 10
		{in: (*css.Token)(nil), s: ``},            // 11

		// Test individual tokens.
		{in: &css.Token{Tok: css.IdentToken, Value: "foo"}, s: `foo`},                  // 12
		{in: &css.Token{Tok: css.FunctionToken, Value: "foo"}, s: `foo(`},              // 12
		{in: &css.Token{Tok: css.AtKeywordToken, Value: "☃"}, s: `@☃`},                 // 12
		{in: &css.Token{Tok: css.HashToken, Value: "foo"}, s: `#foo`},                  // 12
		{in: &css.Token{Tok: css.StringToken, Value: "foo", Ending: '"'}, s: `"foo"`},  // 12
		{in: &css.Token{Tok: css.StringToken, Value: "foo", Ending: '\''}, s: `'foo'`}, // 12
		{in: &css.Token{Tok: css.BadStringToken}, s: `''`},                             // 12
		{in: &css.Token{Tok: css.URLToken, Value: "foo"}, s: `url(foo)`},               // 12
		{in: &css.Token{Tok: css.BadURLToken, Value: "foo"}, s: `url()`},               // 12
		{in: &css.Token{Tok: css.DelimToken, Value: "."}, s: `.`},                      // 12
		{in: &css.Token{Tok: css.NumberToken, Value: "-20.3E2"}, s: `-20.3E2`},         // 12
		{in: &css.Token{Tok: css.PercentageToken, Value: "100%"}, s: `100%`},           // 12
		{in: &css.Token{Tok: css.DimensionToken, Value: "10cm"}, s: `10cm`},            // 12
		{in: &css.Token{Tok: css.WhitespaceToken, Value: "  "}, s: `  `},               // 12
		{in: &css.Token{Tok: css.DelimToken, Value: "."}, s: `.`},                      // 12
		{in: &css.Token{Tok: css.IncludeMatchToken}, s: `~=`},                          // 12
		{in: &css.Token{Tok: css.DashMatchToken}, s: `|=`},                             // 12
		{in: &css.Token{Tok: css.PrefixMatchToken}, s: `^=`},                           // 12
		{in: &css.Token{Tok: css.SuffixMatchToken}, s: `$=`},                           // 12
		{in: &css.Token{Tok: css.SubstringMatchToken}, s: `*=`},                        // 12
		{in: &css.Token{Tok: css.ColumnToken}, s: `||`},                                // 12
		{in: &css.Token{Tok: css.CDOToken}, s: `<!--`},                                 // 12
		{in: &css.Token{Tok: css.CDCToken}, s: `-->`},                                  // 12
		{in: &css.Token{Tok: css.ColonToken}, s: `:`},                                  // 12
		{in: &css.Token{Tok: css.SemicolonToken}, s: `;`},                              // 12
		{in: &css.Token{Tok: css.CommaToken}, s: `,`},                                  // 12
		{in: &css.Token{Tok: css.LBrackToken}, s: `[`},                                 // 12
		{in: &css.Token{Tok: css.RBrackToken}, s: `]`},                                 // 12
		{in: &css.Token{Tok: css.LParenToken}, s: `(`},                                 // 12
		{in: &css.Token{Tok: css.RParenToken}, s: `)`},                                 // 12
		{in: &css.Token{Tok: css.LBraceToken}, s: `{`},                                 // 12
		{in: &css.Token{Tok: css.RBraceToken}, s: `}`},                                 // 12

		{in: &css.Token{Tok: css.UnicodeRangeToken, Start: 10, End: 10}, s: `U+00000a`},          // 12
		{in: &css.Token{Tok: css.UnicodeRangeToken, Start: 10, End: 20}, s: `U+00000a-U+000014`}, // 12

		{in: &css.Token{Tok: css.EOFToken}, s: `EOF`}, // 12
	}

	for i, tt := range tests {
		var buf bytes.Buffer
		var p css.Printer
		err := p.Print(&buf, tt.in)

		if err != nil {
			t.Errorf("%d. unexpected error: %s", i, tt.s)
		} else if tt.s != buf.String() {
			t.Errorf("%d. \n\nexp: %s\n\ngot: %s\n\n", i, tt.s, buf.String())
		}
	}
}

// TODO(benbjohnson): Example: Printer.Print()
