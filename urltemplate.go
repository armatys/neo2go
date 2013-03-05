package neo2go

import (
	"regexp"
)

var paramsRegExp *regexp.Regexp

func init() {
	paramsRegExp = regexp.MustCompile("{[^}]+}")
}

type UrlTemplate struct {
	template string
}

func (u *UrlTemplate) paramIndices() [][]int {
	return paramsRegExp.FindAllStringIndex(u.template, -1)
}

func (u *UrlTemplate) parse() {
	_ = u.paramIndices()
}

func (u *UrlTemplate) Render(params ...string) string {
	return ""
}

func (u *UrlTemplate) String() string {
	return u.template
}

func (u *UrlTemplate) UnmarshalJSON(data []byte) error {
	u.template = string(data)
	return nil
}
