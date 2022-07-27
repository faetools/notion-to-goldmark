package markdown

import (
	"bytes"
	"fmt"

	"github.com/faetools/format/writers"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"

	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	_yaml "gopkg.in/yaml.v2"
)

var myParser = goldmark.New(
	goldmark.WithExtensions(meta.Meta),
	goldmark.WithExtensions(extension.GFM)).Parser()

// Format formats Markdown.
func Format(src []byte, options ...renderer.Option) ([]byte, error) {
	ctx := parser.NewContext()
	parsed := myParser.Parse(text.NewReader(src), parser.WithContext(ctx))

	return Render(meta.GetItems(ctx), src, parsed, options...)
}

// Render renders a given parsed document.
func Render(metaData _yaml.MapSlice, src []byte, doc ast.Node, options ...renderer.Option) ([]byte, error) {
	b := &bytes.Buffer{}

	if len(metaData) > 0 {
		b.WriteString("---\n")

		if err := _yaml.NewEncoder(b).Encode(metaData); err != nil {
			return nil, fmt.Errorf("marshalling metadata: %w", err)
		}

		b.WriteString("---\n\n")
	}

	w := writers.NewTrimWriter(b, sNewLine)

	nr := NewNodeRenderer()
	nr.additionalOptions = options

	opts := append([]renderer.Option{
		renderer.WithNodeRenderers(util.Prioritized(nr, 0)),
	}, options...)

	err := renderer.NewRenderer(opts...).Render(w, src, doc)
	if err != nil {
		return nil, errors.Wrap(err, "rendering markdown")
	}

	// Trim all new lines but add one at the end.
	b.Write(newLine)

	return b.Bytes(), nil
}

// A NodeRenderer struct is an implementation of renderer.NodeRenderer that renders
// nodes as markdown.
type NodeRenderer struct {
	Config
	additionalOptions []renderer.Option
}

// NewNodeRenderer returns a new Renderer with given options.
func NewNodeRenderer(opts ...Option) *NodeRenderer {
	r := &NodeRenderer{Config: NewConfig()}

	for _, opt := range opts {
		opt.SetMarkdownOption(&r.Config)
	}

	return r
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs .
func (r *NodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// Blocks.
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)
	reg.Register(extast.KindStrikethrough, r.renderStrikethrough)

	// Inlines.
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)
}

func (r *NodeRenderer) renderHeading(w util.BufWriter,
	source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		_, _ = w.Write(twoNewLines)
		return ast.WalkContinue, nil
	}

	//nolint:forcetypeassert // function is only called with that
	repeat(w, bHash, node.(*ast.Heading).Level)
	_ = w.WriteByte(bSpace)

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderBlockquote(w util.BufWriter,
	source []byte, n ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		_, _ = w.Write(twoNewLines)
		return ast.WalkContinue, nil
	}

	indent := getIndentOfListItem(n)
	if indent > 0 {
		// Separate from the list with a new line.
		_ = w.WriteByte(bNewLine)
	}

	// Write all indents.
	iw := writers.NewIndentWriter(w, indent)

	// Write as blockquote.
	bw := newBlockquoteWriter(iw)

	err := r.renderChildren(bw, source, n)
	if err != nil {
		return 0, err
	}

	return ast.WalkSkipChildren, nil
}

func (r *NodeRenderer) renderCodeBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	indent := getIndentOfListItem(n)
	if indent == 0 {
		// No indent, so we can transform into a fenced code block.
		return r.renderFencedCodeBlock(w, source, n, entering)
	}

	if !entering {
		return ast.WalkContinue, nil
	}

	// Separate from the rest.
	_ = w.WriteByte(bNewLine)

	l := n.Lines().Len()
	lines := make([][]byte, l)

	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		lines[i] = line.Value(source)
	}

	// Shift all lines to the left.
	for hasSpacePrefix(lines) {
		lines = trimSpacePrefix(lines)
	}

	// Write lines with indent.
	iw := writers.NewIndentWriter(w, indent+1)
	for _, line := range lines {
		_, _ = iw.Write(line)
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderFencedCodeBlock(w util.BufWriter,
	source []byte, node ast.Node, entering bool,
) (ast.WalkStatus, error) {
	_, _ = w.Write(codeBlockFence)

	if !entering {
		_, _ = w.Write(twoNewLines)
		return ast.WalkContinue, nil
	}

	lang := getLanguage(source, node)

	// E.g. "```go\n".
	_, _ = w.Write(lang)
	_ = w.WriteByte(bNewLine)

	writeFormattedLines(w, source, node, string(lang))

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		// Write the closure line, if any.
		//nolint:forcetypeassert // function is only called with that
		if cl := node.(*ast.HTMLBlock).ClosureLine; cl.Start > -1 {
			_, _ = w.Write(cl.Value(source))
		}

		_ = w.WriteByte(bNewLine)

		return ast.WalkContinue, nil
	}

	writeFormattedLines(w, source, node, langHTML)

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering && node.Parent().Kind() != ast.KindListItem {
		// Two new lines at the end if not inside another list.
		_, _ = w.Write(twoNewLines)
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	indent := getIndentOfListItem(node)

	if !entering {
		if // There is no next list item.
		node.NextSibling() != nil &&
			// No content.
			(node.LastChild() == nil ||
				// We have not rendered a blockquote.
				node.LastChild().Kind() != ast.KindBlockquote) {

			_ = w.WriteByte(bNewLine)
		}

		return ast.WalkSkipChildren, nil
	}

	indentWriter := writers.NewIndentWriter(w, indent)

	p := node.Parent()
	if !p.(*ast.List).IsOrdered() { //nolint:forcetypeassert // function is only called with that
		_, _ = indentWriter.Write(preUListItem)
	} else {
		_, _ = indentWriter.WriteString(fmt.Sprintf(preOListItemFormat, numElement(node)))
	}

	if err := r.renderChildren(indentWriter, source, node); err != nil {
		return 0, err
	}

	return ast.WalkSkipChildren, nil
}

func (r *NodeRenderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	indent := getIndentOfListItem(n)
	partOfList := indent > 0

	if !entering {
		if partOfList && n.NextSibling() != nil {
			_ = w.WriteByte(bNewLine)
		} else if !partOfList {
			_, _ = w.Write(twoNewLines)
		}

		return ast.WalkContinue, nil
	}

	if partOfList && n.PreviousSibling() != nil {
		_ = w.WriteByte(bNewLine)
		iw := writers.NewIndentWriter(w, indent)

		if err := r.renderChildren(iw, source, n); err != nil {
			return 0, err
		}
		_ = w.WriteByte(bNewLine)
	} else {
		if err := r.renderChildren(w, source, n); err != nil {
			return 0, err
		}
	}

	return ast.WalkSkipChildren, nil
}

func (r *NodeRenderer) renderTextBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering && n.NextSibling() != nil {
		_ = w.WriteByte(bNewLine)
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderThematicBreak(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		_, _ = w.Write(twoNewLines)
		return ast.WalkContinue, nil
	}

	_, _ = w.Write(thematicBreak)

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		_ = w.WriteByte(bGreaterThan)
		return ast.WalkContinue, nil
	}

	_ = w.WriteByte(bLessThan)

	n := node.(*ast.AutoLink)
	_, _ = w.Write(n.URL(source))

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// Check if backticks need to be escaped.
	if c := node.FirstChild(); c != nil &&
		//nolint:forcetypeassert // is always that
		bytes.Contains(c.(*ast.Text).Segment.Value(source), []byte{bGraveAccent}) {

		_ = w.WriteByte(bGraveAccent)
	}

	if !entering {
		_ = w.WriteByte(bGraveAccent)
		return ast.WalkContinue, nil
	}

	_ = w.WriteByte(bGraveAccent)

	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		segment := c.(*ast.Text).Segment //nolint:forcetypeassert // is always that
		_, _ = w.Write(segment.Value(source))
	}

	return ast.WalkSkipChildren, nil
}

const (
	levelItalics = 1
	levelBold    = 2
)

var errInvalidEmphLevel = errors.New("invalid emphasis level")

func (r *NodeRenderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)

	if n.Attributes() != nil {
		panic(n.Attributes())
	}

	if r.Config.Terminal {
		if entering {
			switch n.Level {
			case levelItalics:
				_, _ = w.Write(tItalic)
			case levelBold:
				_, _ = w.Write(tBold)
			}
		} else {
			_, _ = w.Write(tReset)
		}

		return ast.WalkContinue, nil
	}

	switch n.Level {
	case levelItalics:
		_ = w.WriteByte(bAsterisk)
	case levelBold:
		_, _ = w.Write(twoAsterisks)
	default:
		return 0, fmt.Errorf("%w (level %d)", errInvalidEmphLevel, n.Level)
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_ = w.WriteByte(bLeftSquareBracket)
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Link)
	_, _ = w.Write(linkTransition)
	_, _ = w.Write(n.Destination) // Link.

	if len(n.Title) != 0 {
		_, _ = w.Write(linkTitleStart)
		_, _ = w.Write(n.Title)
		_ = w.WriteByte(bQuotationMark)
	}

	_ = w.WriteByte(bRightParenthesis)

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Image)
	_, _ = w.Write(imageStart)
	_, _ = w.Write(n.Text(source)) // Alt.
	_, _ = w.Write(linkTransition)
	_, _ = w.Write(n.Destination) // Link.

	if len(n.Title) != 0 {
		_, _ = w.Write(linkTitleStart)
		_, _ = w.Write(n.Title)
		_ = w.WriteByte(bQuotationMark)
	}

	_ = w.WriteByte(bRightParenthesis)

	return ast.WalkSkipChildren, nil
}

func (r *NodeRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}

	_, _ = w.Write(node.Text(source))

	//nolint:forcetypeassert // function is only called with that
	segs := node.(*ast.RawHTML).Segments
	for i, l := 0, segs.Len(); i < l; i++ {
		seg := segs.At(i)
		_, _ = w.Write(seg.Value(source))
	}

	return ast.WalkSkipChildren, nil
}

func (r *NodeRenderer) renderString(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.Write(n.(*ast.String).Value)
	} else if s := n.NextSibling(); s != nil && s.Kind() == ast.KindList {
		w.WriteString("\n")
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Text)
	_, _ = w.Write(n.Segment.Value(source))

	if n.IsRaw() {
		return ast.WalkContinue, nil
	}

	if n.SoftLineBreak() {
		switch {
		case r.HardWraps:
			_, _ = w.Write(twoNewLines)
		default:
			_ = w.WriteByte(bNewLine)
		}

		return ast.WalkContinue, nil
	}

	if n.HardLineBreak() {
		_, _ = w.Write(twoNewLines)
	}

	return ast.WalkContinue, nil
}

func (r *NodeRenderer) renderStrikethrough(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	return ast.WalkContinue, nil
}
