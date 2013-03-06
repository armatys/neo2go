package neo2go

import (
	"fmt"
	"regexp"
	"strings"
)

var paramsRe, listParamRe *regexp.Regexp

type urlParameter struct {
	// Only used for list parameters.
	Delimiter string
	// Name of the parameter.
	Name string
	// If true, it means the parameter is to be used as URL query, after the '?' sign, as name=value pair.
	Queried bool
}

type UrlTemplate struct {
	parameters []urlParameter
	template   string
}

func (u *UrlTemplate) parse() error {
	indices := u.paramIndices()
	for _, indexPair := range indices {
		paramString := u.template[indexPair[0]:indexPair[1]]
		s := strings.Trim(paramString, "{}")

		var err error = nil
		if strings.HasPrefix(s, "-list|") {
			err = u.parseList(s)
		} else {
			err = u.parseNamedParams(s)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (u *UrlTemplate) paramIndices() [][]int {
	return paramsRe.FindAllStringIndex(u.template, -1)
}

func (u *UrlTemplate) parseList(s string) error {
	matches := listParamRe.FindStringSubmatch(s)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Could not find match for the url parameter: %s", s)
	}

	var param urlParameter
	param.Delimiter = matches[1]
	param.Name = matches[2]
	u.parameters = append(u.parameters, param)

	return nil
}

func (u *UrlTemplate) parseNamedParams(s string) error {
	queried := false

	if s[0] == '?' {
		queried = true
		s = s[1:]
	}

	paramNames := strings.Split(s, ",")
	for _, name := range paramNames {
		var param urlParameter
		param.Queried = queried
		param.Name = name
		u.parameters = append(u.parameters, param)
	}

	return nil
}

func (u *UrlTemplate) Render(params ...string) string {
	return ""
}

func (u *UrlTemplate) String() string {
	return u.template
}

func (u *UrlTemplate) UnmarshalJSON(data []byte) error {
	u.template = string(data)
	return u.parse()
}

func init() {
	paramsRe = regexp.MustCompile("{[^}]+}")
	listParamRe = regexp.MustCompile(`^-list\|([^\|]+)\|([^\s]+)$`)
}
