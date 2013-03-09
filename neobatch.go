package neo2go

import (
	"encoding/json"
	"fmt"
)

//var _ Grapher = (*NeoBatch)(nil)

type NeoBatchId uint32

type neoBatchElement struct {
	body   interface{}
	id     NeoBatchId
	method string
	to     string
}

type NeoBatch struct {
	currentBatchId       NeoBatchId
	service              *GraphDatabaseService
	requests             []*NeoRequest
	requestBuilderErrors []error
	responses            []interface{}
}

func (n *NeoBatch) nextBatchId() NeoBatchId {
	n.currentBatchId += 1
	return n.currentBatchId
}

func (n *NeoBatch) CreateNode() (*NeoNode, *NeoResponse) {
	result, req, err := n.service.builder.CreateNode()
	batchId := n.nextBatchId()
	result.batchId = batchId
	req.batchId = batchId
	n.requests = append(n.requests, req)

	if err != nil {
		n.requestBuilderErrors = append(n.requestBuilderErrors, err)
		return result, &NeoResponse{600, err}
	}

	return result, nil
}

func (n *NeoBatch) Commit() *NeoResponse {
	if n.currentBatchId == 0 {
		return &NeoResponse{600, fmt.Errorf("This batch does not contain any operations.")}
	}

	if len(n.requestBuilderErrors) > 0 {
		firstError := n.requestBuilderErrors[0]
		return &NeoResponse{600, fmt.Errorf("Errors during construction of requests: %v", firstError.Error())}
	}

	elements := make([]*neoBatchElement, len(n.requests))
	baseUrlLength := len(n.service.builder.self.String()) - 1
	for _, req := range n.requests {
		batchElem := new(neoBatchElement)
		batchElem.body = nil
		batchElem.id = req.batchId
		batchElem.method = req.Method
		batchElem.to = req.RequestURI[baseUrlLength:]
		elements = append(elements, batchElem)
	}

	bodyData, err := json.Marshal(elements)
	if n.currentBatchId == 0 {
		return &NeoResponse{600, fmt.Errorf("Could not serialize batch element: %v", err.Error())}
	}

	results := make([]interface{}, len(n.requests))
	for i, req := range n.requests {
		results[i] = req.result
	}
	neoRequest, err := NewNeoRequest("POST", n.service.builder.root.Batch.String(), bodyData, results)
	neoResponse := n.service.execute(neoRequest, err)

	return neoResponse
}
