package ast

import "github.com/yuin/goldmark/ast"

// KindPropertyIcon is a ast.NodeKind of the PropertyIcon node.
var KindPropertyIcon = ast.NewNodeKind("PropertyIcon")

// A PropertyIcon represents a property icon in Notion.
type PropertyIcon struct{ ast.BaseInline }

// Kind returns a kind of this node.
func (n *PropertyIcon) Kind() ast.NodeKind { return KindPropertyIcon }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *PropertyIcon) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
