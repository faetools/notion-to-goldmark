package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindCalloutText is a ast.NodeKind of the CalloutText node.
var KindCalloutText = ast.NewNodeKind("CalloutText")

// A CalloutText represents a callout text in Notion.
type CalloutText struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *CalloutText) Kind() ast.NodeKind { return KindCalloutText }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *CalloutText) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
