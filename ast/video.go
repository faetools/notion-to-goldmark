package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindVideo is a ast.NodeKind of the Video node.
var KindVideo = ast.NewNodeKind("Video")

// A Video represents a video in Notion.
type Video struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *Video) Kind() ast.NodeKind { return KindVideo }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Video) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
