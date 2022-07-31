package ast

import "github.com/yuin/goldmark/ast"

// KindSVG is a ast.NodeKind of the SVG node.
var KindSVG = ast.NewNodeKind("SVG")

// A SVG represents an svg icon in Notion.
type SVG struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *SVG) Kind() ast.NodeKind { return KindSVG }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *SVG) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
