package neo2go

import (
	"fmt"
	"testing"
)

var _ = fmt.Println

func TestConnecting(t *testing.T) {
	service, err := NewGraphDatabaseService("http://localhost:7474/db/data")
	if err != nil {
		t.Fatal(err)
	}
	neoResp := service.Connect()
	if neoResp.StatusCode != 200 {
		t.Fatalf("Server returned status code %d, but 200 was expected.", neoResp.StatusCode)
	}
	if len(service.builder.root.Neo4jVersion) == 0 {
		t.Fatalf("Expected to receive neo4j version identifier.")
	}
}

func TestConnectingConnectionRefused(t *testing.T) {
	service, err := NewGraphDatabaseService("http://localhost:38479/db/data")
	if err != nil {
		t.Fatal(err)
	}
	neoResp := service.Connect()
	if neoResp.NeoError == nil {
		t.Fatalf("Connection succeeded, but should not.")
	}
}

func TestConnectingGetNodeWithNoConnection(t *testing.T) {
	service, err := NewGraphDatabaseService("http://localhost:7474/db/data")
	if err != nil {
		t.Fatal(err)
	}

	_, resp := service.GetNode("http://localhost:38479/db/data/node/1")
	if resp.StatusCode != 600 {
		t.Fatal(resp.NeoError)
	}
}
