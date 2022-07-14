package ast

import (
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindColor is a ast.NodeKind of the Color node.
var KindColor = ast.NewNodeKind("Color")

// A Color represents a color in Notion.
type Color struct {
	ast.BaseInline
	Color notion.Color
}

// Kind returns a kind of this node.
func (n *Color) Kind() ast.NodeKind { return KindColor }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Color) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Color": string(n.Color)}, nil)
}
