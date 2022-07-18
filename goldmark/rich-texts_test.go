package goldmark

import (
	"encoding/json"
	"testing"

	"github.com/faetools/go-notion/pkg/notion"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
)

func richText(content string, bold, italic bool) notion.RichText {
	return notion.RichText{
		Type: notion.RichTextTypeText,
		Text: &notion.Text{Content: content},
		Annotations: notion.Annotations{
			Bold:   bold,
			Italic: italic,
			Color:  notion.ColorDefault,
		},
	}
}

func bold(nodes ...ast.Node) ast.Node {
	return wrap(ast.NewEmphasis(2), nodes...)
}

func italic(nodes ...ast.Node) ast.Node {
	return wrap(ast.NewEmphasis(1), nodes...)
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
		func() { _ = toNodeRichTextWithoutAnnotations(notion.RichText{}) })

	assert.Equal(t, newString("foo"), toNodeRichTextWithoutAnnotations(richText("foo", true, true)))

	for _, tt := range []struct {
		name string
		rts  notion.RichTexts
		res  []ast.Node
	}{
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
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ns := toNodeRichTexts(tt.rts)
			assert.Equal(t, tt.res, ns)
		})

	}

	// we want to wrap the nodes
	// - n1(bold, italic), n2(italic), n3(bold, italic) -> italic{bold{n1}, n2, bold{n3}}
	// - n1(bold), n2(bold, italic), n3(italic) -> bold{n1, italic{n2}}, italic{n3}
}

func TestTricky(t *testing.T) {
	t.Parallel()

	rts := notion.RichTexts{}
	assert.NoError(t, json.Unmarshal([]byte(tricky), &rts))

	// TODO continue here
	_ = toNodeRichTexts(rts)
}

const tricky = `[
          {
            "type": "text",
            "text": {
              "content": "This ",
              "link": null
            },
            "annotations": {
              "bold": false,
              "italic": false,
              "strikethrough": false,
              "underline": false,
              "code": false,
              "color": "default"
            },
            "plain_text": "This ",
            "href": null
          },
          {
            "type": "text",
            "text": {
              "content": "is",
              "link": null
            },
            "annotations": {
              "bold": false,
              "italic": true,
              "strikethrough": false,
              "underline": false,
              "code": false,
              "color": "purple"
            },
            "plain_text": "is",
            "href": null
          },
          {
            "type": "text",
            "text": {
              "content": " a big ",
              "link": null
            },
            "annotations": {
              "bold": false,
              "italic": false,
              "strikethrough": false,
              "underline": false,
              "code": false,
              "color": "purple"
            },
            "plain_text": " a big ",
            "href": null
          },
          {
            "type": "text",
            "text": {
              "content": "heading",
              "link": null
            },
            "annotations": {
              "bold": true,
              "italic": false,
              "strikethrough": false,
              "underline": false,
              "code": false,
              "color": "purple"
            },
            "plain_text": "heading",
            "href": null
          },
          {
            "type": "text",
            "text": {
              "content": " 1",
              "link": null
            },
            "annotations": {
              "bold": false,
              "italic": false,
              "strikethrough": false,
              "underline": false,
              "code": false,
              "color": "default"
            },
            "plain_text": " 1",
            "href": null
          }
        ]`
