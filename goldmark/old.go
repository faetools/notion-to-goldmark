package goldmark

import (
	"fmt"

	"github.com/faetools/go-notion/pkg/notion"
	n_ast "github.com/faetools/notion-to-goldmark/ast"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/util"
)

func toNodeBookmark(b *notion.Bookmark) ast.Node {
	n := &n_ast.Bookmark{URL: b.Url}

	addCaption(n, &b.Caption)

	return n
}

func toNodeText(t *notion.Text) ast.Node {
	if t.Link != nil {
		// NOTE: Relative links will be notion pages.
		// Walk over all nodes and transform URLs to fit your purposes.
		return newLink(t.Content, t.Link.Url)
	}

	return newString(t.Content)
}

func toNodeTable(table *notion.Table) ast.Node {
	n := extast.NewTable()

	// TODO
	// for i, r := range table.Children {
	// 	row := extast.NewTableRow(nil)

	// 	// NOTE: len(r.TableRow.Cells) == table.TableWidth

	// 	for i, cellContent := range r.TableRow.Cells {
	// 		cell := extast.NewTableCell()

	// 		if i == 0 && table.HasColumnHeader {
	// 			cell.SetAttributeString("header", true)
	// 		}

	// 		appendCommon(cell, cellContent, nil)

	// 		row.AppendChild(row, cell)
	// 	}

	// 	if i == 0 && table.HasRowHeader {
	// 		n.AppendChild(n, extast.NewTableHeader(row))
	// 	} else {
	// 		n.AppendChild(n, row)
	// 	}
	// }

	return n
}

func toNodeEquation(eq *notion.Equation) ast.Node {
	return &n_ast.Equation{Expression: eq.Expression}
}

func appendCommon(n ast.Node, rts notion.RichTexts, children notion.Blocks) {
	for _, child := range toNodeRichTexts(rts) {
		n.AppendChild(n, child)
	}

	addBlockChildren(n, children)
}

// func toNodeLinkPreview(pr *notion.LinkPreview) ast.Node {
// 	n := &n_ast.LinkPreview{}

// 	link := ast.NewLink()
// 	link.Destination = []byte(pr.Url)
// 	n.AppendChild(n, link)

// 	return n
// }

func toNodeQuote(q *notion.Paragraph) ast.Node {
	n := ast.NewBlockquote()

	// TODO
	// appendCommon(n, q.RichText, q.Children)

	return wrapInColor(q.Color, n)
}

func toNodeTableOfContents(toc *notion.TableOfContents) ast.Node {
	return wrapInColor(toc.Color, &n_ast.TableOfContents{})
}

func toNodeVideo(v *notion.Video) ast.Node {
	n := &n_ast.Video{}

	// TODO can notion.Video be notion.FileWithCaption from the beginning?
	n.AppendChild(n, n_ast.NewFile(notion.FileWithCaption{
		External: v.External,
		File:     v.File,
		Type:     notion.FileWithCaptionType(v.Type),
		Caption:  &v.Caption,
	}, n_ast.FileTypeVideo))

	if v.Caption != nil {
		addCaption(n, &v.Caption)
	}

	return n
}

// setParentChild is a convenience method to not name the parent twice.
func setParentChild(parent, child ast.Node) {
	parent.AppendChild(parent, child)
}

func newLink(content, dest string) *ast.Link {
	n := ast.NewLink()
	n.Destination = util.URLEscape([]byte(dest), true)

	n.AppendChild(n, newString(content))

	return n
}

func newString(s string) *ast.String { return ast.NewString([]byte(s)) }

func newLinkToPage(id *notion.UUID) ast.Node {
	return newLink("", fmt.Sprintf("/%s", *id))
}

func wrapInColor(c notion.Color, child ast.Node) ast.Node {
	if c == notion.ColorDefault || c == "" {
		return child
	}

	n := &n_ast.Color{Color: c}
	n.AppendChild(n, child)

	return n
}

// addBlockChildren adds any children blocks
// TODO: call the API go get the children
func addBlockChildren(n ast.Node, bs notion.Blocks) {
	if len(bs) == 0 {
		return
	}

	// TODO
	// children := &n_ast.Children{}
	// for _, child := range FromBlocks(bs) {
	// 	children.AppendChild(children, child)
	// }

	// n.AppendChild(n, children)
}

func addCaption(n ast.Node, rts *notion.RichTexts) {
	if rts == nil || len(*rts) == 0 {
		return
	}

	caption := &n_ast.Caption{}

	for _, child := range toNodeRichTexts(*rts) {
		caption.AppendChild(caption, child)
	}

	setParentChild(n, caption)
}
