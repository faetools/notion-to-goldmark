package ast

import "github.com/yuin/goldmark/ast"

// KindPolygon is a ast.NodeKind of the Polygon node.
var KindPolygon = ast.NewNodeKind("Polygon")

// A Polygon represents a polygon in Notion.
type Polygon struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *Polygon) Kind() ast.NodeKind { return KindPolygon }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Polygon) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
