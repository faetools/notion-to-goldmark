package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindCallout is a ast.NodeKind of the Callout node.
var KindCallout = ast.NewNodeKind("Callout")

// A Callout represents a callout in Notion.
type Callout struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *Callout) Kind() ast.NodeKind { return KindCallout }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Callout) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
