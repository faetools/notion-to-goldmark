package goldmark

import (
	"context"

	"github.com/faetools/go-notion/pkg/docs"
	"github.com/faetools/go-notion/pkg/notion"
	n_ast "github.com/faetools/notion-to-goldmark/ast"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

// FromBlocks returns the goldmark nodes of the notion blocks.
func FromBlocks(ctx context.Context, g notion.Getter, id notion.Id) ([]ast.Node, error) {
	c := &nodeCollector{}

	v := docs.NewVisitor(g,
		func(p *notion.Page) error { return nil },
		func(blocks notion.Blocks) error {
			for _, b := range blocks {
				c.collectBlock(b)
			}

			if c.list != nil {
				c.res = append(c.res, c.list)
			}

			return nil
		},
		func(db *notion.Database) error { return nil },
		nil)

	if err := docs.Walk(ctx, v, docs.TypeBlocks, id); err != nil {
		return nil, err
	}

	return c.collectBlocks(ctx, id)
}

type nodeCollector struct {
	list *ast.List
	res  []ast.Node
}

func (c *nodeCollector) collectBlocks(ctx context.Context, id notion.Id) ([]ast.Node, error) {
	if c.list != nil {
		return append(c.res, c.list), nil
	}

	return c.res, nil
}

func (c *nodeCollector) getList(ordered bool) *ast.List {
	switch {
	case c.list == nil:
		// create new list
		c.list = newList(ordered)
	case c.list.IsOrdered() != ordered:
		// finish list
		c.res = append(c.res, c.list)

		// create a new list of different type
		c.list = newList(ordered)
	}

	return c.list
}

func newList(ordered bool) *ast.List {
	if ordered {
		return ast.NewList('.')
	}

	return ast.NewList('-')
}

func (c *nodeCollector) collectBlock(b notion.Block) {
	switch b.Type {
	case notion.BlockTypeNumberedListItem:
		setParentChild(c.getList(true), c.toNode(b))
	case notion.BlockTypeBulletedListItem:
		setParentChild(c.getList(false), c.toNode(b))
	default:
		// non-list, so finish existing list
		if c.list != nil {
			c.res = append(c.res, c.list)
			c.list = nil
		}

		c.res = append(c.res, c.toNode(b))
	}
}

func (c *nodeCollector) toNode(b notion.Block) ast.Node {
	switch b.Type {
	case notion.BlockTypeBookmark:
		return toNodeBookmark(b.Bookmark)
	case notion.BlockTypeBulletedListItem:
		return toNodeListItem(b.BulletedListItem)
	case notion.BlockTypeNumberedListItem:
		return toNodeListItem(b.NumberedListItem)
	case notion.BlockTypeCallout:
		return toNodeCallout(b.Callout)
	case notion.BlockTypeChildDatabase:
		return toNodeChildDatabase(b.ChildDatabase)
	case notion.BlockTypeChildPage:
		return n_ast.NewChildPage(*b.ChildPage)
	case notion.BlockTypeTable:
		// NOTE: toNode should never be called with notion.BlockTypeTableRow
		// the below function will call the appropriate methods
		return toNodeTable(b.Table)
	case notion.BlockTypeDivider:
		return ast.NewThematicBreak()
	case notion.BlockTypeEmbed:
		return toNodeEmbed(b.Embed)
	case notion.BlockTypeEquation:
		return toNodeEquation(b.Equation)
	case notion.BlockTypeFile:
		return n_ast.NewFile(*b.File, n_ast.FileTypeGeneric)
	case notion.BlockTypeHeading1:
		return toNodeHeading(b.Heading1, 1)
	case notion.BlockTypeHeading2:
		return toNodeHeading(b.Heading2, 2)
	case notion.BlockTypeHeading3:
		return toNodeHeading(b.Heading3, 3)
	case notion.BlockTypeImage:
		return n_ast.NewFile(*b.File, n_ast.FileTypeImage)
	case notion.BlockTypeLinkPreview:
		return toNodeLinkPreview(b.LinkPreview)
	case notion.BlockTypeLinkToPage:
		return n_ast.NewLinkToPage(*b.LinkToPage)
	case notion.BlockTypeParagraph:
		return toNodeParagraph(b.Paragraph)
	case notion.BlockTypePdf:
		return n_ast.NewFile(*b.Pdf, n_ast.FileTypePDF)
	case notion.BlockTypeQuote:
		return toNodeQuote(b.Quote)
	case notion.BlockTypeSyncedBlock:
		return toNodeSyncedBlock(b.SyncedBlock)
	case notion.BlockTypeTableOfContents:
		return toNodeTableOfContents(b.TableOfContents)
	case notion.BlockTypeToDo:
		return toNodeToDo(b.ToDo)
	case notion.BlockTypeToggle:
		return toNodeToggle(b.Toggle)
	case notion.BlockTypeVideo:
		return toNodeVideo(b.Video)
	default: // includes
		// notion.BlockTypeUnsupported (which we'll never support by its nature)
		// notion.BlockTypeColumn, notion.BlockTypeColumnList (which we plan to support in the future)
		// notion.BlockTypeTemplate (we're unsure when this is returned)
		panic("unknown node notion block type " + b.Type)
	}
}

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

func toNodeListItem(item *notion.ListItem) ast.Node {
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

func toNodeHeading(h *notion.Heading, level int) ast.Node {
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

func toNodeParagraph(p *notion.Paragraph) ast.Node {
	n := ast.NewParagraph()

	// TODO
	// appendCommon(n, p.RichText, p.Children)

	return wrapInColor(p.Color, n)
}

func toNodeQuote(q *notion.Quote) ast.Node {
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

func toNodeToggle(t *notion.Toggle) ast.Node {
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
		Caption:  v.Caption,
	}, n_ast.FileTypeVideo))

	if v.Caption != nil {
		addCaption(n, v.Caption)
	}

	return n
}
