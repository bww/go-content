package content

import (
	"database/sql/driver"
	"fmt"

	"github.com/bww/go-mime/v1"
	"github.com/lib/pq"
)

type Set []*Content

// Choose finds the first content in the set that matches the first
// acceptable mimetype. The acceptable types are provided in order of
// preference as are the content elements.
//
// If no acceptable types are provided, the effect is that the first
// element in the set is chosen, or if the set is empty, nil.
func (s Set) Choose(accept ...mime.Type) *Content {
	for _, e := range accept {
		for _, c := range s {
			if e.Equals(c.Type) {
				return c
			}
		}
	}
	if len(s) > 0 {
		return s[0] // no match, choose the first alternate
	} else {
		return nil
	}
}

func (s Set) Value() (driver.Value, error) {
	var c = make([]string, len(s))
	for i, e := range s {
		v, err := e.MarshalColumn()
		if err != nil {
			return nil, fmt.Errorf("Could not marshal element #%d: %w", i, err)
		}
		c[i] = string(v)
	}
	return pq.Array(c).Value()
}

func (s *Set) Scan(src interface{}) error {
	a := pq.StringArray{}
	err := a.Scan(src)
	if err != nil {
		return err
	}
	r := make(Set, len(a))
	for i, e := range a {
		v := &Content{}
		err := v.UnmarshalColumn([]byte(e))
		if err != nil {
			return fmt.Errorf("Could not unmarshal element #%d: %w", i, err)
		}
		r[i] = v
	}
	*s = r
	return nil
}
