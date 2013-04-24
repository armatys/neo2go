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

func checkResponseSucceeded(t *testing.T, resp *NeoResponse, expectedStatus int) {
	if !responseHasSucceededWithCode(resp, expectedStatus) {
		t.Fatalf("Expected %d response, but got error %d: %v", expectedStatus, resp.StatusCode, resp.NeoError)
	}
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
	node.Self = NewUrlTemplate(databaseAddress + "node/909090")

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

func TestSimpleRelationships(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	source, resp := service.CreateNode()
	if !resp.Ok() {
		t.Fatalf(resp.NeoError.Error())
	}
	target, resp := service.CreateNode()
	if !resp.Ok() {
		t.Fatalf(resp.NeoError.Error())
	}

	rel1, resp := service.CreateRelationshipWithType(source, target, "likes")
	if !resp.Ok() {
		t.Fatalf("Error creating relationship: %v", resp.NeoError.Error())
	}

	properties := map[string]float64{"weight": 30}
	relType := "has"
	rel, resp := service.CreateRelationshipWithPropertiesAndType(source, target, properties, relType)
	if !resp.Ok() {
		t.Fatalf("Error creating relationship (with properties): %v", resp.NeoError.Error())
	}

	var weight float64
	resp = service.GetPropertyForRelationship(rel, "weight", &weight)
	if !resp.Ok() {
		t.Fatalf("Could not get relationship property: %v", resp.NeoError.Error())
	}
	expected := 30.0
	if weight != expected {
		t.Fatalf("Expected relationship property value %d, but got: %d", expected, weight)
	}

	expected = 45.0
	service.SetPropertyForRelationship(rel, "weight", expected)
	if !resp.Ok() {
		t.Fatalf("Could not set relationship property: %v", resp.NeoError.Error())
	}

	var props map[string]float64
	resp = service.GetPropertiesForRelationship(rel, &props)
	if !resp.Ok() {
		t.Fatalf("Could not get relationship property: %v", resp.NeoError.Error())
	}

	val := props["weight"]
	if val != expected {
		t.Fatalf("Expected relationship property value %d, but got: %d", expected, val)
	}

	rels, resp := service.GetRelationshipsWithTypesForNode(source, NeoTraversalOut, []string{relType})
	if !resp.Ok() {
		t.Fatalf("Could not get node relationships: %v", resp.NeoError.Error())
	}

	if len(*rels) != 1 {
		t.Fatalf("Expected to get 1-element array of relationships but got %d elements.", len(*rels))
	}

	relTypes, resp := service.GetRelationshipTypes()
	if !resp.Ok() {
		t.Fatalf("Could not get relationship types: %v", resp.NeoError.Error())
	}

	if len(*relTypes) < 2 {
		t.Fatalf("Expected at least 2 relationship types, but got: %d", len(*relTypes))
	}

	if len(*relTypes) != 2 {
		t.Logf("Expected at 2 relationship types, but got: %d", len(*relTypes))
	}

	resp = service.DeleteRelationship(rel1)
	if !resp.Ok() {
		t.Fatalf("Error deleting relationship: %v", resp.NeoError.Error())
	}
	resp = service.DeleteRelationship(rel)
	if !resp.Ok() {
		t.Fatalf("Error deleting relationship: %v", resp.NeoError.Error())
	}
	resp = service.DeleteNode(source)
	if !resp.Ok() {
		t.Fatalf("Unexpected response (%v): %v", resp.StatusCode, resp.NeoError)
	}
	resp = service.DeleteNode(target)
	if !resp.Ok() {
		t.Fatalf("Unexpected response (%v): %v", resp.StatusCode, resp.NeoError)
	}
}

func TestSimpleRelationships2(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	source, resp := service.CreateNode()
	if !resp.Ok() {
		t.Fatalf(resp.NeoError.Error())
	}
	target, resp := service.CreateNode()
	if !resp.Ok() {
		t.Fatalf(resp.NeoError.Error())
	}

	rel, resp := service.CreateRelationshipWithPropertiesAndType(source, target, &map[string]string{"one": "two"}, "likes")
	if !resp.Ok() {
		t.Fatalf("Error creating relationship: %v", resp.NeoError.Error())
	}

	rel, resp = service.GetRelationship(rel.Self.String())
	if !resp.Ok() {
		t.Fatalf("Expected 200 code but got: %d", resp.StatusCode)
	}

	if v, ok := rel.Data.(map[string]interface{}); ok {
		expected := "two"
		if v["one"] != expected {
			t.Fatalf("Expected the property value to be `%v`, but got: `%v`", expected, v["one"])
		}
	} else {
		t.Fatalf("Could not convert properties to map[string]string")
	}

	resp = service.ReplacePropertiesForRelationship(rel, &map[string]string{"three": "four"})
	if !resp.Ok() {
		t.Fatalf("Error updating relationship properties: %v", resp.NeoError.Error())
	}

	var props map[string]string
	resp = service.GetPropertiesForRelationship(rel, &props)
	if !resp.Ok() {
		t.Fatalf("Error getting relationship properties: %v", resp.NeoError.Error())
	}
	expected := "four"
	if props["three"] != expected {
		t.Fatalf("Expected the property value to be `%v`, but got: `%v`", expected, props["three"])
	}
	expected = ""
	if props["one"] != expected {
		t.Fatalf("Expected the property value to be `%v`, but got: `%v`", expected, props["one"])
	}

	batch := service.Batch()
	batch.DeleteRelationship(rel)
	batch.DeleteNode(source)
	batch.DeleteNode(target)
	resp = batch.Commit()
	if !resp.Ok() {
		t.Fatalf("Could not execute batch (deletion): %v", resp.NeoError.Error())
	}

	_, resp = service.GetNode(source.Self.String())
	if !responseHasFailedWithCode(resp, 404) {
		t.Fatalf("Expected 404 code but got: %d", resp.StatusCode)
	}
	_, resp = service.GetNode(target.Self.String())
	if !responseHasFailedWithCode(resp, 404) {
		t.Fatalf("Expected 404 code but got: %d", resp.StatusCode)
	}
	_, resp = service.GetRelationship(rel.Self.String())
	if !responseHasFailedWithCode(resp, 404) {
		t.Fatalf("Expected 404 code but got: %d", resp.StatusCode)
	}
}

func TestCreateDeleteNodeIndex(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	indexName := "test-n"
	index, resp := service.CreateNodeIndex(indexName)
	if !responseHasSucceededWithCode(resp, 201) {
		t.Fatalf("Expected 201 response, but got error: %v", resp.NeoError)
	}
	indexes, resp := service.GetNodeIndexes()
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Expected 200 response, but got error: %v", resp.NeoError)
	}
	if (*indexes)[indexName] == nil {
		t.Fatalf("Expected to get a '%v' index but got nil.", indexName)
	}

	resp = service.DeleteIndex(index)
	if !responseHasSucceededWithCode(resp, 204) {
		t.Fatalf("Expected 204 response, but got error: %v", resp.NeoError)
	}
}

func TestFindExactNodeNoMatches(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	checkResponseSucceeded(t, resp, 200)

	indexName := "test-n"
	index, resp := service.CreateNodeIndex(indexName)
	checkResponseSucceeded(t, resp, 201)

	nodes, resp := service.FindNodeByExactMatch(index, "name", "text-value")
	checkResponseSucceeded(t, resp, 200)

	if len(*nodes) > 0 {
		t.Fatalf("Expected to get 0 nodes but got %d", len(*nodes))
	}

	resp = service.DeleteIndex(index)
	checkResponseSucceeded(t, resp, 204)
}

func TestFindExactNodeMatches(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	checkResponseSucceeded(t, resp, 200)

	indexName := "test-n"
	index, resp := service.CreateNodeIndex(indexName)
	checkResponseSucceeded(t, resp, 201)

	node, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	indexedNode, resp := service.AddNodeToIndex(index, node, "name", "text-value")
	checkResponseSucceeded(t, resp, 201)
	if node.Self.String() != indexedNode.Self.String() {
		t.Fatalf("Expected to get the same node after indexing")
	}

	nodes, resp := service.FindNodeByExactMatch(index, "name", "text-value")
	checkResponseSucceeded(t, resp, 200)

	if len(*nodes) != 1 {
		t.Fatalf("Expected to get 1 nodes but got %d", len(*nodes))
	}

	resp = service.DeleteIndex(index)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(node)
	checkResponseSucceeded(t, resp, 204)
}

func TestCreateUniqueNode(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	checkResponseSucceeded(t, resp, 200)

	indexName := "test-n"
	index, resp := service.CreateNodeIndex(indexName)
	checkResponseSucceeded(t, resp, 201)

	createdNode, resp := service.GetOrCreateUniqueNode(index, "name", "text-value")
	if !resp.Created() || resp.StatusCode != 201 {
		t.Fatalf("Unexpected response %d: %v", resp.StatusCode, resp.NeoError)
	}

	fetchedNode, resp := service.GetOrCreateUniqueNode(index, "name", "text-value")
	if !resp.Ok() || resp.StatusCode != 200 {
		t.Fatalf("Unexpected response %d: %v", resp.StatusCode, resp.NeoError)
	}

	if createdNode.Self.String() != fetchedNode.Self.String() {
		t.Fatalf("Expected to get the same nodes, but got (created): %v and (fetched): %v", createdNode.Self.String(), fetchedNode.Self.String())
	}

	nodes, resp := service.FindNodeByExactMatch(index, "name", "text-value")
	checkResponseSucceeded(t, resp, 200)

	if len(*nodes) != 1 {
		t.Fatalf("Expected to get 1 nodes but got %d", len(*nodes))
	}

	if (*nodes)[0].Self.String() != createdNode.Self.String() {
		t.Fatalf("Expected to get the same node but got (create) %v and (by exact match) %v", createdNode.Self.String(), (*nodes)[0].Self.String())
	}

	resp = service.DeleteIndex(index)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(createdNode)
	checkResponseSucceeded(t, resp, 204)
}

func TestCreateUniqueNodeOrFail(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	checkResponseSucceeded(t, resp, 200)

	indexName := "test-n"
	index, resp := service.CreateNodeIndex(indexName)
	checkResponseSucceeded(t, resp, 201)

	createdNode, resp := service.CreateUniqueNodeOrFail(index, "name", "text-value")
	if !resp.Created() || resp.StatusCode != 201 {
		t.Fatalf("Unexpected response %d: %v", resp.StatusCode, resp.NeoError)
	}

	_, resp = service.CreateUniqueNodeOrFail(index, "name", "text-value")
	if resp.Ok() || resp.StatusCode == 201 {
		t.Fatalf("Unexpected response %d: %v", resp.StatusCode, resp.NeoError)
	}

	nodes, resp := service.FindNodeByExactMatch(index, "name", "text-value")
	checkResponseSucceeded(t, resp, 200)

	if len(*nodes) != 1 {
		t.Fatalf("Expected to get 1 nodes but got %d", len(*nodes))
	}

	if (*nodes)[0].Self.String() != createdNode.Self.String() {
		t.Fatalf("Expected to get the same node but got (create) %v and (by exact match) %v", createdNode.Self.String(), (*nodes)[0].Self.String())
	}

	resp = service.DeleteIndex(index)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(createdNode)
	checkResponseSucceeded(t, resp, 204)
}

func TestCreateDeleteRelationshipIndex(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	indexName := "test-r"
	index, resp := service.CreateRelationshipIndex(indexName)
	if !responseHasSucceededWithCode(resp, 201) {
		t.Fatalf("Expected 201 response, but got error: %v", resp.NeoError)
	}
	indexes, resp := service.GetRelationshipIndexes()
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Expected 200 response, but got error: %v", resp.NeoError)
	}
	if (*indexes)[indexName] == nil {
		t.Fatalf("Expected to get a '%v' index but got nil.", indexName)
	}

	resp = service.DeleteIndex(index)
	if !responseHasSucceededWithCode(resp, 204) {
		t.Fatalf("Expected 204 response, but got error: %v", resp.NeoError)
	}
}

func TestFindExactRelationshipNoMatches(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	indexName := "test-r"
	index, resp := service.CreateRelationshipIndex(indexName)
	if !responseHasSucceededWithCode(resp, 201) {
		t.Fatalf("Expected 201 response, but got error: %v", resp.NeoError)
	}

	rels, resp := service.FindRelationshipByExactMatch(index, "name", "text-value")
	checkResponseSucceeded(t, resp, 200)

	if len(*rels) > 0 {
		t.Fatalf("Expected to get 0 nodes but got %d", len(*rels))
	}

	resp = service.DeleteIndex(index)
	if !responseHasSucceededWithCode(resp, 204) {
		t.Fatalf("Expected 204 response, but got error: %v", resp.NeoError)
	}
}

func TestFindExactRelationshipMatches(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	checkResponseSucceeded(t, resp, 200)

	indexName := "test-n"
	index, resp := service.CreateRelationshipIndex(indexName)
	checkResponseSucceeded(t, resp, 201)

	source, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	target, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	rel, resp := service.CreateRelationshipWithType(source, target, "likes")
	checkResponseSucceeded(t, resp, 201)

	indexedRel, resp := service.AddRelationshipToIndex(index, rel, "name", "text-value")
	checkResponseSucceeded(t, resp, 201)
	if rel.Self.String() != indexedRel.Self.String() {
		t.Fatalf("Expected to get the same relationship after indexing (%v; %v).", rel.Self.String(), indexedRel.Self.String())
	}

	rels, resp := service.FindRelationshipByExactMatch(index, "name", "text-value")
	checkResponseSucceeded(t, resp, 200)

	if len(*rels) != 1 {
		t.Fatalf("Expected to get 1 relationships but got %d", len(*rels))
	}

	resp = service.DeleteIndex(index)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteRelationship(rel)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(source)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(target)
	checkResponseSucceeded(t, resp, 204)
}

func TestCreateUniqueRelationship(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	checkResponseSucceeded(t, resp, 200)

	indexName := "test-r"
	index, resp := service.CreateRelationshipIndex(indexName)
	checkResponseSucceeded(t, resp, 201)

	source, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	target, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	createdRel, resp := service.GetOrCreateUniqueRelationship(index, "name", "text-value", source, target, "rel-type")
	if !resp.Created() || resp.StatusCode != 201 {
		t.Fatalf("Unexpected response %d: %v", resp.StatusCode, resp.NeoError)
	}

	fetchedRel, resp := service.GetOrCreateUniqueRelationship(index, "name", "text-value", source, target, "rel-type")
	if !resp.Ok() || resp.StatusCode != 200 {
		t.Fatalf("Unexpected response %d: %v", resp.StatusCode, resp.NeoError)
	}

	if createdRel.Self.String() != fetchedRel.Self.String() {
		t.Fatalf("Expected to get the same nodes, but got (created): %v and (fetched): %v", createdRel.Self.String(), fetchedRel.Self.String())
	}

	rels, resp := service.FindRelationshipByExactMatch(index, "name", "text-value")
	checkResponseSucceeded(t, resp, 200)

	if len(*rels) != 1 {
		t.Fatalf("Expected to get 1 relationships but got %d", len(*rels))
	}

	if (*rels)[0].Self.String() != createdRel.Self.String() {
		t.Fatalf("Expected to get the same rel but got (create) %v and (by exact match) %v", createdRel.Self.String(), (*rels)[0].Self.String())
	}

	resp = service.DeleteIndex(index)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteRelationship(createdRel)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(source)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(target)
	checkResponseSucceeded(t, resp, 204)
}

func TestCreateUniqueRelationshipOrFail(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	checkResponseSucceeded(t, resp, 200)

	indexName := "test-r"
	index, resp := service.CreateRelationshipIndex(indexName)
	checkResponseSucceeded(t, resp, 201)

	source, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	target, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	createdRel, resp := service.CreateUniqueRelationshipOrFail(index, "name", "text-value", source, target, "likes")
	if !resp.Created() || resp.StatusCode != 201 {
		t.Fatalf("Unexpected response %d: %v", resp.StatusCode, resp.NeoError)
	}

	_, resp = service.CreateUniqueRelationshipOrFail(index, "name", "text-value", source, target, "likes")
	if resp.Created() || resp.Ok() || resp.StatusCode == 201 {
		t.Fatalf("Unexpected response %d: %v", resp.StatusCode, resp.NeoError)
	}

	rels, resp := service.FindRelationshipByExactMatch(index, "name", "text-value")
	checkResponseSucceeded(t, resp, 200)

	if len(*rels) != 1 {
		t.Fatalf("Expected to get 1 relationships but got %d", len(*rels))
	}

	if (*rels)[0].Self.String() != createdRel.Self.String() {
		t.Fatalf("Expected to get the same rel but got (create) %v and (by exact match) %v", createdRel.Self.String(), (*rels)[0].Self.String())
	}

	resp = service.DeleteIndex(index)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteRelationship(createdRel)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(source)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(target)
	checkResponseSucceeded(t, resp, 204)
}

func TestPathFinder(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	start, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	target, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	rel, resp := service.CreateRelationshipWithType(start, target, "likes")
	checkResponseSucceeded(t, resp, 201)

	spec := NewNeoPathFinderSpecWithRelationships(&NeoTraversalRelationship{Type: "likes", Direction: NeoTraversalOut})
	path, resp := service.FindPathFromNode(start, target, spec)
	if !resp.Ok() || resp.StatusCode != 200 {
		t.Fatalf("Unexpected response (%d): %v", resp.StatusCode, resp.NeoError)
	}

	if path.Length != 1 {
		t.Errorf("Expected the path length to be 1, but is %d", path.Length)
	}
	if len(path.Nodes) != 2 || path.Nodes[0] != start.Self.String() || path.Nodes[1] != target.Self.String() {
		t.Errorf("Expected %v for path nodes, but got %v", []string{start.Self.String(), target.Self.String()}, path.Nodes)
	}
	if len(path.Relationships) != 1 || path.Relationships[0] != rel.Self.String() {
		t.Errorf("Expected %v for path relationships, but got %v", []string{rel.Self.String()}, path.Relationships)
	}

	resp = service.DeleteRelationship(rel)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(start)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(target)
	checkResponseSucceeded(t, resp, 204)
}

func TestPathsFinder(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	start, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	target, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	rel, resp := service.CreateRelationshipWithPropertiesAndType(start, target, map[string]int{"cost": 1}, "likes")
	checkResponseSucceeded(t, resp, 201)

	spec := NewNeoPathFinderSpecWithRelationships(&NeoTraversalRelationship{Type: "likes", Direction: NeoTraversalOut})
	spec.Algorithm = NeoDijkstra
	spec.CostProperty = "cost"
	paths, resp := service.FindPathsFromNode(start, target, spec)
	if !resp.Ok() || resp.StatusCode != 200 {
		t.Fatalf("Unexpected response (%d): %v", resp.StatusCode, resp.NeoError)
	}

	if len(paths) != 1 {
		t.Errorf("Expected to get 1 path, but got %d", len(paths))
	}

	path := paths[0]

	if path.Length != 1 {
		t.Errorf("Expected the path length to be 1, but is %d", path.Length)
	}
	if len(path.Nodes) != 2 || path.Nodes[0] != start.Self.String() || path.Nodes[1] != target.Self.String() {
		t.Errorf("Expected %v for path nodes, but got %v", []string{start.Self.String(), target.Self.String()}, path.Nodes)
	}
	if len(path.Relationships) != 1 || path.Relationships[0] != rel.Self.String() {
		t.Errorf("Expected %v for path relationships, but got %v", []string{rel.Self.String()}, path.Relationships)
	}

	resp = service.DeleteRelationship(rel)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(start)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(target)
	checkResponseSucceeded(t, resp, 204)
}

func TestTraverseByNodes(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	start, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	middle, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	target, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	rel1, resp := service.CreateRelationshipWithType(start, middle, "likes")
	checkResponseSucceeded(t, resp, 201)

	rel2, resp := service.CreateRelationshipWithType(middle, target, "likes")
	checkResponseSucceeded(t, resp, 201)

	traversal := &NeoTraversal{}
	traversal.MaxDepth = 5
	traversal.ReturnFilter = NewNeoReturnFilterAllButStartNode()
	nodes, resp := service.TraverseByNodes(traversal, start)
	if !resp.Ok() || resp.StatusCode != 200 {
		t.Errorf("Unexpected server response (%d): %v", resp.StatusCode, resp.NeoError)
	}

	if len(nodes) != 2 {
		t.Errorf("Expected to get just 2 nodes, but got %d", len(nodes))
	}

	resp = service.DeleteRelationship(rel1)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteRelationship(rel2)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(start)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(middle)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(target)
	checkResponseSucceeded(t, resp, 204)
}

func TestTraverseByRelationships(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	start, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	middle, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	target, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	rel1, resp := service.CreateRelationshipWithType(start, middle, "likes")
	checkResponseSucceeded(t, resp, 201)

	rel2, resp := service.CreateRelationshipWithType(middle, target, "likes")
	checkResponseSucceeded(t, resp, 201)

	traversal := &NeoTraversal{}
	traversal.MaxDepth = 5
	rels, resp := service.TraverseByRelationships(traversal, start)
	if !resp.Ok() || resp.StatusCode != 200 {
		t.Errorf("Unexpected server response (%d): %v", resp.StatusCode, resp.NeoError)
	}

	if len(rels) != 2 {
		t.Errorf("Expected to get just 2 relationships, but got %d", len(rels))
	}

	if rels[0].Self.String() != rel1.Self.String() {
		t.Errorf("Expected to get the same relationship (%v), but got %v", rel1.Self.String(), rels[0].Self.String())
	}

	if rels[1].Self.String() != rel2.Self.String() {
		t.Errorf("Expected to get the same relationship (%v), but got %v", rel2.Self.String(), rels[1].Self.String())
	}

	resp = service.DeleteRelationship(rel1)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteRelationship(rel2)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(start)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(middle)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(target)
	checkResponseSucceeded(t, resp, 204)
}

func TestTraverseByPaths(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	start, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	middle, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	target, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	rel1, resp := service.CreateRelationshipWithType(start, middle, "likes")
	checkResponseSucceeded(t, resp, 201)

	rel2, resp := service.CreateRelationshipWithType(middle, target, "likes")
	checkResponseSucceeded(t, resp, 201)

	traversal := &NeoTraversal{}
	traversal.MaxDepth = 5
	paths, resp := service.TraverseByPaths(traversal, start)
	if !resp.Ok() || resp.StatusCode != 200 {
		t.Errorf("Unexpected server response (%d): %v", resp.StatusCode, resp.NeoError)
	}

	if len(paths) != 2 {
		t.Errorf("Expected to get just 2 paths, but got %d", len(paths))
	}

	resp = service.DeleteRelationship(rel1)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteRelationship(rel2)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(start)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(middle)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(target)
	checkResponseSucceeded(t, resp, 204)
}

func TestTraverseByFullPaths(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	start, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	middle, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	target, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	rel1, resp := service.CreateRelationshipWithType(start, middle, "likes")
	checkResponseSucceeded(t, resp, 201)

	rel2, resp := service.CreateRelationshipWithType(middle, target, "likes")
	checkResponseSucceeded(t, resp, 201)

	traversal := &NeoTraversal{}
	traversal.MaxDepth = 5
	paths, resp := service.TraverseByFullPaths(traversal, start)
	if !resp.Ok() || resp.StatusCode != 200 {
		t.Errorf("Unexpected server response (%d): %v", resp.StatusCode, resp.NeoError)
	}

	if len(paths) != 2 {
		t.Errorf("Expected to get just 2 paths, but got %d", len(paths))
	}

	resp = service.DeleteRelationship(rel1)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteRelationship(rel2)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(start)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(middle)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(target)
	checkResponseSucceeded(t, resp, 204)
}

func TestPagedTraverseByNodes(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	start, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	middle, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	target, resp := service.CreateNode()
	checkResponseSucceeded(t, resp, 201)

	rel1, resp := service.CreateRelationshipWithType(start, middle, "likes")
	checkResponseSucceeded(t, resp, 201)

	rel2, resp := service.CreateRelationshipWithType(middle, target, "likes")
	checkResponseSucceeded(t, resp, 201)

	traversal := &NeoTraversal{}
	traversal.MaxDepth = 5
	traversal.ReturnFilter = NewNeoReturnFilterAllButStartNode()

	// TODO paged traversal seems not supported, when used with streaming
	// traversal.PageSize = 1
	// _, _, resp = service.TraverseByNodesWithPaging(traversal, start)
	// if !resp.Ok() || resp.StatusCode != 200 {
	// 	t.Errorf("Unexpected server response (%d): %v", resp.StatusCode, resp.NeoError)
	// }

	// if len(nodes) != 1 {
	// 	t.Errorf("Expected to get just 1 node, but got %d", len(nodes))
	// }

	// nodes, resp = service.TraverseByNodesGetNextPage(traverser)
	// if !resp.Ok() || resp.StatusCode != 200 {
	// 	t.Errorf("Unexpected server response (%d): %v", resp.StatusCode, resp.NeoError)
	// }

	// if len(nodes) != 1 {
	// 	t.Errorf("Expected to get just 1 node, but got %d", len(nodes))
	// }

	// nodes, resp = service.TraverseByNodesGetNextPage(traverser)
	// if resp.Ok() || resp.StatusCode != 404 {
	// 	t.Errorf("Unexpected server response (%d): %v", resp.StatusCode, resp.NeoError)
	// }

	resp = service.DeleteRelationship(rel1)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteRelationship(rel2)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(start)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(middle)
	checkResponseSucceeded(t, resp, 204)

	resp = service.DeleteNode(target)
	checkResponseSucceeded(t, resp, 204)
}
