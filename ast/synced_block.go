package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindSyncedBlock is a ast.NodeKind of the SyncedBlock node.
var KindSyncedBlock = ast.NewNodeKind("SyncedBlock")

// A SyncedBlock represents a syned block in Notion.
type SyncedBlock struct {
	ast.BaseInline
}

// NewSyncedBlock returns a new syned block node.
func NewSyncedBlock() ast.Node {
	return &SyncedBlock{}
}

// Kind returns a kind of this node.
func (n *SyncedBlock) Kind() ast.NodeKind { return KindSyncedBlock }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *SyncedBlock) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
