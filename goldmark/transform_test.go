package goldmark_test

import (
	"context"
	"testing"

	"github.com/faetools/go-notion/pkg/fake"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/faetools/notion-to-goldmark/goldmark"
	"github.com/stretchr/testify/assert"
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
	// for _, b := range goldmark.FromBlocks(testBlocks) {
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

func TestTransform(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Skip()

	cli, fs, err := fake.NewClient()
	assert.NoError(t, err)

	assert.PanicsWithValue(t, "invalid memory address or nil pointer dereference", func() {
		_, err = goldmark.FromBlocks(ctx, cli, fake.PageID)
		assert.NoError(t, err)
	})

	assert.Len(t, fs.Unseen(), 30)
}
