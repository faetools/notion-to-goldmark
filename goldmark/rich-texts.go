package goldmark

import (
	"fmt"

	n_ast "github.com/faetools/notion-to-goldmark/ast"

	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

// toNodeRichTextWithoutAnnotations returns the node without annotation information
func toNodeRichTextWithoutAnnotations(t notion.RichText) ast.Node {
	// NOTE: we are ignoring t.Href for now because it seems to be just
	// a repeat of t.Text.Link.Url, t.Mention.LinkPreview.Url etc.

	switch t.Type {
	case notion.RichTextTypeText:
		return toNodeText(t.Text)
	case notion.RichTextTypeEquation:
		return toNodeEquation(t.Equation)
	case notion.RichTextTypeMention:
		return &n_ast.Mention{Content: t.Mention}
	default:
		panic(fmt.Sprintf("invalid RichText of type %q", t.Type))
	}
}

// we want to wrap the nodes
// - n1,n2,n3->n1,n2,n3
// - n1(bold,italic),n2(bold),n3(bold,italic)->bold{italic{n1},n2,italic{n3}}
// - n1(bold,italic),n2(italic),n3(bold,italic)->italic{bold{n1},n2,bold{n3}}
// - n1(bold),n2(bold,italic),n3(italic)->bold{n1,italic{n2}},italic{n3}

// TODO we might have to work with empty spaces:
// - trim them and add back in without formatting

// annotations with extra methods
type annotations notion.Annotations

// node returns the first node representing an annotation
// this should only be called on annotations that only have annotation
func (a annotations) node() ast.Node {
	switch {
	case a.Bold:
		return ast.NewEmphasis(2)
	case a.Italic:
		return ast.NewEmphasis(1)
	case a.Underline:
		return ast.NewEmphasis(3)
	case a.Strikethrough:
		return extast.NewStrikethrough()
	case a.Code:
		return ast.NewCodeSpan()
	case a.Color != notion.ColorDefault && a.Color != "":
		return &n_ast.Color{Color: a.Color}
	default:
		return nil
	}
}

type annotationWrappers []*annotationWrapper

type annotationWrapper struct {
	// represents first all annotations
	// then objects will be wrapped so in the end it's just one annotation
	ann annotations

	// a node that is wrapped
	// should only be present if there are no annotations or sub wrappers
	node ast.Node

	// the objects that this object wraps
	subs annotationWrappers
}

func toNodeRichTexts(rts notion.RichTexts) (ns []ast.Node) {
	// construct the initial wrappers
	ws := make(annotationWrappers, len(rts))

	for i, rt := range rts {
		n := toNodeRichTextWithoutAnnotations(rt)
		ann := annotations(rt.Annotations)

		if ann.node() == nil {
			// pure node
			ws[i] = &annotationWrapper{node: n}
		} else {
			// a wrapper with all annotations, wrapping the node
			ws[i] = &annotationWrapper{ann: ann, subs: annotationWrappers{{node: n}}}
		}
	}

	// merge siblings that have the same annotation
	ws = ws.mergeSiblings()

	// wrap wrappers until each wrapper has either a node or just one annotation
	ws.wrapSelves()

	return ws.toNodes()
}

// mergeIfSameAnnotation merges two wrappers if they have the same annotation
func mergeIfSameAnnotation(w1, w2 *annotationWrapper) *annotationWrapper {
	ann := annotations{Color: notion.ColorDefault}

	for _, v := range []struct {
		this, next, merged *bool
	}{
		{&w1.ann.Bold, &w2.ann.Bold, &ann.Bold},
		{&w1.ann.Italic, &w2.ann.Italic, &ann.Italic},
		{&w1.ann.Code, &w2.ann.Code, &ann.Code},
		{&w1.ann.Strikethrough, &w2.ann.Strikethrough, &ann.Strikethrough},
		{&w1.ann.Underline, &w2.ann.Underline, &ann.Underline},
	} {
		if *v.this && *v.next {
			// set the flag for the new wrapper
			*v.merged = true

			// disable the flag for the old wrappers
			*v.this = false
			*v.next = false

			return &annotationWrapper{
				ann:  ann,
				subs: merge(w1, w2),
			}
		}
	}

	switch w1.ann.Color {
	case notion.ColorDefault, "": // no color, ignore
	case w2.ann.Color:
		w1.ann.Color = notion.ColorDefault
		w2.ann.Color = notion.ColorDefault

		return &annotationWrapper{
			ann:  ann,
			subs: merge(w1, w2),
		}
	}

	return nil
}

// merge merges two wrappers so they can form the subwrappers to a new wrapper
func merge(ws ...*annotationWrapper) annotationWrappers {
	// make sure all wrappers contain either node or annotation
	subs := make(annotationWrappers, 0, len(ws))

	for _, w := range ws {
		if w.ann.node() != nil {
			// w has annotation information
			subs = append(subs, w)
			continue
		}

		// if len(w.subs) == 0 || w.node != nil {
		// 	panic("something went wrong")
		// }

		// w is empty wrapper
		subs = append(subs, w.subs...)
	}

	return subs
}

// remove removes an element from the slice or array.
func remove[T any](collection []T, i int) []T {
	copy(collection[i:], collection[i+1:])

	end := len(collection) - 1

	var zero T
	collection[end] = zero

	return collection[:end]
}

// mergeSiblings merges siblings if they have the same annotations
func (ws annotationWrappers) mergeSiblings() annotationWrappers {
	for i, l := 0, len(ws)-1; i < l; i++ {
		if merger := mergeIfSameAnnotation(ws[i], ws[i+1]); merger != nil {
			ws[i] = merger       // write the merger
			ws = remove(ws, i+1) // remove i+1

			return ws.mergeSiblings()
		}
	}

	return ws
}

// wrapSelves calls wrapSelf on all elements
func (ws annotationWrappers) wrapSelves() {
	for _, w := range ws {
		w.wrapSelf()
	}
}

func (w *annotationWrapper) wrap(child *annotationWrapper) {
	if w.node == nil && w.ann.node() == nil {
		// parent would be empty wrapper
		w.ann = child.ann
		return
	}

	child.subs = w.subs
	w.subs = annotationWrappers{child}
}

// wrapSelf wraps this wrapper in another wrapper for each annotation
// so that each wrapper only has one annotation
func (w *annotationWrapper) wrapSelf() {
	w.subs.wrapSelves()

	if w.ann.Bold {
		w.ann.Bold = false
		w.wrap(&annotationWrapper{
			ann: annotations{Bold: true},
		})
	}

	if w.ann.Italic {
		w.ann.Italic = false
		w.wrap(&annotationWrapper{
			ann: annotations{Italic: true},
		})
	}

	if w.ann.Code {
		w.ann.Code = false
		w.wrap(&annotationWrapper{
			ann: annotations{Code: true},
		})
	}

	if w.ann.Strikethrough {
		w.ann.Strikethrough = false
		w.wrap(&annotationWrapper{
			ann: annotations{Strikethrough: true},
		})
	}

	if w.ann.Underline {
		w.ann.Underline = false
		w.wrap(&annotationWrapper{
			ann: annotations{Underline: true},
		})
	}

	if w.ann.Color != notion.ColorDefault &&
		w.ann.Color != "" {
		w.ann.Color = notion.ColorDefault

		w.wrap(&annotationWrapper{
			ann: annotations{Color: w.ann.Color},
		})
	}
}

func (ws annotationWrappers) toNodes() []ast.Node {
	out := make([]ast.Node, len(ws))

	for i, wr := range ws {
		out[i] = wr.toNode()
	}

	return out
}

func (w *annotationWrapper) toNode() ast.Node {
	if w.node != nil {
		return w.node
	}

	n := w.ann.node()

	doPanic := n == nil

	for _, sub := range w.subs {
		if doPanic {
			// TODO
			sub.toNode().Dump(nil, 0)
			continue
		}

		n.AppendChild(n, sub.toNode())
	}

	if doPanic {
		panic("aaaaaa!")
	}

	return n
}
