package content

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContent(t *testing.T) {
	tests := []struct {
		Content *Content
		Expect  []byte
		Error   error
	}{
		{
			&Content{
				Type: "",
				Data: "",
			},
			[]byte(""),
			nil,
		},
		{
			&Content{
				Type: "hello/there",
				Data: "",
			},
			[]byte("hello/there,"),
			nil,
		},
		{
			&Content{
				Type: "",
				Data: "Cooltown okaysville",
			},
			[]byte(",Q29vbHRvd24gb2theXN2aWxsZQ=="),
			nil,
		},
		{
			&Content{
				Type: "hello/there;charset=utf-8",
				Data: "",
			},
			[]byte("hello/there;charset=utf-8,"),
			nil,
		},
		{
			&Content{
				Type: "hello/there",
				Data: "Cooltown okaysville",
			},
			[]byte("hello/there,Q29vbHRvd24gb2theXN2aWxsZQ=="),
			nil,
		},
		{
			&Content{
				Type: "hello/there;charset=utf-8",
				Data: "Cooltown okaysville",
			},
			[]byte("hello/there;charset=utf-8,Q29vbHRvd24gb2theXN2aWxsZQ=="),
			nil,
		},
		{
			&Content{
				Type: " hello / there; charset=utf-8 ",
				Data: "Cooltown okaysville",
			},
			[]byte(" hello / there; charset=utf-8 ,Q29vbHRvd24gb2theXN2aWxsZQ=="),
			nil,
		},
	}
	for _, e := range tests {
		enc, err := e.Content.MarshalColumn()
		if e.Error != nil {
			fmt.Println("***", err)
			assert.Equal(t, e.Error, err)
		} else if assert.Nil(t, err, fmt.Sprint(err)) {
			fmt.Println("-->", string(enc))
			assert.Equal(t, e.Expect, enc)
		}
		if err != nil {
			dec := &Content{}
			err = dec.UnmarshalColumn(enc)
			if assert.Nil(t, err, fmt.Sprint(err)) {
				assert.Equal(t, e.Content, dec)
			}
		}
	}
}

func TestJoinContent(t *testing.T) {
	tests := []struct {
		Data   []*Content
		Sep    string
		Expect *Content
		Error  func(error) error
	}{
		{
			[]*Content{
				{Type: "", Data: ""},
			},
			"",
			&Content{Type: "", Data: ""},
			nil,
		},
		{
			[]*Content{
				{Type: "a", Data: "First"},
				{Type: "a", Data: "Second"},
			},
			"/",
			&Content{Type: "a", Data: "First/Second"},
			nil,
		},
		{
			[]*Content{
				{Type: "a", Data: "First"},
				{Type: "a", Data: "Second"},
				{Type: "a", Data: "Third"},
			},
			"/",
			&Content{Type: "a", Data: "First/Second/Third"},
			nil,
		},
		{
			[]*Content{
				{Type: "a", Data: "First"},
				{Type: "b", Data: "Second"},
				{Type: "a", Data: "Third"},
			},
			"/",
			&Content{},
			func(err error) error {
				if errors.Is(err, ErrIncompatibleFormats) {
					return nil
				} else {
					return err
				}
			},
		},
	}
	for _, e := range tests {
		res, err := Join(e.Data, e.Sep)
		if e.Error != nil {
			assert.Nil(t, e.Error(err), fmt.Sprint(err))
		} else {
			assert.Equal(t, e.Expect, res)
		}
	}
}
