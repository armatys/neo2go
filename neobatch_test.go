package neo2go

import (
	"testing"
)

func TestEmptyBatchShouldFail(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	batch := service.Batch()
	resp = batch.Commit()

	if !responseHasFailedWithCode(resp, 600) {
		t.Fatalf("Empty batch should fail, but didn't.")
	}
}

func TestBatchCreateDelete(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	batch := service.Batch()

	node, resp := batch.CreateNode()
	if resp != nil && resp.NeoError != nil {
		t.Fatalf("Error when creating batch: %v", resp.NeoError)
	}
	if node.batchId <= 0 {
		t.Fatalf("The `batchId` value should be greater than 0.")
	}

	batch.DeleteNode(node)
	resp = batch.Commit()

	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Batch did return an error (%d; %v): %v", resp.StatusCode, resp.Ok(), resp.NeoError)
	}

	_, resp = service.GetNode(node.Self.String())
	if !responseHasFailedWithCode(resp, 404) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}
}

func TestBatchCreateThenDelete(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	batch := service.Batch()

	node, resp := batch.CreateNode()
	if resp != nil && resp.NeoError != nil {
		t.Fatalf("Error when creating batch: %v", resp.NeoError)
	}
	if node.batchId <= 0 {
		t.Fatalf("The `batchId` value should be greater than 0.")
	}

	resp = batch.Commit()

	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Batch did return an error: %v", resp.NeoError)
	}

	_, resp = service.GetNode(node.Self.String())
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}

	resp = service.DeleteNode(node)
	if !responseHasSucceededWithCode(resp, 204) {
		t.Fatalf("Server returned unexpected response: %v", resp.NeoError.Error())
	}
}

func TestBatchDeleteNonExistentAndCreate(t *testing.T) {
	service := NewGraphDatabaseService()
	resp := service.Connect(databaseAddress)
	if !responseHasSucceededWithCode(resp, 200) {
		t.Fatalf("Error while connecting: %v", resp.NeoError.Error())
	}

	batch := service.Batch()
	node0 := new(NeoNode)
	node0.Self = NewPlainUrlTemplate(databaseAddress + "node/909090")
	batch.DeleteNode(node0)

	node, resp := batch.CreateNode()
	if resp.NeoError != nil {
		t.Fatalf("Error when creating batch: %v", resp.NeoError)
	}
	if node.batchId <= 0 {
		t.Fatalf("The `batchId` value should be greater than 0.")
	}

	resp = batch.Commit()

	if !responseHasFailedWithCode(resp, 600) {
		t.Fatalf("Batch did return an error: %v", resp.NeoError)
	}
}
