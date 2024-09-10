package html

import (
	"github.com/microcosm-cc/bluemonday"
)

var defaultPolicy = bluemonday.UGCPolicy()

func Sanitize(t string) string {
	return defaultPolicy.Sanitize(t)
}
