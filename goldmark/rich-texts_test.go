package goldmark

import (
	"testing"

	"github.com/faetools/go-notion/pkg/notion"
	n_ast "github.com/faetools/notion-to-goldmark/ast"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
)

func richTextWithColor(content string, bold, italic bool, color notion.Color) notion.RichText {
	rt := richText(content, bold, italic)
	rt.Annotations.Color = color

	return rt
}

func richText(content string, bold, italic bool) notion.RichText {
	rt := notion.NewRichText(content)
	rt.Annotations.Bold = bold
	rt.Annotations.Italic = italic

	return rt
}

func bold(nodes ...ast.Node) ast.Node {
	return wrap(ast.NewEmphasis(2), nodes...)
}

func italic(nodes ...ast.Node) ast.Node {
	return wrap(ast.NewEmphasis(1), nodes...)
}

func color(c notion.Color, nodes ...ast.Node) ast.Node {
	return wrap(&n_ast.Color{Color: c}, nodes...)
}

func wrap(n ast.Node, nodes ...ast.Node) ast.Node {
	for _, child := range nodes {
		n.AppendChild(n, child)
	}

	return n
}

const t1, t2, t3 = "foo", "bar", "blub"

func TestRichTexts(t *testing.T) {
	t.Parallel()

	assert.PanicsWithValue(t, `invalid RichText of type ""`,
		func() { _ = newAnnotationWrapper(notion.RichText{}) })

	for _, tt := range []struct {
		name string
		rts  notion.RichTexts
		res  []ast.Node
	}{
		{
			"n1->n1",
			notion.RichTexts{richText(t1, false, false)},
			[]ast.Node{newString(t1)},
		},
		{
			"n1,n2,n3->n1,n2,n3",
			notion.RichTexts{
				richText(t1, false, false),
				richText(t2, false, false),
				richText(t3, false, false),
			},
			[]ast.Node{newString(t1), newString(t2), newString(t3)},
		},
		{
			"n1(bold,italic)->italic{bold{n1}}",
			notion.RichTexts{
				richText(t1, true, true),
			},
			[]ast.Node{italic(bold(newString(t1)))},
		},
		{
			"n1(bold,italic),n2(bold),n3(bold,italic)->bold{italic{n1},n2,italic{n3}}",
			notion.RichTexts{
				richText(t1, true, true),
				richText(t2, true, false),
				richText(t3, true, true),
			},
			[]ast.Node{bold(
				italic(
					newString(t1),
				),
				newString(t2),
				italic(
					newString(t3),
				),
			)},
		},
		{
			"n1(bold,italic),n2(italic),n3(bold,italic)->italic{bold{n1},n2,bold{n3}}",
			notion.RichTexts{
				richText(t1, true, true),
				richText(t2, false, true),
				richText(t3, true, true),
			},
			[]ast.Node{italic(
				bold(
					newString(t1),
				),
				newString(t2),
				bold(
					newString(t3),
				),
			)},
		},
		{
			"n1(bold),n2(bold,italic),n3(italic)->bold{n1,italic{n2}},italic{n3}",
			notion.RichTexts{
				richText(t1, true, false),
				richText(t2, true, true),
				richText(t3, false, true),
			},
			[]ast.Node{
				bold(
					newString(t1),
					italic(
						newString(t2),
					),
				),
				italic(
					newString(t3),
				),
			},
		},
		{
			"n1,n2(bold)-> n1,bold{n2}",
			notion.RichTexts{
				richText(t1, false, false),
				richText(t2, true, false),
			},
			[]ast.Node{
				newString(t1),
				bold(newString(t2)),
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// ns := toNodeRichTexts(tt.rts)
			// assert.Equal(t, tt.res, ns)
		})

	}

	// we want to wrap the nodes
	// - n1(bold, italic), n2(italic), n3(bold, italic) -> italic{bold{n1}, n2, bold{n3}}
	// - n1(bold), n2(bold, italic), n3(italic) -> bold{n1, italic{n2}}, italic{n3}
	// - n1, n2(bold), -> n1, bold{n2}
}

// func TestRichTexts_FromExamplePage(t *testing.T) {
// 	t.Parallel()

// 	ns := toNodeRichTexts(notion.RichTexts{
// 		notion.NewRichText("This "),
// 		richTextWithColor("is", false, true, notion.ColorPurple),
// 		richTextWithColor(" a big ", false, false, notion.ColorPurple),
// 		richTextWithColor("heading", true, false, notion.ColorPurple),
// 		notion.NewRichText(" 1"),
// 	})

// 	assert.Equal(t, []ast.Node{
// 		newString("This "),
// 		color(notion.ColorPurple,
// 			italic(newString("is")),
// 			newString(" a big "),
// 			bold(newString("heading")),
// 		),
// 		newString(" 1"),
// 	}, ns)
// }
