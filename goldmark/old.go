package goldmark

import (
	"github.com/faetools/go-notion/pkg/notion"
	n_ast "github.com/faetools/notion-to-goldmark/ast"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
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

func toNodeListItem(item *notion.Paragraph) ast.Node {
	n := ast.NewListItem(0)

	// TODO
	// appendCommon(n, item.RichText, item.Children)

	return wrapInColor(item.Color, n)
}

func toNodeCallout(callout *notion.Callout) ast.Node {
	n := n_ast.NewCallout(callout.Icon)

	// TODO
	// appendCommon(n, callout.RichText, callout.Children)

	return wrapInColor(callout.Color, n)
}

func toNodeChildDatabase(db *notion.Child) ast.Node {
	return &n_ast.ChildDatabase{DB: *db}
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

func toNodeEmbed(embed *notion.Embed) ast.Node {
	n := &n_ast.Embed{}

	link := ast.NewLink()
	link.Destination = []byte(embed.Url)

	n.AppendChild(n, link)

	addCaption(n, &embed.Caption)

	return n
}

func toNodeEquation(eq *notion.Equation) ast.Node {
	return &n_ast.Equation{Expression: eq.Expression}
}

func toNodeHeading(h *notion.Paragraph, level int) ast.Node {
	n := ast.NewHeading(level)

	appendCommon(n, h.RichText, nil)

	return wrapInColor(h.Color, n)
}

func appendCommon(n ast.Node, rts notion.RichTexts, children notion.Blocks) {
	for _, child := range toNodeRichTexts(rts) {
		n.AppendChild(n, child)
	}

	addBlockChildren(n, children)
}

func toNodeLinkPreview(pr *notion.LinkPreview) ast.Node {
	n := &n_ast.LinkPreview{}

	link := ast.NewLink()
	link.Destination = []byte(pr.Url)
	n.AppendChild(n, link)

	return n
}

func toNodeQuote(q *notion.Paragraph) ast.Node {
	n := ast.NewBlockquote()

	// TODO
	// appendCommon(n, q.RichText, q.Children)

	return wrapInColor(q.Color, n)
}

func toNodeSyncedBlock(b *notion.SyncedBlock) ast.Node {
	n := n_ast.NewSyncedBlock()

	// TODO
	// addBlockChildren(n, b.Children)

	return n
}

func toNodeTableOfContents(toc *notion.TableOfContents) ast.Node {
	return wrapInColor(toc.Color, &n_ast.TableOfContents{})
}

func toNodeToDo(todo *notion.ToDo) ast.Node {
	n := extast.NewTaskCheckBox(todo.Checked)

	// TODO
	// appendCommon(n, todo.RichText, todo.Children)

	return wrapInColor(todo.Color, n)
}

func toNodeToggle(t *notion.Paragraph) ast.Node {
	n := &n_ast.Toggle{}

	// TODO
	// appendCommon(n, t.RichText, t.Children)

	return wrapInColor(t.Color, n)
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

func toNodeCode(c *notion.Code) ast.Node {
	lang := ast.NewText()

	n := ast.NewFencedCodeBlock(lang)

	n.SetAttributeString("language", c.Language)
	panic("code not yet implemented")
}
