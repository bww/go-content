package text

import (
	"testing"

	"github.com/bww/go-content/v1"
	"github.com/bww/go-mime/v1"
	"github.com/stretchr/testify/assert"
)

func TestFromHTML(t *testing.T) {
	tests := []struct {
		In     *content.Content
		Expect string
		Err    func(error) error
	}{
		{
			In:     content.NewText("Hello"),
			Expect: "Hello",
		},
		{
			In:     content.NewMarkdown("Hello"),
			Expect: "Hello",
		},
		{
			In:     content.New(mime.HTML, "<b>Hello</b>"),
			Expect: "Hello",
		},
		{
			In:     content.New(mime.HTML, "<html><body><b>Hello</b></body></html>"),
			Expect: "Hello",
		},
		{
			In:     content.New(mime.HTML, "<html> <body> <b>Hello</b> </body> </html>"),
			Expect: "Hello",
		},
		{
			In:     content.New(mime.HTML, "<html> <body> <b>Hello, there.</b> </body> Yep. </html>"),
			Expect: "Hello, there. Yep.",
		},
	}
	for i, e := range tests {
		v, err := From(e.In)
		if e.Err != nil {
			assert.NoError(t, e.Err(err))
		} else if assert.NoError(t, err) {
			assert.Equal(t, e.Expect, v.Data, "#%d", i)
		}
	}
}
