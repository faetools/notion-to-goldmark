package goldmark

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/faetools/go-notion/pkg/notion"
	n_ast "github.com/faetools/notion-to-goldmark/ast"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/util"
)

type pageCollector struct {
	root string

	ctx context.Context
	cli notion.Getter
}

type blockCollector struct {
	p *pageCollector

	list     *ast.List
	listType notion.BlockType
	res      []ast.Node
}

// GetPage returns the goldmark nodes of a notion page.
func GetPage(ctx context.Context, cli notion.Getter, id notion.Id, max int) ([]ast.Node, error) {
	p, err := cli.GetNotionPage(ctx, id)
	if err != nil {
		return nil, err
	}

	c := &pageCollector{
		root: getDir(p.Title(), p.Id),
		ctx:  ctx, cli: cli,
	}

	return c.getBlocks(id, max)
}

// getBlocks returns the goldmark nodes of a notion page.
func (c *pageCollector) getBlocks(id notion.Id, max int) ([]ast.Node, error) {
	blocks, err := c.cli.GetAllBlocks(c.ctx, id)
	if err != nil {
		return nil, err
	}

	bc := &blockCollector{p: c}

	for i, b := range blocks {
		if i == max {
			break
		}

		if err := bc.collectBlock(b); err != nil {
			return nil, err
		}
	}

	return bc.collectBlocks(c.ctx, id)
}

func (c *blockCollector) collectBlocks(ctx context.Context, id notion.Id) ([]ast.Node, error) {
	if c.list != nil {
		return append(c.res, c.list), nil
	}

	return c.res, nil
}

func (c *blockCollector) getList(tp notion.BlockType) *ast.List {
	if tp == notion.BlockTypeToggle {
		tp = notion.BlockTypeBulletedListItem
	}

	switch {
	case c.list == nil:
		// create new list
		c.list = newList(tp)
		c.listType = tp
	case tp != c.listType:
		// finish list
		c.res = append(c.res, c.list)

		// create a new list of different type
		c.list = newList(tp)
		c.listType = tp
	}

	return c.list
}

func newList(tp notion.BlockType) *ast.List {
	switch tp {
	case notion.BlockTypeNumberedListItem:
		n := ast.NewList('.')
		n.SetAttributeString(attrClass, classNumberedList)
		return n
	case notion.BlockTypeBulletedListItem:
		n := ast.NewList('-')
		n.SetAttributeString(attrClass, classBulletedList)
		return n
	default:
		n := ast.NewList('-')
		n.SetAttributeString(attrClass, classToDoList)
		return n
	}
}

type listType int

func (c *blockCollector) collectBlock(b notion.Block) error {
	n, err := c.toNodeWithChildren(b)
	if err != nil {
		return err
	}

	switch b.Type {
	case notion.BlockTypeNumberedListItem, notion.BlockTypeBulletedListItem,
		notion.BlockTypeToDo, notion.BlockTypeToggle:
		setParentChild(c.getList(b.Type), n)
	default:
		// non-list, so finish existing list
		if c.list != nil {
			c.res = append(c.res, c.list)
			c.list = nil
			c.listType = ""
		}

		c.res = append(c.res, n)
	}

	return nil
}

func (c *blockCollector) toNodeWithChildren(b notion.Block) (ast.Node, error) {
	n := c.toNode(b)

	switch b.Type {
	case notion.BlockTypeChildPage:
		return n, nil
	case notion.BlockTypeChildDatabase:
		table, err := c.p.getTable(notion.Id(b.Id))
		if err != nil {
			return nil, err
		}

		n.AppendChild(n, table)

		return n, nil
	}

	if !b.HasChildren {
		return c.toNode(b), nil
	}

	bc := &n_ast.BlockChildren{}

	switch n.Kind() {
	case n_ast.KindCallout:
		// append to callout text node instead of the callout node
		text := n.LastChild()
		text.AppendChild(text, bc)
	default:
		n.AppendChild(n, bc)
	}

	children, err := c.p.getBlocks(notion.Id(b.Id), -1)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		bc.AppendChild(bc, child)
	}

	return n, nil
}

func (c *blockCollector) toNode(b notion.Block) ast.Node {
	switch b.Type {
	case notion.BlockTypeParagraph:
		return toNodeParagraph(ast.NewParagraph(), b.Id, b.Paragraph)
	case notion.BlockTypeHeading1:
		return toNodeParagraph(ast.NewHeading(1), b.Id, b.Heading1)
	case notion.BlockTypeHeading2:
		return toNodeParagraph(ast.NewHeading(2), b.Id, b.Heading2)
	case notion.BlockTypeHeading3:
		return toNodeParagraph(ast.NewHeading(3), b.Id, b.Heading3)
	case notion.BlockTypeCallout:
		return toNodeCallout(b.Id, b.Callout)
	case notion.BlockTypeQuote:
		return toNodeParagraph(ast.NewBlockquote(), b.Id, b.Quote)
	case notion.BlockTypeSyncedBlock:
		return n_ast.NewSyncedBlock()
	case notion.BlockTypeToDo:
		return toNodeToDo(b.Id, b.ToDo)
	case notion.BlockTypeNumberedListItem:
		return toNodeListItem(b.Id, b.NumberedListItem)
	case notion.BlockTypeBulletedListItem:
		return toNodeListItem(b.Id, b.BulletedListItem)
	case notion.BlockTypeToggle:
		return toNodeToggle(b.Id, b.Toggle)
	case notion.BlockTypeCode:
		return toNodeCode(b.Id, b.Code)
	case notion.BlockTypeChildPage:
		return c.p.toNodeChildPage(b.Id, b.ChildPage)
	case notion.BlockTypeChildDatabase:
		return toNodeChildDatabase(b.Id, b.ChildDatabase)
	case notion.BlockTypeEmbed:
		return c.p.toNodeEmbed(b.Id, b.Embed.Url, &b.Embed.Caption)
	case notion.BlockTypePdf:
		return c.p.toNodeEmbed(b.Id, b.Pdf.URL(), b.Pdf.Caption)

	// 	// TODO validate:
	// case notion.BlockTypeBookmark:
	// 	return toNodeBookmark(b.Bookmark)

	// case notion.BlockTypeTable:
	// 	// NOTE: toNode should never be called with notion.BlockTypeTableRow
	// 	// the below function will call the appropriate methods
	// 	return toNodeTable(b.Table)
	// case notion.BlockTypeDivider:
	// 	return ast.NewThematicBreak()

	// case notion.BlockTypeEquation:
	// 	return toNodeEquation(b.Equation)
	// case notion.BlockTypeFile:
	// 	return n_ast.NewFile(*b.File, n_ast.FileTypeGeneric)

	// case notion.BlockTypeImage:
	// 	return n_ast.NewFile(*b.File, n_ast.FileTypeImage)
	// case notion.BlockTypeLinkPreview:
	// 	return toNodeLinkPreview(b.LinkPreview)
	// case notion.BlockTypeLinkToPage:
	// 	return n_ast.NewLinkToPage(*b.LinkToPage)

	// case notion.BlockTypeTableOfContents:
	// 	return toNodeTableOfContents(b.TableOfContents)

	// case notion.BlockTypeVideo:
	// 	return toNodeVideo(b.Video)
	default: // includes
		// notion.BlockTypeUnsupported (which we'll never support by its nature)
		// notion.BlockTypeColumn, notion.BlockTypeColumnList (which we plan to support in the future)
		// notion.BlockTypeTemplate (we're unsure when this is returned)
		panic(fmt.Sprintf("unknown node notion block type %q", b.Type))
	}
}

func toNodeParagraph(n ast.Node, id notion.UUID, p *notion.Paragraph) ast.Node {
	n.SetAttributeString(attrID, []byte(id))
	setClasses(n, p.Color)
	appendRichTexts(n, p.RichText)

	return n
}

func toNodeCallout(id notion.UUID, callout *notion.Callout) ast.Node {
	n := &n_ast.Callout{}

	setClasses(n, callout.Color, []byte("callout"))
	n.SetAttributeString("style", []byte("white-space:pre-wrap;display:flex"))
	n.SetAttributeString("id", []byte(id))

	n.AppendChild(n, n_ast.NewIcon(callout.Icon))

	text := &n_ast.CalloutText{}
	n.AppendChild(n, text)

	appendRichTexts(text, callout.RichText)

	return n
}

func toNodeToDo(id notion.UUID, todo *notion.ToDo) ast.Node {
	checkBox := newCheckbox(todo.Checked, todo.Color)

	txt := &n_ast.CheckboxText{Checked: todo.Checked}
	appendRichTexts(txt, todo.RichText)
	setClasses(txt, "", func() []byte {
		if todo.Checked {
			return classCheckboxTextChecked
		}

		return classCheckboxTextUnchecked
	}())

	n := ast.NewListItem(0)
	n.SetAttributeString(attrID, []byte(id))

	n.AppendChild(n, checkBox)
	n.AppendChild(n, ast.NewString([]byte{' '}))
	n.AppendChild(n, txt)

	return n
}

func newCheckbox(checked bool, c notion.Color) *extast.TaskCheckBox {
	mode := classCheckboxOff
	if checked {
		mode = classCheckboxOn
	}

	n := extast.NewTaskCheckBox(checked)
	setClasses(n, c, classCheckbox, mode)

	return n
}

func toNodeListItem(id notion.UUID, item *notion.Paragraph) ast.Node {
	// list := newList(ordered)
	// class := classBulletList
	// if ordered {
	// 	class = classNumberedList
	// }

	// setClasses(list, color, class)

	// if ordered {
	// 	list.SetAttributeString("start", []byte(fmt.Sprintf("%d", c.list.ChildCount())))
	// }

	// return list

	n := ast.NewListItem(0)
	n.SetAttributeString(attrID, []byte(id))
	appendRichTexts(n, item.RichText)
	setClasses(n, item.Color)
	return n
}

func toNodeToggle(id notion.UUID, t *notion.Paragraph) ast.Node {
	n := &n_ast.Toggle{}

	n.SetAttributeString(attrID, []byte(id))
	setClasses(n, t.Color, classToggle)

	txt := &n_ast.ToggleText{}
	appendRichTexts(txt, t.RichText)
	n.AppendChild(n, txt)

	return n
}

func toNodeCode(id notion.UUID, c *notion.Code) ast.Node {
	lang := ast.NewText() // TODO use source to create language

	n := ast.NewFencedCodeBlock(lang)

	n.SetAttributeString(attrID, []byte(id))
	setClasses(n, "", classCode)

	if c.Language != "" {
		n.SetAttributeString("language", []byte(c.Language))
	}

	appendRichTexts(n, c.RichText)

	if c.Caption != nil && len(*c.Caption) > 0 {
		cap := &n_ast.Caption{}
		appendRichTexts(cap, *c.Caption)
		n.AppendChild(n, cap)
	}

	return n
}

func (c *pageCollector) toNodeChildPage(id notion.UUID, child *notion.Child) ast.Node {
	n := n_ast.NewChildPage(*child)

	n.SetAttributeString(attrID, []byte(id))
	setClasses(n, "", classLinkToPage)

	n.AppendChild(n, linkToPage(child.Title, id, c.root))

	return n
}

func linkToPage(title string, id notion.UUID, dir ...string) *ast.Link {
	n := ast.NewLink()

	if title == "" {
		title = "Untitled"
	}

	fileName := fmt.Sprintf("%s %s.html", title,
		strings.ReplaceAll(string(id), "-", ""))

	n.Destination = util.URLEscape([]byte(filepath.Join(append(dir, fileName)...)), true)

	n.AppendChild(n, newString(title))

	return n
}

func toNodeChildDatabase(id notion.UUID, db *notion.Child) ast.Node {
	n := &n_ast.ChildDatabase{}
	n.SetAttributeString(attrID, []byte(id))
	setClasses(n, "", classCollectionContent)

	title := ast.NewHeading(4)
	n.AppendChild(n, title)

	setClasses(title, "", classCollectionTitle)
	title.AppendChild(title, ast.NewString([]byte(db.Title)))

	return n
}

func (p *pageCollector) toNodeEmbed(id notion.UUID, rawURL string, caption *notion.RichTexts) ast.Node {
	n := &n_ast.Embed{}
	n.SetAttributeString(attrID, []byte(id))

	src := &n_ast.EmbedSource{}
	n.AppendChild(n, src)

	link := ast.NewLink()
	src.AppendChild(src, link)

	u, _ := url.Parse(rawURL) // TODO don't underscore error

	fileName := filepath.Base(u.Path)

	if u.Host == "s3.us-west-2.amazonaws.com" &&
		strings.HasPrefix(u.Path, "/secure.notion-static.com/") {
		// TODO download

		link.Destination = []byte(filepath.Join(p.root, fileName))
		u.RawQuery = ""
	} else {
		link.Destination = []byte(rawURL)
	}

	link.AppendChild(link, newString(u.String()))

	if caption == nil || len(*caption) == 0 {
		return n
	}

	capt := &n_ast.Caption{}
	n.AppendChild(n, capt)

	appendRichTexts(capt, *caption)

	return n
}
