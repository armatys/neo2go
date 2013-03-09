package neo2go

import (
	"fmt"
	"testing"
)

var _ = fmt.Println

func TestConnecting(t *testing.T) {
	service, err := NewGraphDatabaseService("http://localhost:7474/db/data")
	if err != nil {
		t.Fatalf("Error while connecting: %v", err.Error())
	}
	if len(service.builder.root.Neo4jVersion) == 0 {
		t.Fatalf("Expected to receive neo4j version identifier.")
	}
}

func TestConnectingConnectionRefused(t *testing.T) {
	_, err := NewGraphDatabaseService("http://localhost:38479/db/data")
	if err == nil {
		t.Fatalf("Connection succeeded, but should not.")
	}
}

func TestCreateGetDeleteNode(t *testing.T) {
	service, err := NewGraphDatabaseService("http://localhost:7474/db/data")
	if err != nil {
		t.Fatalf("Error while connecting: %v", err.Error())
	}

	node, resp := service.CreateNode()
	if resp.StatusCode != 201 {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	node2, resp := service.GetNode(node.Self.String())
	if resp.StatusCode != 200 {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	resp = service.DeleteNode(node2)
	if resp.StatusCode != 204 {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	_, resp = service.GetNode(node.Self.String())
	if resp.StatusCode != 404 {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}
}
