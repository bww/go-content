package content

import (
	"github.com/bww/go-mime/v1"
)

type Set []Content

func (s Set) Choose(accept []mime.Type) Content {
	for _, e := range accept {
		for _, c := range s {
			if e.Equals(c.Type) {
				return c
			}
		}
	}
	if len(s) > 0 {
		return s[0]
	} else {
		return Content{} // 🤷‍♂️
	}
}
