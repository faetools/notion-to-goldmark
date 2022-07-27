package goldmark_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	n_ast "github.com/faetools/notion-to-goldmark/ast"

	"github.com/faetools/go-notion/pkg/fake"
	"github.com/faetools/go-notion/pkg/notion"
	. "github.com/faetools/notion-to-goldmark/goldmark"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func paragraphBlock(p *notion.Paragraph) notion.Block {
	return notion.Block{
		Object:    "block",
		Type:      "paragraph",
		Paragraph: p,
	}
}

// var testBlocks = notion.Blocks{
// 	paragraphBlock(&notion.Paragraph{
// 		Color: notion.ColorGreen,
// 		RichText: []notion.RichText{
// 			notion.NewRichText("This is the "),
// 			{
// 				Type:      notion.RichTextTypeText,
// 				PlainText: "home",
// 				Text: &notion.Text{
// 					Content: "home",
// 				},
// 				Annotations: notion.Annotations{
// 					Bold:  true,
// 					Color: notion.ColorDefault,
// 				},
// 			},
// 			notion.NewRichText(" page."),
// 		},
// 	}),
// 	paragraphBlock(&notion.Paragraph{
// 		Color:    notion.ColorDefault,
// 		RichText: []notion.RichText{},
// 	}),
// 	paragraphBlock(&notion.Paragraph{
// 		Color: notion.ColorDefault,
// 		RichText: []notion.RichText{
// 			{
// 				Type:      notion.RichTextTypeEquation,
// 				PlainText: `\sum 2+2 = 4\sigma`,
// 				Equation: &notion.Equation{
// 					Expression: `\sum 2+2 = 4\sigma`,
// 				},
// 				Annotations: notion.Annotations{Color: notion.ColorDefault},
// 			},
// 		},
// 	}),
// 	paragraphBlock(&notion.Paragraph{
// 		Color: notion.ColorDefault,
// 		RichText: []notion.RichText{
// 			{
// 				Type:      notion.RichTextTypeMention,
// 				PlainText: "@Mark Rösler",
// 				Mention: &notion.Mention{
// 					Type: notion.MentionTypeUser,
// 					User: &notion.User{
// 						Object:    "user",
// 						Type:      notion.UserTypePerson,
// 						Id:        "af171d5d-c36f-45bc-a0a3-6086c0dafa45",
// 						Name:      "Mark Rösler",
// 						AvatarUrl: "https://lh3.googleusercontent.com/a-/AOh14Gi54BUKkLrZ2IX8ORURI__6avK9zjCYXdhbmthj=s100",

// 						Person: &notion.Person{
// 							Email: "mark@faetools.com",
// 						},
// 					},
// 				},
// 				Annotations: notion.Annotations{Color: notion.ColorDefault},
// 			},
// 			notion.NewRichText(" is a great guy."),
// 		},
// 	}),
// }

func TestMarkdown(t *testing.T) {
	t.Parallel()

	// doc := ast.NewDocument()
	// for _, b := range FromBlocks(testBlocks) {
	// 	doc.AppendChild(doc, b)
	// }

	// b, err := hugo.Render(nil, doc)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, `<span class=green>This is the **home** page.

	// </span>

	// $\sum 2+2 = 4\sigma$

	// Mark Rösler is a great guy.
	// `, string(b))
}

var r = goldmark.DefaultRenderer()

func init() {
	r.AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&htmlRenderer{}, 100),
	))
}

type htmlRenderer struct{}

func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(n_ast.KindColor, func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			w.WriteString("</mark>")
			return ast.WalkContinue, nil
		}

		w.WriteString(`<mark class="highlight-`)
		w.Write([]byte(n.(*n_ast.Color).Color))
		w.WriteString(`">`)

		return ast.WalkContinue, nil
	})
}

func TestFromBlock(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cli, _, err := fake.NewClient()
	assert.NoError(t, err)

	doc := ast.NewDocument()

	for _, id := range []notion.Id{
		"8fe3dd1a-c8d8-47ff-b2fb-70b0269c4e9f",
		"bb9ff6f9-c230-4688-b15c-15a6648e1c8a",
	} {
		n, err := FromBlock(ctx, cli, id)
		assert.NoError(t, err)

		doc.AppendChild(doc, n)
	}

	w := &bytes.Buffer{}

	doc.Dump(nil, 0)

	assert.NoError(t, r.Render(w, nil, doc))

	html, err := fake.HTMLExport.ReadFile("html/Example Page 96245c8f178444a482ad1941127c3ec3.html")
	assert.NoError(t, err)

	start := 14764
	end := start + 202
	assert.Equal(t, string(html[start:end]), strings.ReplaceAll(w.String(), "\n", ""))
}

// func TestTransform(t *testing.T) {
// 	t.Parallel()
// 	ctx := context.Background()

// 	cli, fs, err := fake.NewClient()
// 	assert.NoError(t, err)

// 	// TODO continue here
// 	assert.PanicsWithValue(t, `code not yet implemented`, func() {
// 		_, err = FromBlocks(ctx, cli, fake.PageID)
// 		assert.NoError(t, err)
// 	})

// 	assert.Len(t, fs.Unseen(), 164)
// }

func TestTransform_Errors(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	testError := errors.New("some error")

	_, err := FromBlocks(ctx, &testGetter{
		blocks: func(id notion.Id) (notion.Blocks, error) {
			return nil, testError
		},
	}, "")

	assert.ErrorIs(t, err, testError)
}

func TestTransform_Panics(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	assert.PanicsWithValue(t, `unknown node notion block type "foo"`, func() {
		_, _ = FromBlocks(ctx, &testGetter{
			blocks: func(id notion.Id) (notion.Blocks, error) {
				return notion.Blocks{{Type: "foo"}}, nil
			},
		}, "")
	})
}

type testRoundtripper struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (rt *testRoundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.roundTrip(req)
}

type testGetter struct {
	page            func(id notion.Id) (*notion.Page, error)
	blocks          func(id notion.Id) (notion.Blocks, error)
	database        func(id notion.Id) (*notion.Database, error)
	databaseEntries func(id notion.Id) (notion.Pages, error)
}

// GetNotionPage implements notion.Getter
func (g *testGetter) GetNotionPage(ctx context.Context, id notion.Id) (*notion.Page, error) {
	return g.page(id)
}

// GetAllBlocks implements notion.Getter
func (g *testGetter) GetAllBlocks(_ context.Context, id notion.Id) (notion.Blocks, error) {
	return g.blocks(id)
}

// GetNotionDatabase implements notion.Getter
func (g *testGetter) GetNotionDatabase(ctx context.Context, id notion.Id) (*notion.Database, error) {
	return g.database(id)
}

// GetAllDatabaseEntries implements notion.Getter
func (g *testGetter) GetAllDatabaseEntries(ctx context.Context, id notion.Id) (notion.Pages, error) {
	return g.databaseEntries(id)
}
