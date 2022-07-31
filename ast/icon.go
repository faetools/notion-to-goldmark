package ast

import (
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindIcon is a ast.NodeKind of the Icon node.
var KindIcon = ast.NewNodeKind("Icon")

// A Icon represents a icon in Notion.
type Icon struct {
	ast.BaseInline
	Emoji string
}

func NewIcon(icon notion.Icon) *Icon {
	n := &Icon{}

	switch icon.Type {
	case notion.IconTypeEmoji:
		n.Emoji = *icon.Emoji
	case notion.IconTypeExternal, notion.IconTypeFile:
		link := ast.NewLink()
		link.Destination = []byte(icon.URL())
		img := ast.NewImage(link)

		img.SetAttributeString("class", []byte("icon"))

		if icon.Type == notion.IconTypeFile {
			img.SetAttributeString("expires", icon.File.ExpiryTime)
		}

		n.AppendChild(n, img)
	}

	return n
}

// Kind returns a kind of this node.
func (n *Icon) Kind() ast.NodeKind { return KindIcon }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Icon) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
