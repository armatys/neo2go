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
