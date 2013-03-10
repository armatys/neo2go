package neo2go

import (
	"fmt"
	"testing"
)

var _ = fmt.Println

const (
	databaseAddress                = "http://localhost:7474/db/data/"
	databaseAddressWithInvalidPort = "http://localhost:38479/db/data/"
)

func responseHasSucceededWithCode(resp *NeoResponse, expectedStatus int) bool {
	return ((resp.StatusCode == expectedStatus) == resp.Ok()) && resp.Ok()
}

func responseHasFailedWithCode(resp *NeoResponse, unexpectedStatus int) bool {
	return ((resp.StatusCode == unexpectedStatus) == !resp.Ok()) && !resp.Ok()
}

func TestConnecting(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}
	if len(service.builder.root.Neo4jVersion) == 0 {
		t.Fatalf("Expected to receive neo4j version identifier.")
	}
}

func TestConnectingConnectionRefused(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddressWithInvalidPort)
	if resp.NeoError == nil {
		t.Fatalf("Connection succeeded, but should not.")
	}
}

func TestCreateNodeWithoutConnecting(t *testing.T) {
	service := NewGraphDatabaseService()

	_, resp := service.GetNode(databaseAddress + "node/1")
	if !responseHasFailedWithCode(resp, 600) {
		t.Fatalf("Request should not have succeeded, since Connect method was not called")
	}
}

func TestCreateGetDeleteNode(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	node, resp := service.CreateNode()
	if !responseHasSucceededWithCode(resp, 201) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	node2, resp := service.GetNode(node.Self.String())
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	resp = service.DeleteNode(node2)
	if !responseHasSucceededWithCode(resp, 204) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	_, resp = service.GetNode(node.Self.String())
	if !responseHasFailedWithCode(resp, 404) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}
}

func TestCreateNodeWithProperties(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	node, resp := service.CreateNodeWithProperties(map[string]interface{}{"name": "jan", "age": 99})
	if !responseHasSucceededWithCode(resp, 201) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	node2, resp := service.GetNode(node.Self.String())
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	resp = service.DeleteNode(node2)
	if !responseHasSucceededWithCode(resp, 204) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	_, resp = service.GetNode(node.Self.String())
	if !responseHasFailedWithCode(resp, 404) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}
}

func TestSimpleCypherQueryFail(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	_, resp = service.Cypher("START x = node(186) RETURN x.dummy", nil)
	if !responseHasFailedWithCode(resp, 400) {
		t.Fatalf("Expected the cypher query to fail, but the status code was %d", resp.StatusCode)
	}
}

func TestDeleteNonExistendNode(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	node := new(NeoNode)
	node.Self = NewPlainUrlTemplate(databaseAddress + "node/909090")

	resp = service.DeleteNode(node)
	if !responseHasFailedWithCode(resp, 404) {
		t.Fatalf("Unexpected response: %v", resp.StatusCode)
	}
}

func TestSimpleCypherQuery(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	node, resp := service.CreateNodeWithProperties(map[string]interface{}{"name": "jan", "age": 99})
	if !responseHasSucceededWithCode(resp, 201) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	cypherResult, resp := service.Cypher("START x = node({nid}) RETURN x.name", map[string]interface{}{"nid": node.Id()})
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Unexpected result %v", resp.NeoError)
	}

	numOfColumnts := 1
	if len(cypherResult.Columns) != numOfColumnts {
		t.Fatalf("Expected %d columnts, but got %d", numOfColumnts, len(cypherResult.Columns))
	}

	colName := "x.name"
	if cypherResult.Columns[0] != colName {
		t.Fatalf("Expected column name to have a name %v, but got %v", colName, cypherResult.Columns[0])
	}

	if value, ok := cypherResult.Data[0][0].(string); ok {
		expectedValue := "jan"
		if value != expectedValue {
			t.Fatalf("Expected the value returned by Cypher to be %v, but got %v", expectedValue, value)
		}
	} else {
		t.Fatalf("Expected cypher data to be a string.")
	}

	resp = service.DeleteNode(node)
	if !responseHasSucceededWithCode(resp, 204) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}
}
