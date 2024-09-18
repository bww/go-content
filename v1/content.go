package content

import (
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/bww/go-mime/v1"
	textutil "github.com/bww/go-util/v1/text"
)

var ErrIncompatibleFormats = errors.New("Content formats are not compatible")

var errContentTypeRequired = errors.New("Content type is required")

type Content struct {
	Type mime.Type `json:"type"`
	Data string    `json:"data"`
}

func NewText(d string) *Content {
	return New(mime.Text, d)
}

func NewMarkdown(d string) *Content {
	return New(mime.Markdown, d)
}

func New(t mime.Type, d string) *Content {
	return &Content{
		Type: t,
		Data: d,
	}
}

func Join(data []*Content, sep string) (*Content, error) {
	if len(data) == 0 {
		return nil, nil
	}
	var res Content
	for i, e := range data {
		if i == 0 {
			res.Type = e.Type
		} else if e.Type != res.Type {
			return nil, ErrIncompatibleFormats
		} else {
			res.Data += sep
		}
		res.Data += e.Data
	}
	return &res, nil
}

func (d *Content) Validate() error {
	if len(d.Data) > 0 && len(d.Type) < 1 { // if we have data but no type, oh, you better believe that's an error
		return errContentTypeRequired
	} else {
		return nil
	}
}

func (d *Content) IsZero() bool {
	return len(d.Data) == 0
}

func (d *Content) String() string {
	if d == nil {
		return "<nil>"
	}
	switch d.Type {
	case mime.Text, mime.Markdown:
		return d.Data
	default:
		v, err := d.MarshalColumn()
		if err != nil {
			return fmt.Sprintf("<could not encode: %v>", err)
		} else {
			return string(v)
		}
	}
}

func (d *Content) MarshalColumn() ([]byte, error) {
	if d == nil || d.Type == "" && len(d.Data) < 1 {
		return []byte{}, nil
	}
	b := &strings.Builder{}
	b.WriteString(d.Type.String())
	b.WriteString(",")
	if d.Data != "" {
		b.WriteString(base64.StdEncoding.EncodeToString([]byte(d.Data)))
	}
	return []byte(b.String()), nil
}

func (d *Content) UnmarshalColumn(text []byte) error {
	if len(text) < 1 {
		return nil
	}
	x := strings.Index(string(text), ",")
	if x < 0 {
		return fmt.Errorf("Malformed content data: %s", textutil.Truncate(string(text), 100, " [...]"))
	}
	dec, err := base64.StdEncoding.DecodeString(string(text[x+1:]))
	if err != nil {
		return err
	}
	*d = Content{
		Type: mime.Type(text[:x]),
		Data: string(dec),
	}
	return nil
}

func (d *Content) Value() (driver.Value, error) {
	data, err := d.MarshalColumn()
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (d *Content) Scan(src interface{}) error {
	var err error
	switch c := src.(type) {
	case []byte:
		err = d.UnmarshalColumn(c)
	case string:
		err = d.UnmarshalColumn([]byte(c))
	default:
		err = fmt.Errorf("Unsupported type: %T", src)
	}
	return err
}
