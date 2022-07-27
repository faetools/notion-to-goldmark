package markdown

import (
	"bytes"
	"io"
	"strings"

	"github.com/faetools/format/golang"
	"github.com/faetools/format/writers"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

const (
	// Bytes.
	bNewLine           = '\n'
	bHash              = '#'
	bTab               = '\t'
	bSpace             = ' '
	bLeftSquareBracket = '['
	bRightParenthesis  = ')'
	bGreaterThan       = '>'
	bLessThan          = '<'
	bQuotationMark     = '"'
	bAsterisk          = '*'
	bGraveAccent       = '`'
	bMinus             = '-'
	bTerminalControl   = '\x1b'
	bM                 = 'm'

	// Strings.
	sNewLine           = "\n"
	preOListItemFormat = "%d. "

	langHTML = "html"
	langGo   = "go"
)

var (
	// Byte arrays.
	arrayEmptySpace = []byte{bSpace}
	arrayTab        = []byte{bTab}
	newLine         = []byte{bNewLine}
	twoNewLines     = []byte{bNewLine, bNewLine}
	threeSpaces     = []byte{bSpace, bSpace, bSpace}
	preUListItem    = []byte{bMinus, bSpace}
	twoAsterisks    = []byte{bAsterisk, bAsterisk}
	thematicBreak   = []byte{bMinus, bMinus, bMinus}
	codeBlockFence  = []byte{bGraveAccent, bGraveAccent, bGraveAccent}
	imageStart      = []byte{'!', bLeftSquareBracket}
	linkTitleStart  = []byte{bSpace, bQuotationMark}
	linkTransition  = []byte{']', '('}

	// Terminal.
	tReset  = []byte{bTerminalControl, bLeftSquareBracket, '0', bM}
	tBold   = []byte{bTerminalControl, bLeftSquareBracket, '1', bM}
	tItalic = []byte{bTerminalControl, bLeftSquareBracket, '3', bM}
	// TUnderline     = []byte{bTerminalControl, bLeftSquareBracket, '4', bM}
	// tStrikethrough = []byte{bTerminalControl, bLeftSquareBracket, '9', bM}.
)

func newBlockquoteWriter(w io.Writer) writers.Writer {
	return writers.NewTrimWriter(writers.NewBlockquoteWriter(w), sNewLine)
}

func repeat(w util.BufWriter, c byte, count int) {
	for i := 0; i < count; i++ {
		_ = w.WriteByte(c)
	}
}

func hasSpacePrefix(lines [][]byte) bool {
	if len(lines) == 0 {
		return false
	}

	for _, line := range lines {
		switch line[0] {
		case bSpace, bTab:
		default:
			return false
		}
	}

	return true
}

func trimSpacePrefix(lines [][]byte) [][]byte {
	for i, line := range lines {
		switch line[0] {
		case bSpace:
			lines[i] = line[1:]
		case bTab:
			lines[i] = append(threeSpaces, line[1:]...)
		default:
			panic("trimSpacePrefix called where one line doesn't have space prefix")
		}
	}

	return lines
}

func writeLines(w io.Writer, source []byte, n ast.Node) {
	for i, l := 0, n.Lines().Len(); i < l; i++ {
		line := n.Lines().At(i)
		_, _ = w.Write(line.Value(source))
	}
}

func writeFormattedLines(w util.BufWriter, source []byte, n ast.Node, lang string) {
	defer func() { _ = w.WriteByte(bNewLine) }()

	trimWriter := writers.NewTrimWriter(w, sNewLine)
	if len(lang) == 0 {
		// No formatting necessary.
		writeLines(trimWriter, source, n)
		return
	}

	b := &bytes.Buffer{}
	writeLines(b, source, n)
	_, _ = trimWriter.Write(formatCode(lang, b.Bytes()))
}

func formatCode(lang string, src []byte) []byte {
	switch strings.ToLower(lang) {
	case langGo:
		gofmt, err := golang.Format("", src)
		if err == nil {
			return gofmt
		}
	}

	return src
}

func getIndentOfListItem(node ast.Node) int {
	indent := 0

	for {
		if node = node.Parent(); node == nil {
			return indent / 2 //nolint:gomnd // halfing
		}

		switch node.Kind() {
		case ast.KindList, ast.KindListItem:
			indent++
		default:
			return indent / 2 //nolint:gomnd // halfing
		}
	}
}

// renderChildren renders the children of a node with the exact same options but with a different writer.
func (r *NodeRenderer) renderChildren(w io.Writer, source []byte, n ast.Node) error {
	fullRenderer := renderer.NewRenderer(r.getOptions()...)

	for c := n.FirstChild(); c != nil; {
		if err := fullRenderer.Render(w, source, c); err != nil {
			return errors.Wrap(err, "rendering children")
		}

		c = c.NextSibling()
	}

	return nil
}

func getLanguage(source []byte, node ast.Node) []byte {
	if n, ok := node.(*ast.FencedCodeBlock); ok {
		return n.Language(source)
	}

	return nil
}

func numElement(node ast.Node) (i int) {
	for s := node; s != nil; i++ {
		s = s.PreviousSibling()
	}

	return
}
