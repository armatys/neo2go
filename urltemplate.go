package neo2go

import (
	"bytes"
	"fmt"
	"net/url"
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
	// Contains the string as passed by the server.
	template string
	// Contains either a 2-element int array of indices that correspond
	// to sections of the template which do not need to be rendered;
	// or a urlParameter.
	sections []interface{}
}

func NewUrlTemplate(url string) (*UrlTemplate, error) {
	u := new(UrlTemplate)
	u.template = url
	return u, u.parse()
}

func NewPlainUrlTemplate(url string) *UrlTemplate {
	u := new(UrlTemplate)
	u.template = url
	u.sections = append(u.sections, [2]int{0, len(url)})
	return u
}

func (u *UrlTemplate) parse() error {
	indices := u.paramIndices()
	prevIndex := 0

	for _, indexPair := range indices {
		if indexPair[0] > prevIndex {
			u.sections = append(u.sections, [2]int{prevIndex, indexPair[0]})
		}
		prevIndex = indexPair[1]

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
	u.sections = append(u.sections, param)

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
		u.sections = append(u.sections, param)
	}

	return nil
}

func (u *UrlTemplate) renderUrlParameterIntoBuffer(urlparam *urlParameter, buf *bytes.Buffer, paramValue interface{}, questionMarkInserted *bool, wasLastSectionQueried *bool) error {
	if paramValue == nil && !urlparam.Queried {
		return fmt.Errorf("The value for key '%v' is required.", urlparam.Name)
	}

	if s, ok := paramValue.(string); ok {
		if len(urlparam.Delimiter) > 0 {
			return fmt.Errorf("The type of the value for key '%v' is a `string`, but `[]string` was expected.")
		}

		if urlparam.Queried {
			if !*questionMarkInserted {
				buf.WriteString("?")
				*questionMarkInserted = true
			}
			if *wasLastSectionQueried {
				buf.WriteString("&")
			}

			buf.WriteString(url.QueryEscape(urlparam.Name))
			buf.WriteString("=")
			buf.WriteString(url.QueryEscape(s))
			*wasLastSectionQueried = true
		} else {
			buf.WriteString(s)
		}
	} else if arr, ok := paramValue.([]string); ok {
		if len(urlparam.Delimiter) == 0 {
			return fmt.Errorf("The type of the value for key '%v' is a `[]string`, but `string` was expected.")
		}

		maxIndexForPlacingDelimiter := len(arr) - 2
		escapedDelimiter := url.QueryEscape(urlparam.Delimiter)
		for i, s := range arr {
			buf.WriteString(url.QueryEscape(s))
			if i <= maxIndexForPlacingDelimiter {
				buf.WriteString(escapedDelimiter)
			}
		}
	} else if paramValue != nil {
		return fmt.Errorf("The type of the value for key '%v' is not supported (use `string` or `[]string`).", urlparam.Name)
	}

	return nil
}

func (u *UrlTemplate) Render(params map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	questionMarkInserted := false
	wasLastSectionQueried := false

	for _, section := range u.sections {
		if indices, ok := section.([2]int); ok {
			buf.WriteString(u.template[indices[0]:indices[1]])
		} else if urlparam, ok := section.(urlParameter); ok {
			err := u.renderUrlParameterIntoBuffer(&urlparam, &buf, params[urlparam.Name], &questionMarkInserted, &wasLastSectionQueried)
			if err != nil {
				return "", err
			}
		}
	}

	return buf.String(), nil
}

func (u *UrlTemplate) String() string {
	return u.template
}

func (u *UrlTemplate) UnmarshalJSON(data []byte) error {
	s := string(data)
	u.template = s[1 : len(s)-1]
	return u.parse()
}

func init() {
	paramsRe = regexp.MustCompile("{[^}]+}")
	listParamRe = regexp.MustCompile(`^-list\|([^\|]+)\|([^\s]+)$`)
}
