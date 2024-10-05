package html

import (
	"html"
	"html/template"

	"github.com/bww/go-content/v1"
	"github.com/bww/go-mime/v1"
	"github.com/russross/blackfriday/v2"
)

// HTML is a view of Content as HTML
type HTML content.Content

// Of produces a view of the provided content as HTML
func Of(c *content.Content) *HTML {
	return (*HTML)(c)
}

// HTML is a convenience method that is intended to be useful in templates.
// The content is converted to HTML if necessary, sanitized, and returned
// as a [template.HTML] so it will be rendered verbatim in a template.
func (d *HTML) HTML() (template.HTML, error) {
	h, err := From((*content.Content)(d))
	if err != nil {
		return "", err
	}
	return template.HTML(h.Data), nil
}

// Convert converts the provided content to an HTML representation, if possible.
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
