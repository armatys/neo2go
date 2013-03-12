package neo2go

import (
	"reflect"
	"testing"
)

func TestParsingIndices(t *testing.T) {
	excepted := [][]int{
		[]int{52, 64},
		[]int{64, 85},
	}

	tpl := new(UrlTemplate)
	tpl.template = "http://localhost:7474/db/data/node/1/paged/traverse/{returnType}{?pageSize,leaseTime}"
	indices := tpl.paramIndices()

	if len(indices) != len(excepted) {
		t.Fatalf("The number of parsed indices should be 2, but is %d", len(indices))
	}

	for i := 0; i < len(indices); i++ {
		if indices[i][0] != excepted[i][0] || indices[i][1] != excepted[i][1] {
			t.Fatalf("The first index pair should be %v, but is %v", excepted[i], indices[i])
		}
	}
}

func TestParsingList(t *testing.T) {
	tpl := new(UrlTemplate)
	tpl.template = "http://localhost:7474/db/data/node/9/relationships/all/{-list|&|types}"
	err := tpl.parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(tpl.sections) != 2 {
		t.Fatalf("Expected 2 sections, but got %d", len(tpl.sections))
	}

	if section0, ok := tpl.sections[0].([2]int); ok {
		sExpected := [2]int{0, 55}
		if section0[0] != sExpected[0] || section0[1] != sExpected[1] {
			t.Fatalf("Expected the first section to be '%v', but got %v", sExpected, section0)
		}
	} else {
		t.Fatalf("Expected the first section to be of type `[]int`.")
	}

	if param, ok := tpl.sections[1].(urlParameter); ok {
		if param.Delimiter != "&" {
			t.Fatalf("Expected '&' delimiter, but got %v", param.Delimiter)
		}

		if param.Queried {
			t.Fatalf("Expected the 'Queried' property to be false.")
		}

		name := "types"
		if param.Name != name {
			t.Fatalf("Expected the parameter name to be %v, but got %v.", name, param.Name)
		}
	} else {
		t.Fatalf("Expected second section to be of type `urlParameter`.")
	}
}

func TestParsingNamedParams(t *testing.T) {
	tpl := new(UrlTemplate)
	tpl.template = "http://localhost:7474/db/data/node/1/paged/traverse/{returnType}{?pageSize,leaseTime}"
	err := tpl.parse()
	if err != nil {
		t.Fatal(err)
	}

	expected := []interface{}{
		[2]int{0, 52},
		urlParameter{"", "returnType", false},
		urlParameter{"", "pageSize", true},
		urlParameter{"", "leaseTime", true},
	}

	if len(tpl.sections) != len(expected) {
		t.Fatalf("Expected %d sections, but got %d", len(expected), len(tpl.sections))
	}

	for i := 0; i < len(expected); i++ {
		exemplary := expected[i]
		parsed := tpl.sections[i]

		exemplaryType := reflect.TypeOf(exemplary)
		if exemplaryType != reflect.TypeOf(parsed) {
			t.Fatalf("Types of parsed sections do not match.")
		}

		if exemplaryType == reflect.TypeOf([]int{}) {
			exemplaryArray, _ := exemplary.([2]int)
			parsedArray, _ := parsed.([2]int)
			if exemplaryArray[0] != parsedArray[0] || exemplaryArray[1] != parsedArray[1] {
				t.Fatalf("Expected '%v' section, but got %v", exemplaryArray, parsedArray)
			}
		} else {
			exemplaryParam, _ := exemplary.(urlParameter)
			parsedParam, _ := parsed.(urlParameter)

			if exemplaryParam.Delimiter != parsedParam.Delimiter {
				t.Fatalf("Expected '%v' delimiter, but got %v", exemplaryParam.Delimiter, parsedParam.Delimiter)
			}

			if exemplaryParam.Queried != parsedParam.Queried {
				t.Fatalf("Expected the 'Queried' property to be %v, but got %v.", exemplaryParam.Queried, parsedParam.Queried)
			}

			if exemplaryParam.Name != parsedParam.Name {
				t.Fatalf("Expected the parameter name to be %v, but got %v.", exemplaryParam.Name, parsedParam.Name)
			}
		}
	}
}

func TestRenderingWithNoParameters(t *testing.T) {
	tpl := new(UrlTemplate)
	tpl.UnmarshalJSON([]byte(`"http://localhost:7474/db/data/node/1/paged/traverse/{returnType}{?pageSize,leaseTime}"`))
	s, err := tpl.Render(nil)
	expected := "http://localhost:7474/db/data/node/1/paged/traverse/"
	if err != nil {
		t.Fatalf("Unexpected error: '%v'", err)
	}

	if s != expected {
		t.Fatalf("Expected to render '%v', but got '%v'", expected, s)
	}
}

func TestRenderingWithQueryParam(t *testing.T) {
	tpl := new(UrlTemplate)
	tpl.UnmarshalJSON([]byte(`"http://localhost:7474/db/data/node/1/paged/traverse/{returnType}{?pageSize,leaseTime}"`))

	params := map[string]interface{}{
		"leaseTime": "100",
	}
	s, err := tpl.Render(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}

	expected := "http://localhost:7474/db/data/node/1/paged/traverse/?leaseTime=100"
	if s != expected {
		t.Fatalf("Expected to render '%v', but got '%v'", expected, s)
	}
}

func TestRenderingWithQueryParams(t *testing.T) {
	tpl := new(UrlTemplate)
	tpl.UnmarshalJSON([]byte(`"http://localhost:7474/db/data/node/1/paged/traverse/{returnType}{?pageSize,leaseTime}"`))

	params := map[string]interface{}{
		"leaseTime": "100",
		"pageSize":  "33",
	}
	s, err := tpl.Render(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}

	// The ordere of query params is determined by the template.
	expected := "http://localhost:7474/db/data/node/1/paged/traverse/?pageSize=33&leaseTime=100"
	if s != expected {
		t.Fatalf("Expected to render '%v', but got '%v'", expected, s)
	}
}

func TestRenderingSimpleParam(t *testing.T) {
	tpl := new(UrlTemplate)
	tpl.UnmarshalJSON([]byte(`"http://localhost:7474/db/data/node/1/paged/traverse/{returnType}{?pageSize,leaseTime}"`))

	params := map[string]interface{}{
		"returnType": "node",
	}
	s, err := tpl.Render(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}

	expected := "http://localhost:7474/db/data/node/1/paged/traverse/node"
	if s != expected {
		t.Fatalf("Expected to render '%v', but got '%v'", expected, s)
	}
}

func TestRenderingSimpleParamAndQueried(t *testing.T) {
	tpl := new(UrlTemplate)
	tpl.UnmarshalJSON([]byte(`"http://localhost:7474/db/data/node/1/paged/traverse/{returnType}{?pageSize,leaseTime}"`))

	params := map[string]interface{}{
		"leaseTime":  "100",
		"pageSize":   "10",
		"returnType": "node",
	}
	s, err := tpl.Render(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}

	// Note that the order of queried params is determined by the template.
	expected := "http://localhost:7474/db/data/node/1/paged/traverse/node?pageSize=10&leaseTime=100"
	if s != expected {
		t.Fatalf("Expected to render '%v', but got '%v'", expected, s)
	}
}

func TestRenderingArray(t *testing.T) {
	tpl := new(UrlTemplate)
	tpl.UnmarshalJSON([]byte(`"http://localhost:7474/db/data/node/9/relationships/all/{-list|&|types}"`))

	params := map[string]interface{}{
		"types": []string{"T1", "T2", "T3"},
	}
	s, err := tpl.Render(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}

	// The ampersand has to be url-encoded.
	expected := `http://localhost:7474/db/data/node/9/relationships/all/T1%26T2%26T3`
	if s != expected {
		t.Fatalf("Expected to render '%v', but got '%v'", expected, s)
	}
}
