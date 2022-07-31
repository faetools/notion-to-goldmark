package goldmark_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	extast "github.com/yuin/goldmark/extension/ast"

	n_ast "github.com/faetools/notion-to-goldmark/ast"

	"github.com/faetools/go-notion/pkg/fake"
	"github.com/faetools/go-notion/pkg/notion"
	. "github.com/faetools/notion-to-goldmark/goldmark"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

const max = 34

func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(n_ast.KindColor, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			w.WriteString("</mark>")
			return ast.WalkContinue, nil
		}

		w.WriteString(`<mark class="highlight-`)
		w.Write([]byte(n.(*n_ast.Color).Color))
		w.WriteString(`">`)

		return ast.WalkContinue, nil
	})

	reg.Register(ast.KindCodeSpan, renderTag("code", html.CodeAttributeFilter))

	reg.Register(n_ast.KindUnderline, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			_, _ = w.WriteString(`<span style="border-bottom:0.05em solid">`)
			return ast.WalkContinue, nil
		}

		_, _ = w.WriteString("</span>")
		return ast.WalkContinue, nil
	})

	reg.Register(ast.KindParagraph, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString("</p>")
			return ast.WalkContinue, nil
		}

		_, _ = w.WriteString("<p")
		html.RenderAttributes(w, n, html.ParagraphAttributeFilter)
		_ = w.WriteByte('>')

		if !n.HasChildren() {
			// for some reason, notion adds a new line here
			_ = w.WriteByte('\n')
		}

		return ast.WalkContinue, nil
	})

	reg.Register(n_ast.KindBlockChildren, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch n.Parent().Kind() {
		case ast.KindParagraph: // ok
		case n_ast.KindCalloutText, n_ast.KindSyncedBlock,
			ast.KindListItem, n_ast.KindToggle, n_ast.KindChildPage:
			// NOTE: maybe we change the implementation so as to not use KindBlockChildren
			// and instead call it "indented" or something
			return ast.WalkContinue, nil
		default:
			fmt.Printf("unknown block children parent: %v\n", n.Parent().Kind())
		}

		if entering {
			_, _ = w.WriteString(`<div class="indented">`)
		} else {
			_, _ = w.WriteString(`</div>`)
		}

		return ast.WalkContinue, nil
	})

	reg.Register(ast.KindHeading, func(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		n := node.(*ast.Heading)
		if entering {
			_, _ = w.WriteString("<h")
			_ = w.WriteByte("0123456"[n.Level])
			html.RenderAttributes(w, node, html.HeadingAttributeFilter)
		} else {
			_, _ = w.WriteString("</h")
			_ = w.WriteByte("0123456"[n.Level])
		}

		_ = w.WriteByte('>')
		return ast.WalkContinue, nil
	})

	renderFigure := renderTag("figure", html.GlobalAttributeFilter)

	reg.Register(n_ast.KindCallout, renderFigure)

	reg.Register(n_ast.KindIcon, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString(`</div>`)
			return ast.WalkContinue, nil
		}

		_, _ = w.WriteString(`<div style="font-size:1.5em">`)

		if emoji := n.(*n_ast.Icon).Emoji; emoji != "" {
			_, _ = w.WriteString(`<span class="icon">`)
			_, _ = w.WriteString(emoji)
			_, _ = w.WriteString(`</span>`)
		}

		return ast.WalkContinue, nil
	})

	reg.Register(n_ast.KindCalloutText, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			_, _ = w.WriteString(`<div style="width:100%">`)
		} else {
			_, _ = w.WriteString(`</div>`)
		}

		return ast.WalkContinue, nil
	})

	reg.Register(ast.KindImage, func(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		n := node.(*ast.Image)

		_, _ = w.WriteString("<img")
		html.RenderAttributes(w, n, html.ImageAttributeFilter)

		_, _ = w.WriteString(` src="`)
		_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		_ = w.WriteByte('"')

		// TODO remove in favor of below code
		html.RenderAttributes(w, n, util.NewBytesFilter([]byte("alt")))

		if alt := n.Text(source); len(alt) > 0 {
			_, _ = w.WriteString(`" alt="`)
			_, _ = w.Write(util.EscapeHTML(alt))
			_ = w.WriteByte('"')
		}

		if n.Title != nil {
			_, _ = w.WriteString(` title="`)
			_, _ = w.Write(n.Title)
			_ = w.WriteByte('"')
		}

		_, _ = w.WriteString("/>")

		return ast.WalkSkipChildren, nil
	})

	reg.Register(ast.KindBlockquote, renderTag("blockquote", html.BlockquoteAttributeFilter))

	renderDiv := renderTag("div", html.GlobalAttributeFilter)
	noop := func(util.BufWriter, []byte, ast.Node, bool) (ast.WalkStatus, error) {
		return ast.WalkContinue, nil
	}

	reg.Register(n_ast.KindSyncedBlock, renderDiv)
	reg.Register(extast.KindTaskCheckBox, renderDiv)

	idFilter := util.NewBytesFilter([]byte("id"))

	// notion prints all list items in their own list
	reg.Register(ast.KindList, noop)
	reg.Register(ast.KindListItem, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		p := n.Parent().(*ast.List)

		tag := "ul"
		if p.IsOrdered() {
			tag = "ol"
		}

		if !entering {
			_, _ = w.WriteString("</li></")
			_, _ = w.WriteString(tag)
			_ = w.WriteByte('>')
			return ast.WalkContinue, nil
		}

		_ = w.WriteByte('<')
		_, _ = w.WriteString(tag)

		depth := listDepth(p, p.IsOrdered())

		if p.IsOrdered() {
			// different types depending on depth
			_, _ = w.WriteString(` type="`)
			_, _ = w.WriteString([]string{"1", "a", "i"}[depth%3])
			_ = w.WriteByte('"')
		}

		// html.RenderAttributes(w, p, typeFilter)
		html.RenderAttributes(w, n, idFilter)

		// e.g. class="block-color-red_background numbered-list" or class="red to-do-list"
		renderClass(w, n, p)

		if p.IsOrdered() {
			_, _ = w.WriteString(` start="`)
			_, _ = w.WriteString(strconv.Itoa(numPrevious(n) + 1))
			_ = w.WriteByte('"')
		}

		isToDo := false
		if c := n.FirstChild(); c != nil && c.Kind() == extast.KindTaskCheckBox {
			isToDo = true
		}

		_, _ = w.WriteString("><li")
		if !p.IsOrdered() && !isToDo {
			// different styles depending on depth
			_, _ = w.WriteString(` style="list-style-type:`)
			_, _ = w.WriteString([]string{"disc", "circle", "square"}[depth%3])
			_ = w.WriteByte('"')
		}

		_ = w.WriteByte('>')

		return ast.WalkContinue, nil
	})

	reg.Register(n_ast.KindCheckboxText, renderTag("span", html.ListItemAttributeFilter))

	reg.Register(n_ast.KindToggle, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString("</details></li></ul>")
			return ast.WalkContinue, nil
		}

		_, _ = w.WriteString("<ul")
		html.RenderAttributes(w, n, idFilter)
		renderClass(w, n)
		_, _ = w.WriteString(`><li><details open="">`)

		return ast.WalkContinue, nil
	})

	reg.Register(n_ast.KindToggleText, renderTag("summary", html.GlobalAttributeFilter))

	reg.Register(ast.KindFencedCodeBlock, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString("</code></pre>")
			return ast.WalkContinue, nil
		}

		_, _ = w.WriteString("<pre")
		html.RenderAttributes(w, n, html.CodeAttributeFilter)
		_, _ = w.WriteString("><code>")

		return ast.WalkContinue, nil
	})

	renderLink := func(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString("</a>")
			return ast.WalkContinue, nil
		}

		n := node.(*ast.Link)

		_, _ = w.WriteString(`<a href="`)
		_, _ = w.Write(util.EscapeHTML(n.Destination))
		_ = w.WriteByte('"')

		if n.Title != nil {
			_, _ = w.WriteString(` title="`)
			_, _ = w.Write(n.Title)
			_ = w.WriteByte('"')
		}

		html.RenderAttributes(w, n, html.LinkAttributeFilter)

		_ = w.WriteByte('>')

		return ast.WalkContinue, nil
	}

	reg.Register(ast.KindLink, func(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		// don't render links in fenced code blocks
		for p := node.Parent(); p != nil; p = p.Parent() {
			if p.Kind() == ast.KindFencedCodeBlock {
				return ast.WalkContinue, nil
			}
		}

		return renderLink(w, nil, node, entering)
	})

	reg.Register(n_ast.KindChildPage, renderFigure)
	reg.Register(n_ast.KindEmbed, renderFigure)

	reg.Register(n_ast.KindEmbedSource, func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString("</div>")
			return ast.WalkContinue, nil
		}

		_, _ = w.WriteString(`<div class="source">`)

		return ast.WalkContinue, nil
	})

	reg.Register(n_ast.KindCaption, func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		// notion ignores the caption of fenced code blocks
		if n.Parent().Kind() == ast.KindFencedCodeBlock {
			return ast.WalkSkipChildren, nil
		}

		return renderTag("figcaption", html.GlobalAttributeFilter)(w, source, n, entering)
	})

	reg.Register(n_ast.KindChildDatabase, renderDiv)

	reg.Register(extast.KindTable, renderTag("table", extension.TableAttributeFilter))

	reg.Register(extast.KindTableHeader, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			_, _ = w.WriteString("<thead")
			html.RenderAttributes(w, n, extension.TableHeaderAttributeFilter)
			_, _ = w.WriteString("><tr>") // Header <tr> has no separate handle
		} else {
			_, _ = w.WriteString("</tr>")
			_, _ = w.WriteString("</thead>")
			if n.NextSibling() != nil {
				_, _ = w.WriteString("<tbody>")
			}
		}

		return ast.WalkContinue, nil
	})

	reg.Register(extast.KindTableRow, func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			_, _ = w.WriteString("<tr")
			html.RenderAttributes(w, n, extension.TableRowAttributeFilter)
			_, _ = w.WriteString(">")
		} else {
			_, _ = w.WriteString("</tr>")

			if n.Parent().LastChild() == n {
				_, _ = w.WriteString("</tbody>")
			}
		}
		return ast.WalkContinue, nil
	})

	reg.Register(extast.KindTableCell, func(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		n := node.(*extast.TableCell)

		tag := "td"
		if n.Parent().Kind() == extast.KindTableHeader {
			tag = "th"
		}

		if entering {
			fmt.Fprintf(w, "<%s", tag)
			if n.Alignment != extast.AlignNone {
				// amethod := r.TableConfig.TableCellAlignMethod
				// if amethod == TableCellAlignDefault {
				// 	if r.Config.XHTML {
				// 		amethod = TableCellAlignAttribute
				// 	} else {
				// 		amethod = TableCellAlignStyle
				// 	}
				// }
				// switch amethod {
				// case TableCellAlignAttribute:
				if _, ok := n.AttributeString("align"); !ok { // Skip align render if overridden
					fmt.Fprintf(w, ` align="%s"`, n.Alignment.String())
				}
				// case TableCellAlignStyle:
				// 	v, ok := n.AttributeString("style")
				// 	var cob util.CopyOnWriteBuffer
				// 	if ok {
				// 		cob = util.NewCopyOnWriteBuffer(v.([]byte))
				// 		cob.AppendByte(';')
				// 	}
				// 	style := fmt.Sprintf("text-align:%s", n.Alignment.String())
				// 	cob.AppendString(style)
				// 	n.SetAttributeString("style", cob.Bytes())
				// }
			}
			if tag == "td" {
				html.RenderAttributes(w, n, extension.TableTdCellAttributeFilter) // <td>
			} else {
				html.RenderAttributes(w, n, extension.TableThCellAttributeFilter) // <th>
			}
			_ = w.WriteByte('>')
		} else {
			fmt.Fprintf(w, "</%s>", tag)
		}

		return ast.WalkContinue, nil
	})

	reg.Register(n_ast.KindPropertyIcon, renderTag("span", html.GlobalAttributeFilter))
	reg.Register(n_ast.KindSVG, renderTag("svg", svgFilter))
	reg.Register(n_ast.KindSVGPath, renderTag("path", pathFilter))
	reg.Register(n_ast.KindPolygon, renderTag("polygon", polygonFilter))
	reg.Register(n_ast.KindStatus, func(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString(`</span>`)
			return ast.WalkContinue, nil
		}

		data := node.(*n_ast.Status).Data

		_, _ = w.WriteString(`<span class="status-value select-value-color-`)
		_, _ = w.WriteString(string(data.Color))
		_, _ = w.WriteString(`"><div class="status-dot status-dot-color-`)
		_, _ = w.WriteString(string(data.Color))
		_, _ = w.WriteString(`"></div>`)
		_, _ = w.WriteString(data.Name)

		return ast.WalkSkipChildren, nil
	})

	reg.Register(n_ast.KindDate, func(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString(`</time>`)
			return ast.WalkContinue, nil
		}

		_, _ = w.WriteString(`<time>`)

		n := node.(*n_ast.Date)

		if n.Date == nil {
			return ast.WalkContinue, nil
		}

		_ = w.WriteByte('@')

		_, _ = w.WriteString(formatTime(n.Date.Start, n.TwelveHourClock))

		if n.Date.End == nil {
			return ast.WalkContinue, nil
		}

		_, _ = w.WriteString(" → ")
		_, _ = w.WriteString(formatTime(*n.Date.End, n.TwelveHourClock))

		return ast.WalkContinue, nil
	})

	reg.Register(n_ast.KindFileInCell, func(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString(`</span>`)
			return ast.WalkContinue, nil
		}

		_, _ = w.WriteString(`<span style="margin-right:6px">`)

		return ast.WalkContinue, nil
	})

	reg.Register(n_ast.KindUser, func(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString(`</span>`)
			return ast.WalkContinue, nil
		}

		user := node.(*n_ast.User).Data

		_, _ = w.WriteString(`<span class="user">`)

		if user.AvatarUrl != nil {
			_, _ = w.WriteString(`<img src="`)
			_, _ = w.WriteString(*user.AvatarUrl)
			_, _ = w.WriteString(`" class="icon user-icon"/>`)
		}

		if user.Name != nil {
			_, _ = w.WriteString(*user.Name)
		}

		return ast.WalkContinue, nil
	})

	reg.Register(n_ast.KindSelect, func(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			_, _ = w.WriteString(`</span>`)
			return ast.WalkContinue, nil
		}

		val := node.(*n_ast.Select).Data

		_, _ = w.WriteString(`<span class="selected-value`)

		switch val.Color {
		case notion.ColorDefault, "": // no color
		default:
			_, _ = w.WriteString(` select-value-color-`)
			_, _ = w.WriteString(string(val.Color))
		}

		_, _ = w.WriteString(`">`)
		_, _ = w.WriteString(val.Name)

		return ast.WalkContinue, nil
	})
}

var (
	svgFilter     = html.GlobalAttributeFilter.Extend([]byte("viewBox"))
	pathFilter    = util.NewBytesFilter([]byte("d"))
	polygonFilter = util.NewBytesFilter([]byte("points"))
)

func formatTime(ts time.Time, twelveHourClock bool) string {
	if ts.Minute() == 0 && ts.Hour() == 0 {
		return ts.Format("January 2, 2006")
	}

	if twelveHourClock {
		return ts.Local().Format("January 2, 2006 3:04 PM")
	}

	// we need to remove the zero
	// e.g. "August 12, 2022 03:00" -> "August 12, 2022 3:00"
	return strings.Replace(ts.Local().Format("January 2, 2006 15:04"), " 0", " ", 1)
}

// renderTag factories out a simple function to render a tag
func renderTag(tagName string, filter util.BytesFilter) func(util.BufWriter, []byte, ast.Node, bool) (ast.WalkStatus, error) {
	start := "<" + tagName
	end := "</" + tagName + ">"

	return func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			_, _ = w.WriteString(start)
			html.RenderAttributes(w, n, filter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString(end)
		}

		return ast.WalkContinue, nil
	}
}

func listDepth(l ast.Node, ordered bool) int {
	switch p := l.Parent().(type) {
	case *ast.List:
		if p.IsOrdered() != ordered {
			return 0
		}

		return listDepth(p, ordered) + 1
	case *n_ast.BlockChildren, *ast.ListItem:
		return listDepth(p, ordered)
	default:
		return 0
	}
}

func renderClass(w util.BufWriter, nodes ...ast.Node) {
	_, _ = w.WriteString(` class="`)

	classes := [][]byte{}

	for _, n := range nodes {
		cl, _ := n.AttributeString("class")
		class, _ := cl.([]byte)

		if len(class) == 0 {
			continue
		}

		classes = append(classes, class)
	}

	_, _ = w.Write(bytes.Join(classes, []byte{' '}))
	_ = w.WriteByte('"')
}

func TestFromBlock(t *testing.T) {
	t.Parallel()

	// set local time zone to the time zone
	// the example page was exported to
	loc, err := time.LoadLocation("Europe/Berlin")
	assert.NoError(t, err)

	time.Local = loc

	ctx := context.Background()

	cli, _, err := fake.NewClient()
	assert.NoError(t, err)

	doc := ast.NewDocument()

	nodes, err := GetPage(ctx, cli, fake.PageID, max)
	assert.NoError(t, err)

	for _, n := range nodes {
		doc.AppendChild(doc, n)
	}

	pageRoot := func(name string, id notion.Id) string {
		return fmt.Sprintf("%s %s", name, strings.ReplaceAll(string(id), "-", ""))
	}

	root := pageRoot("Example Page", fake.PageID)

	assert.NoError(t, ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		switch n := node.(type) {
		case *ast.Image:
			if _, ok := n.AttributeString("expires"); ok {
				u, err := url.Parse(string(n.Destination))
				if err != nil {
					return 0, err
				}

				// download the file to root with fileName here
				fileName := filepath.Base(u.Path)

				n.Destination = util.URLEscape([]byte(filepath.Join(root, fileName)), true)
			}
		}

		return 0, nil
	}))

	// doc.Dump(nil, 0)

	w := &bytes.Buffer{}
	assert.NoError(t, r.Render(w, nil, doc))
	got := w.Bytes()

	start := 14705
	html, err := fake.HTMLExport.ReadFile("html/" + root + ".html")
	want := html[start:]

	for i, b := range want {
		if len(got) > i && got[i] == b {
			continue
		}

		k := i - 6

		offset := 3
		end := i + bytes.Index(want[i+offset:], []byte(">")) + 1 + offset

		assert.Equal(t, string(want[k:end]), string(got[k:]), "from %d", k)
		break
	}
}

func paragraphBlock(p *notion.Paragraph) notion.Block {
	return notion.Block{
		Object:    "block",
		Type:      "paragraph",
		Paragraph: p,
	}
}

func numPrevious(n ast.Node) int {
	prev := n.PreviousSibling()
	if prev == nil {
		return 0
	}

	return numPrevious(prev) + 1
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
		util.Prioritized(extension.NewStrikethroughHTMLRenderer(), 100),
		// util.Prioritized(extension.NewTableHTMLRenderer(), 99),
		util.Prioritized(&htmlRenderer{}, 100),
	))
}

type htmlRenderer struct{}

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

// func TestTransform_Errors(t *testing.T) {
// 	t.Parallel()
// 	ctx := context.Background()

// 	testError := errors.New("some error")

// 	_, err := FromBlocks(ctx, &testGetter{
// 		blocks: func(id notion.Id) (notion.Blocks, error) {
// 			return nil, testError
// 		},
// 	}, "")

// 	assert.ErrorIs(t, err, testError)
// }

func TestTransform_Panics(t *testing.T) {
	t.Parallel()
	// ctx := context.Background()

	// assert.PanicsWithValue(t, `unknown node notion block type "foo"`, func() {
	// 	_, _ = FromBlocks(ctx, &testGetter{
	// 		blocks: func(id notion.Id) (notion.Blocks, error) {
	// 			return notion.Blocks{{Type: "foo"}}, nil
	// 		},
	// 	}, "")
	// })
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
