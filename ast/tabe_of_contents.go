package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindTableOfContents is a ast.NodeKind of the TableOfContents node.
var KindTableOfContents = ast.NewNodeKind("TableOfContents")

// A TableOfContents represents a child page in Notion.
type TableOfContents struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *TableOfContents) Kind() ast.NodeKind { return KindTableOfContents }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *TableOfContents) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
