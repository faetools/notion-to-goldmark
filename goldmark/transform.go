package goldmark

import (
	"context"
	"fmt"
	"net/http"

	"github.com/faetools/client"
	"github.com/faetools/go-notion/pkg/docs"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

type Getter interface {
	notion.Getter
	GetBlock(ctx context.Context, id notion.Id, reqEditors ...client.RequestEditorFn) (*notion.GetBlockResponse, error)
}

func FromBlock(ctx context.Context, cli Getter, id notion.Id) (ast.Node, error) {
	c := &nodeCollector{}

	resp, err := cli.GetBlock(ctx, id)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status(), string(resp.Body))
	}

	return c.toNode(*resp.JSON200), nil
}

// FromBlocks returns the goldmark nodes of the notion blocks.
func FromBlocks(ctx context.Context, cli notion.Getter, id notion.Id) ([]ast.Node, error) {
	c := &nodeCollector{}

	v := docs.NewVisitor(cli,
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
	case notion.BlockTypeParagraph:
		return toNodeParagraph(b.Id, b.Paragraph)

	// 	// TODO validate:
	// case notion.BlockTypeBookmark:
	// 	return toNodeBookmark(b.Bookmark)
	// case notion.BlockTypeBulletedListItem:
	// 	return toNodeListItem(b.BulletedListItem)
	// case notion.BlockTypeNumberedListItem:
	// 	return toNodeListItem(b.NumberedListItem)
	// case notion.BlockTypeCallout:
	// 	return toNodeCallout(b.Callout)
	// case notion.BlockTypeChildDatabase:
	// 	return toNodeChildDatabase(b.ChildDatabase)
	// case notion.BlockTypeChildPage:
	// 	return n_ast.NewChildPage(*b.ChildPage)
	// case notion.BlockTypeTable:
	// 	// NOTE: toNode should never be called with notion.BlockTypeTableRow
	// 	// the below function will call the appropriate methods
	// 	return toNodeTable(b.Table)
	// case notion.BlockTypeDivider:
	// 	return ast.NewThematicBreak()
	// case notion.BlockTypeEmbed:
	// 	return toNodeEmbed(b.Embed)
	// case notion.BlockTypeEquation:
	// 	return toNodeEquation(b.Equation)
	// case notion.BlockTypeFile:
	// 	return n_ast.NewFile(*b.File, n_ast.FileTypeGeneric)
	// case notion.BlockTypeHeading1:
	// 	return toNodeHeading(b.Heading1, 1)
	// case notion.BlockTypeHeading2:
	// 	return toNodeHeading(b.Heading2, 2)
	// case notion.BlockTypeHeading3:
	// 	return toNodeHeading(b.Heading3, 3)
	// case notion.BlockTypeImage:
	// 	return n_ast.NewFile(*b.File, n_ast.FileTypeImage)
	// case notion.BlockTypeLinkPreview:
	// 	return toNodeLinkPreview(b.LinkPreview)
	// case notion.BlockTypeLinkToPage:
	// 	return n_ast.NewLinkToPage(*b.LinkToPage)
	// case notion.BlockTypePdf:
	// 	return n_ast.NewFile(*b.Pdf, n_ast.FileTypePDF)
	// case notion.BlockTypeQuote:
	// 	return toNodeQuote(b.Quote)
	// case notion.BlockTypeSyncedBlock:
	// 	return toNodeSyncedBlock(b.SyncedBlock)
	// case notion.BlockTypeTableOfContents:
	// 	return toNodeTableOfContents(b.TableOfContents)
	// case notion.BlockTypeToDo:
	// 	return toNodeToDo(b.ToDo)
	// case notion.BlockTypeToggle:
	// 	return toNodeToggle(b.Toggle)
	// case notion.BlockTypeVideo:
	// 	return toNodeVideo(b.Video)
	// case notion.BlockTypeCode:
	// 	return toNodeCode(b.Code)
	default: // includes
		// notion.BlockTypeUnsupported (which we'll never support by its nature)
		// notion.BlockTypeColumn, notion.BlockTypeColumnList (which we plan to support in the future)
		// notion.BlockTypeTemplate (we're unsure when this is returned)
		panic(fmt.Sprintf("unknown node notion block type %q", b.Type))
	}
}

func setID(n ast.Node, id notion.UUID) ast.Node {
	n.SetAttributeString("id", []byte(id))
	n.SetAttributeString("class", []byte(""))

	return n
}

func toNodeParagraph(id notion.UUID, p *notion.Paragraph) ast.Node {
	n := setID(ast.NewParagraph(), id)

	for _, child := range toNodeRichTexts(p.RichText) {
		n.AppendChild(n, child)
	}

	return wrapInColor(p.Color, n)
}
