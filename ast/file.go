package ast

import (
	"strconv"
	"time"

	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindFile is a ast.NodeKind of the File node.
var KindFile = ast.NewNodeKind("File")

// A File represents a file in Notion.
type File struct {
	ast.BaseInline
	Expires  time.Time
	External bool
	tp       FileType
}

// FileType is a type of file object.
type FileType string

const (
	// FileTypePDF is a type for a file that is a PDF.
	FileTypePDF FileType = "pdf"
	// FileTypeImage is a type for a file that is an image.
	FileTypeImage FileType = "image"
	// FileTypeVideo is a type for a file that is a video.
	FileTypeVideo FileType = "video"
	// FileTypeGeneric is a type for any other file.
	FileTypeGeneric FileType = "generic"
)

// NewFile returns a new file node
func NewFile(f notion.File, tp FileType) ast.Node {
	n := &File{}
	switch f.Type {
	case notion.FileTypeFile:
		n.Expires = f.File.ExpiryTime
	case notion.FileTypeExternal:
		n.External = true
	}

	link := ast.NewLink()
	link.Destination = []byte(f.URL())

	if tp == FileTypeImage {
		n.AppendChild(n, ast.NewImage(link))
	} else {
		n.AppendChild(n, link)
	}

	return n
}

// Destination is a convenience method to return the destination of the underlying link.
func (n File) Destination() []byte {
	switch c := n.FirstChild().(type) {
	case *ast.Link:
		return c.Destination
	case *ast.Image:
		return c.Destination
	default:
		return nil
	}
}

// IsImage returns whether or not this file is an image.
func (n File) IsImage() bool { return n.tp == FileTypeImage }

// IsPDF returns whether or not this file is a PDF.
func (n File) IsPDF() bool { return n.tp == FileTypePDF }

// IsVideo returns whether or not this file is a video.
func (n File) IsVideo() bool { return n.tp == FileTypeVideo }

// Kind returns a kind of this node.
func (n *File) Kind() ast.NodeKind { return KindFile }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *File) Dump(source []byte, level int) {
	kv := map[string]string{"External": strconv.FormatBool(n.External)}
	if !n.Expires.IsZero() {
		kv["Expires"] = n.Expires.String()
	}

	ast.DumpHelper(n, source, level, kv, nil)
}
