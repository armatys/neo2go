package neo2go

import (
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

	if len(tpl.parameters) != 1 {
		t.Fatalf("Expected 1 parameter, but got %d", len(tpl.parameters))
	}

	param := tpl.parameters[0]

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
}

func TestParsingNamedParams(t *testing.T) {
	tpl := new(UrlTemplate)
	tpl.template = "http://localhost:7474/db/data/node/1/paged/traverse/{returnType}{?pageSize,leaseTime}"
	err := tpl.parse()
	if err != nil {
		t.Fatal(err)
	}

	expected := []urlParameter{
		urlParameter{"", "returnType", false},
		urlParameter{"", "pageSize", true},
		urlParameter{"", "leaseTime", true},
	}

	if len(tpl.parameters) != len(expected) {
		t.Fatalf("Expected 1 parameter, but got %d", len(tpl.parameters))
	}

	for i := 0; i < len(expected); i++ {
		exemplary := expected[i]
		parsed := tpl.parameters[i]

		if exemplary.Delimiter != parsed.Delimiter {
			t.Fatalf("Expected '%v' delimiter, but got %v", exemplary.Delimiter, parsed.Delimiter)
		}

		if exemplary.Queried != parsed.Queried {
			t.Fatalf("Expected the 'Queried' property to be %v, but got %v.", exemplary.Queried, parsed.Queried)
		}

		if exemplary.Name != parsed.Name {
			t.Fatalf("Expected the parameter name to be %v, but got %v.", exemplary.Name, parsed.Name)
		}
	}
}
