package html

import (
	"html"

	"github.com/bww/go-content/v1"
	"github.com/bww/go-mime/v1"
	"github.com/russross/blackfriday/v2"
)

// Convert converts the provided content to an HTML representation,
// if possible.
func From(c *content.Content) (*content.Content, error) {
	if c == nil {
		return nil, nil
	}
	switch c.Type {
	case mime.HTML:
		return content.New(mime.HTML, Sanitize(string(c.Data))), nil
	case mime.Text:
		return fromText(c)
	case mime.Markdown:
		return fromMarkdown(c)
	default:
		return nil, content.ErrIncompatibleFormats
	}
}

func fromText(c *content.Content) (*content.Content, error) {
	return content.New(
		mime.HTML,
		"<p>"+html.EscapeString(c.Data)+"</p>",
	), nil
}
func fromMarkdown(c *content.Content) (*content.Content, error) {
	return content.New(
		mime.HTML,
		Sanitize(string(blackfriday.Run([]byte(c.Data)))),
	), nil
}
