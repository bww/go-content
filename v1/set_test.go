package content

import (
	"testing"

	"github.com/bww/go-mime/v1"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	tests := []struct {
		Set    Set
		Choose mime.Type
		Expect *Content
	}{
		{
			Set:    nil,
			Choose: mime.Markdown,
			Expect: nil,
		},
		{
			Set:    Set{},
			Choose: mime.Markdown,
			Expect: nil,
		},
		{
			Set: Set{
				&Content{Type: "", Data: ""},
			},
			Choose: mime.Markdown,
			Expect: &Content{Type: "", Data: ""},
		},
		{
			Set: Set{
				&Content{Type: mime.Text, Data: "Hello"},
				&Content{Type: mime.Markdown, Data: "_Hello_"},
			},
			Choose: mime.Text,
			Expect: &Content{Type: mime.Text, Data: "Hello"},
		},
		{
			Set: Set{
				&Content{Type: mime.Text, Data: "Hello"},
				&Content{Type: mime.Markdown, Data: "_Hello_"},
				&Content{Type: mime.Markdown, Data: "_Hello, again_"}, // this is effectively unreachable
			},
			Choose: mime.Markdown,
			Expect: &Content{Type: mime.Markdown, Data: "_Hello_"},
		},
		{
			Set: Set{
				&Content{Type: mime.Text, Data: "Hello"},
				&Content{Type: mime.Markdown, Data: "_Hello_"},
			},
			Choose: mime.JSON,
			Expect: &Content{Type: mime.Text, Data: "Hello"},
		},
	}
	for i, e := range tests {
		alt := e.Set.Choose(e.Choose)
		assert.Equal(t, e.Expect, alt, "#%d", i)
	}
}
