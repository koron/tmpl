package tmpl

import "strings"

func init() {
	AddFunc("hasPrefix", strings.HasPrefix)
	AddFunc("hasSuffix", strings.HasSuffix)
	AddFunc("trimPrefix", strings.TrimPrefix)
	AddFunc("trimSuffix", strings.TrimSuffix)
}
