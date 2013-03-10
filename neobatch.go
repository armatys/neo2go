package neo2go

import (
	"encoding/json"
	"fmt"
	"regexp"
)

//var _ Grapher = (*NeoBatch)(nil)

type NeoBatchId uint32

type neoBatchElement struct {
	Body   interface{} `json:"body"`
	Id     NeoBatchId  `json:"id"`
	Method string      `json:"method"`
	To     string      `json:"to"`
}

func (n *neoBatchElement) String() string {
	return fmt.Sprintf("<neoBatchElement[%d] %v %v>", n.Id, n.Method, n.To)
}

type NeoBatchResultElement struct {
	From     string
	Id       NeoBatchId
	Location string
	Body     interface{}
	Status   int
}

type NeoBatch struct {
	currentBatchId       NeoBatchId
	service              *GraphDatabaseService
	requests             []*NeoRequest
	requestBuilderErrors []error
	responses            []*NeoResponse
}

func (n *NeoBatch) nextBatchId() NeoBatchId {
	n.currentBatchId += 1
	return n.currentBatchId
}

func (n *NeoBatch) CreateNode() (*NeoNode, *NeoResponse) {
	result, req, err := n.service.builder.CreateNode()
	batchId := n.nextBatchId()
	req.batchId = batchId
	n.requests = append(n.requests, req)
	resp := &NeoResponse{req.expectedStatus, 0, nil}
	n.responses = append(n.responses, resp)

	result.batchId = batchId
	result.Self = NewPlainUrlTemplate(fmt.Sprintf("{%d}", batchId))

	if err != nil {
		n.requestBuilderErrors = append(n.requestBuilderErrors, err)
		resp.StatusCode = 600
		resp.NeoError = err
		return result, resp
	}

	return result, resp
}

func (n *NeoBatch) CreateNodeWithProperties(properties map[string]interface{}) (*NeoNode, *NeoResponse) {
	result, req, err := n.service.builder.CreateNodeWithProperties(properties)
	batchId := n.nextBatchId()
	result.batchId = batchId
	req.batchId = batchId
	n.requests = append(n.requests, req)
	resp := &NeoResponse{req.expectedStatus, 0, nil}
	n.responses = append(n.responses, resp)

	if err != nil {
		n.requestBuilderErrors = append(n.requestBuilderErrors, err)
		resp.StatusCode = 600
		resp.NeoError = err
		return result, resp
	}

	return result, resp
}

func (n *NeoBatch) DeleteNode(node *NeoNode) *NeoResponse {
	req, err := n.service.builder.DeleteNode(node)
	batchId := n.nextBatchId()
	req.batchId = batchId
	n.requests = append(n.requests, req)
	resp := &NeoResponse{req.expectedStatus, 0, nil}
	n.responses = append(n.responses, resp)

	if err != nil {
		n.requestBuilderErrors = append(n.requestBuilderErrors, err)
		resp.StatusCode = 600
		resp.NeoError = err
		return resp
	}

	return resp
}

func (n *NeoBatch) GetNode(uri string) (*NeoNode, *NeoResponse) {
	result, req, err := n.service.builder.GetNode(uri)
	batchId := n.nextBatchId()
	req.batchId = batchId
	result.batchId = batchId
	n.requests = append(n.requests, req)
	resp := &NeoResponse{req.expectedStatus, 0, nil}
	n.responses = append(n.responses, resp)

	if err != nil {
		n.requestBuilderErrors = append(n.requestBuilderErrors, err)
		resp.StatusCode = 600
		resp.NeoError = err
		return nil, resp
	}

	return result, resp
}

func (n *NeoBatch) Commit() *NeoResponse {
	expectedStatus := 200
	if n.currentBatchId == 0 {
		return &NeoResponse{expectedStatus, 600, fmt.Errorf("This batch does not contain any operations.")}
	}

	if len(n.requestBuilderErrors) > 0 {
		firstError := n.requestBuilderErrors[0]
		return &NeoResponse{expectedStatus, 600, fmt.Errorf("Errors during construction of requests: %v", firstError.Error())}
	}

	elements := make([]*neoBatchElement, len(n.requests))
	baseUrlLength := len(n.service.builder.self.String()) - 1
	for i, req := range n.requests {
		batchElem := new(neoBatchElement)
		//batchElem.Body = ""
		batchElem.Id = req.batchId
		batchElem.Method = req.Method

		re := regexp.MustCompile(`^%7B([0-9]+)%7D$`)
		if match := re.FindStringSubmatch(req.URL.String()); len(match) > 1 {
			batchElem.To = fmt.Sprintf("{%s}", match[1])
		} else {
			batchElem.To = req.URL.String()[baseUrlLength:]
		}
		elements[i] = batchElem
	}

	bodyData, err := json.Marshal(elements)
	if n.currentBatchId == 0 {
		return &NeoResponse{expectedStatus, 600, fmt.Errorf("Could not serialize batch element: %v", err.Error())}
	}

	results := make([]*NeoBatchResultElement, len(n.requests))
	for i, req := range n.requests {
		resultElem := new(NeoBatchResultElement)
		resultElem.Body = req.result
		results[i] = resultElem
	}

	neoRequest, err := NewNeoRequest("POST", n.service.builder.root.Batch.String(), bodyData, &results, 200)
	neoResponse := n.service.execute(neoRequest, err)

	// It is possible, (if neoResponse.Status is 200) to populate the responses
	// n.responses; with status codes, maybe errors

	for i, resultElem := range results {
		n.responses[i].StatusCode = resultElem.Status
		n.responses[i].ExpectedCode = n.requests[i].expectedStatus
	}

	return neoResponse
}
