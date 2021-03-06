package neo2go

import (
	"testing"
)

func TestEmptyBatchShouldFail(t *testing.T) {
	service := getDefaultDb()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.Err.Error())
	}

	batch := service.Batch()
	resp = batch.Commit()

	if !responseHasFailedWithCode(resp, 600) {
		t.Fatalf("Empty batch should fail, but didn't.")
	}
}

func TestBatchCreateDelete(t *testing.T) {
	service := getDefaultDb()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.Err.Error())
	}

	batch := service.Batch()

	node, resp := batch.CreateNode()
	if resp != nil && resp.Err != nil {
		t.Fatalf("Error when creating batch: %v", resp.Err)
	}
	if node.batchId <= 0 {
		t.Fatalf("The `batchId` value should be greater than 0.")
	}

	batch.DeleteNode(node)
	resp = batch.Commit()

	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Batch did return an error (%d; %v): %v", resp.StatusCode, resp.Ok(), resp.Err)
	}

	_, resp = service.GetNode(node.Self.String())
	if !responseHasFailedWithCode(resp, 404) {
		t.Fatalf("Server returned unexpected response: %v", resp.Err.Error())
	}
}

func TestBatchCreateThenDelete(t *testing.T) {
	service := getDefaultDb()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.Err.Error())
	}

	batch := service.Batch()

	node, resp := batch.CreateNode()
	if resp != nil && resp.Err != nil {
		t.Fatalf("Error when creating batch: %v", resp.Err)
	}
	if node.batchId <= 0 {
		t.Fatalf("The `batchId` value should be greater than 0.")
	}

	resp = batch.Commit()

	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Batch did return an error: %v", resp.Err)
	}

	_, resp = service.GetNode(node.Self.String())
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Server returned unexpected response: %v", resp.Err.Error())
	}

	resp = service.DeleteNode(node)
	if !responseHasSucceededWithCode(resp, 204) {
		t.Fatalf("Server returned unexpected response: %v", resp.Err.Error())
	}
}

func TestBatchDeleteNonExistentAndCreate(t *testing.T) {
	service := getDefaultDb()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.Err.Error())
	}

	batch := service.Batch()
	node0 := new(NeoNode)
	node0.Self = NewUrlTemplate(databaseAddress + "node/909090")
	batch.DeleteNode(node0)

	node, resp := batch.CreateNode()
	if resp.Err != nil {
		t.Fatalf("Error when creating batch: %v", resp.Err)
	}
	if node.batchId <= 0 {
		t.Fatalf("The `batchId` value should be greater than 0.")
	}

	resp = batch.Commit()

	if !responseHasFailedWithCode(resp, 600) {
		t.Fatalf("Batch did return an error: %v", resp.Err)
	}
}

func TestBatchGetNodeIndexes(t *testing.T) {
	service := getDefaultDb()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.Err.Error())
	}

	batch := service.Batch()

	index, _ := batch.CreateNodeIndex("test")
	indexMap, _ := batch.GetNodeIndexes()

	resp = batch.Commit()

	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Batch did return an error (%d; %v): %v", resp.StatusCode, resp.Ok(), resp.Err)
	}

	if len(*indexMap) < 1 {
		t.Fatalf("Excepted to get at least 1 index, but got %d: %v", len(*indexMap), *indexMap)
	}

	resp = service.DeleteIndex(index)
	if !responseHasSucceededWithCode(resp, 204) {
		t.Fatalf("Could not delete the index [%d]: %v", resp.StatusCode, resp.Err)
	}
}

func TestBatchSetNodeProperty(t *testing.T) {
	service := getDefaultDb()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.Err.Error())
	}

	v1, v2 := "John", "Alice"
	node, resp := service.CreateNodeWithProperties(map[string]string{"name": v1})
	if !responseHasSucceededWithCode(resp, 201) {
		t.Fatalf("Cannot create node (%d; %v): %v", resp.StatusCode, resp.Ok(), resp.Err)
	}

	var data map[string]interface{}
	if err := node.ParseData(&data); err != nil {
		t.Fatalf("Expected to get a map[string]interface{}, but could not convert")
	}

	if v, ok := data["name"].(string); ok {
		if v != v1 {
			t.Fatalf("Invalid node property value: expected %v but got %v.")
		}
	} else {
		t.Fatalf("Invalid node property value: cannot convert to string (%v)", data["name"])
	}

	batch := service.Batch()
	node2, _ := batch.CreateNode()
	batch.SetPropertyForNode(node2, "name", v2)
	resp = batch.Commit()

	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Batch did return an error (%d; %v): %v", resp.StatusCode, resp.Ok(), resp.Err)
	}

	node2, resp = service.GetNode(node2.Self.String())
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Could not get the node (%d; %v): %v", resp.StatusCode, resp.Ok(), resp.Err)
	}

	var data2 map[string]interface{}
	if err := node2.ParseData(&data2); err != nil {
		t.Fatalf("Expected to get a map[string]interface{}, but could not convert")
	}

	if v, ok := data2["name"].(string); ok {
		if v != v2 {
			t.Fatalf("Invalid node2 property value: expected %v but got %v.", v2, v)
		}
	} else {
		t.Fatalf("Invalid node2 property value: cannot convert to string (%v)", data2["name"])
	}
}
