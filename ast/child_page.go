package ast

import (
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindChildPage is a ast.NodeKind of the ChildPage node.
var KindChildPage = ast.NewNodeKind("ChildPage")

// A ChildPage represents a child page in Notion.
type ChildPage struct {
	ast.BaseInline
	Page notion.Child
}

// NewChildPage returns a new child page node.
func NewChildPage(p notion.Child) ast.Node {
	return &ChildPage{Page: p}
}

// Kind returns a kind of this node.
func (n *ChildPage) Kind() ast.NodeKind { return KindChildPage }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *ChildPage) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Title": n.Page.Title}, nil)
}
