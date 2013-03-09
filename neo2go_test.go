package neo2go

import (
	"fmt"
	"testing"
)

var _ = fmt.Println

func TestConnecting(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect("http://localhost:7474/db/data")
	if resp.StatusCode != 200 {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}
	if len(service.builder.root.Neo4jVersion) == 0 {
		t.Fatalf("Expected to receive neo4j version identifier.")
	}
}

func TestConnectingConnectionRefused(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect("http://localhost:38479/db/data")
	if resp.NeoError == nil {
		t.Fatalf("Connection succeeded, but should not.")
	}
}

func TestCreateNodeWithoutConnecting(t *testing.T) {
	service := NewGraphDatabaseService()

	_, resp := service.GetNode("http://localhost:38479/db/data/node/1")
	if (resp.StatusCode == 600) != !resp.Ok() {
		t.Fatalf("Request should not have succeeded, since Connect method was not called")
	}
}

func TestCreateGetDeleteNode(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect("http://localhost:7474/db/data")
	if (resp.StatusCode == 200) != resp.Ok() {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	node, resp := service.CreateNode()
	if (resp.StatusCode == 201) != resp.Ok() {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	node2, resp := service.GetNode(node.Self.String())
	if (resp.StatusCode == 200) != resp.Ok() {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	resp = service.DeleteNode(node2)
	if (resp.StatusCode == 204) != resp.Ok() {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	_, resp = service.GetNode(node.Self.String())
	if (resp.StatusCode == 404) != !resp.Ok() {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}
}
