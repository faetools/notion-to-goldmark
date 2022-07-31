package ast

import (
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindDate is a ast.NodeKind of the Date node.
var KindDate = ast.NewNodeKind("Date")

// A Date represents a callout in Notion.
type Date struct {
	ast.BaseInline
	Date            *notion.Date
	TwelveHourClock bool
}

func NewDate(date *notion.Date, twelveHourClock bool) *Date {
	return &Date{
		Date:            date,
		TwelveHourClock: twelveHourClock,
	}
}

// Kind returns a kind of this node.
func (n *Date) Kind() ast.NodeKind { return KindDate }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Date) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
