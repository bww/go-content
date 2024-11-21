package text

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bww/go-content/v1"
	"github.com/bww/go-mime/v1"
	"golang.org/x/net/html"
)

// Text is a view of Content as plain text
type Text content.Content

// Of produces a view of the provided content as plain text. This does not
// convert the content, it must already be of this type.
func Of(c *content.Content) *Text {
	return (*Text)(c)
}

// Convert converts the provided content to a plain text representation, if
// possible.
func From(c *content.Content) (*content.Content, error) {
	if c == nil {
		return nil, nil
	}
	switch c.Type {
	case mime.Text:
		return c, nil
	case mime.Markdown:
		return content.NewText(c.Data), nil
	case mime.HTML:
		return fromHTML(c.Data)
	default:
		return nil, fmt.Errorf("%w (%s -> %s)", content.ErrIncompatibleFormats, c.Type, mime.HTML)
	}
}

func fromHTML(s string) (*content.Content, error) {
	v, err := fromHTMLText(s)
	if err != nil {
		return nil, fmt.Errorf("Conversion failed: %w", err)
	}
	return content.NewText(v), nil
}

func fromHTMLText(s string) (string, error) {
	doc, err := html.Parse(bytes.NewBufferString(s))
	if err != nil {
		return "", err
	}
	buf := &strings.Builder{}
	err = fromHTMLNode(doc, buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func fromHTMLNode(n *html.Node, buf *strings.Builder) error {
	if n.Type == html.TextNode {
		t := strings.TrimSpace(n.Data)
		if len(t) > 0 {
			if buf.Len() > 0 {
				buf.WriteRune(' ')
			}
			buf.WriteString(t)
		}
	}
	return fromHTMLNodeRec(n, buf)
}

func fromHTMLNodeRec(n *html.Node, buf *strings.Builder) error {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := fromHTMLNode(c, buf)
		if err != nil {
			return err
		}
	}
	return nil
}
