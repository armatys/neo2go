package neo2go

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var paramsRe, listParamRe, batchParamRe *regexp.Regexp

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

func NewUrlTemplate(url string) *UrlTemplate {
	u := new(UrlTemplate)
	u.sections = make([]interface{}, 0)
	u.template = url
	u.parse()
	return u
}

func (u *UrlTemplate) parse() {
	indices := u.paramIndices()
	if len(indices) == 0 {
		u.sections = append(u.sections, [2]int{0, len(u.template)})
		return
	}

	prevIndex := 0

	for _, indexPair := range indices {
		if batchParamRe.MatchString(u.template[indexPair[0]:indexPair[1]]) {
			u.sections = append(u.sections, [2]int{indexPair[0], indexPair[1]})
			prevIndex = indexPair[1]
			continue
		}
		if indexPair[0] > prevIndex {
			u.sections = append(u.sections, [2]int{prevIndex, indexPair[0]})
		}
		prevIndex = indexPair[1]

		paramString := u.template[indexPair[0]:indexPair[1]]
		s := strings.Trim(paramString, "{}")

		matches := listParamRe.FindStringSubmatch(s)
		if matches != nil || len(matches) == 3 {
			u.parseList(matches)
		} else {
			u.parseNamedParams(s)
		}
	}
	if prevIndex < len(u.template) {
		u.sections = append(u.sections, [2]int{prevIndex, len(u.template)})
	}
}

func (u *UrlTemplate) paramIndices() [][]int {
	return paramsRe.FindAllStringIndex(u.template, -1)
}

func (u *UrlTemplate) parseList(matches []string) {
	var param urlParameter
	param.Delimiter = matches[1]
	param.Name = matches[2]
	u.sections = append(u.sections, param)
}

func (u *UrlTemplate) parseNamedParams(s string) {
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
}

func (u *UrlTemplate) renderUrlParameterIntoBuffer(urlparam *urlParameter, buf *bytes.Buffer, paramValue interface{}, questionMarkInserted *bool, wasLastSectionQueried *bool) (error, bool) {
	if paramValue == nil && !urlparam.Queried {
		return nil, true // true - should stop processing the template, even if there are more sections to process.
	}

	if s, ok := paramValue.(string); ok {
		if len(urlparam.Delimiter) > 0 {
			return fmt.Errorf("The type of the value for key '%v' is a `string`, but `[]string` was expected."), true
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
			buf.WriteString(url.QueryEscape(s))
		}
	} else if arr, ok := paramValue.([]string); ok {
		if len(urlparam.Delimiter) == 0 {
			return fmt.Errorf("The type of the value for key '%v' is a `[]string`, but `string` was expected."), true
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
		return fmt.Errorf("The type of the value for key '%v' is not supported (use `string` or `[]string`).", urlparam.Name), true
	}

	return nil, false
}

func (u *UrlTemplate) Render(params map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	questionMarkInserted := false
	wasLastSectionQueried := false

	for _, section := range u.sections {
		if indices, ok := section.([2]int); ok {
			buf.WriteString(u.template[indices[0]:indices[1]])
		} else if urlparam, ok := section.(urlParameter); ok {
			err, shouldStop := u.renderUrlParameterIntoBuffer(&urlparam, &buf, params[urlparam.Name], &questionMarkInserted, &wasLastSectionQueried)
			if err != nil {
				return "", err
			}
			if shouldStop {
				break
			}
		}
	}

	return buf.String(), nil
}

func (u *UrlTemplate) String() string {
	return u.template
}

func (u *UrlTemplate) UnmarshalJSON(data []byte) error {
	u.template = ""
	u.sections = u.sections[:0]
	s := string(data)
	u.template = s[1 : len(s)-1]
	u.parse()
	return nil
}

func init() {
	paramsRe = regexp.MustCompile(`{[^}]+}`)
	listParamRe = regexp.MustCompile(`^-list\|([^\|]+)\|([^\s]+)$`)
	batchParamRe = regexp.MustCompile(`{[0-9]+}`)
}
