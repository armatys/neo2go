package neo2go

import (
	"testing"
)

func TestConnecting(t *testing.T) {
	service := NewGraphDatabaseService("http://localhost:7474/db/data")
	neoResp, err := service.Connect()
	if err != nil {
		t.Fatal(err)
	}
	if neoResp.StatusCode != 200 {
		t.Fatal("Server returned status code %d, but 200 was expected.", neoResp.StatusCode)
	}
	if len(service.Neo4jVersion) == 0 {
		t.Fatal("Expected to receive neo4j version identifier.")
	}
}

func TestConnectingConnectionRefused(t *testing.T) {
	service := NewGraphDatabaseService("http://localhost:38479/db/data")
	_, err := service.Connect()
	if err == nil {
		t.Fatal("Connection succeeded, but should not.")
	}
}
