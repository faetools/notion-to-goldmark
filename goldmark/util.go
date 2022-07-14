package goldmark

import (
	"fmt"

	"github.com/faetools/go-notion/pkg/notion"
	n_ast "github.com/faetools/notion-to-goldmark/ast"
	"github.com/yuin/goldmark/ast"
)

// setParentChild is a convenience method to not name the parent twice.
func setParentChild(parent, child ast.Node) {
	parent.AppendChild(parent, child)
}

func newLink(content, dest string) *ast.Link {
	n := ast.NewLink()
	n.Destination = []byte(dest)

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

func appendCommon(n ast.Node, rts notion.RichTexts, children notion.Blocks) {
	for _, child := range toNodeRichTexts(rts) {
		n.AppendChild(n, child)
	}

	addBlockChildren(n, children)
}

// addBlockChildren adds any children blocks
// TODO: call the API go get the children
func addBlockChildren(n ast.Node, bs notion.Blocks) {
	if len(bs) == 0 {
		return
	}

	children := &n_ast.Children{}
	for _, child := range FromBlocks(bs) {
		children.AppendChild(children, child)
	}

	n.AppendChild(n, children)
}

func addCaption(n ast.Node, rts notion.RichTexts) {
	if len(rts) == 0 {
		return
	}

	caption := &n_ast.Caption{}
	appendCommon(caption, rts, nil)
	setParentChild(n, caption)
}
