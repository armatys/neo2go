package neo2go

import (
	"bytes"
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
	Body     interface{}
	From     string
	Id       NeoBatchId
	Location string
	Status   int
}

type NeoBatch struct {
	// The `0` value indicates there are no operations in this batch. Otherwise, a batch operation id starts from `1`.
	currentBatchId NeoBatchId
	service        *GraphDatabaseService
	requests       []*neoRequestData
	responses      []*NeoResponse
}

func (n *NeoBatch) nextBatchId() NeoBatchId {
	n.currentBatchId += 1
	return n.currentBatchId
}

func (n *NeoBatch) CreateNode() (*NeoNode, *NeoResponse) {
	result, req := n.service.builder.CreateNode()
	batchId := n.nextBatchId()
	req.batchId = batchId
	n.requests = append(n.requests, req)
	resp := &NeoResponse{req.expectedStatus, 0, nil}
	n.responses = append(n.responses, resp)

	result.batchId = batchId
	result.Self = NewPlainUrlTemplate(fmt.Sprintf("{%d}", batchId))

	return result, resp
}

func (n *NeoBatch) CreateNodeWithProperties(properties map[string]interface{}) (*NeoNode, *NeoResponse) {
	result, req := n.service.builder.CreateNodeWithProperties(properties)
	batchId := n.nextBatchId()
	result.batchId = batchId
	req.batchId = batchId
	n.requests = append(n.requests, req)
	resp := &NeoResponse{req.expectedStatus, 0, nil}
	n.responses = append(n.responses, resp)

	return result, resp
}

func (n *NeoBatch) DeleteNode(node *NeoNode) *NeoResponse {
	req := n.service.builder.DeleteNode(node)
	batchId := n.nextBatchId()
	req.batchId = batchId
	n.requests = append(n.requests, req)
	resp := &NeoResponse{req.expectedStatus, 0, nil}
	n.responses = append(n.responses, resp)

	return resp
}

func (n *NeoBatch) GetNode(uri string) (*NeoNode, *NeoResponse) {
	result, req := n.service.builder.GetNode(uri)
	batchId := n.nextBatchId()
	req.batchId = batchId
	result.batchId = batchId
	n.requests = append(n.requests, req)
	resp := &NeoResponse{req.expectedStatus, 0, nil}
	n.responses = append(n.responses, resp)

	return result, resp
}

func (n *NeoBatch) Commit() *NeoResponse {
	expectedStatus := 200
	if n.currentBatchId == 0 {
		return &NeoResponse{expectedStatus, 600, fmt.Errorf("This batch does not contain any operations.")}
	}

	elements := make([]*neoBatchElement, len(n.requests))
	baseUrlLength := len(n.service.builder.self.String()) - 1
	for i, req := range n.requests {
		batchElem := new(neoBatchElement)
		batchElem.Body = req.body
		batchElem.Id = req.batchId
		batchElem.Method = req.method

		re := regexp.MustCompile(`^{([0-9]+)}$`)
		if match := re.FindStringSubmatch(req.requestUrl); len(match) > 1 {
			batchElem.To = fmt.Sprintf("{%s}", match[1])
		} else if len(req.requestUrl) >= baseUrlLength {
			batchElem.To = req.requestUrl[baseUrlLength:]
		} else {
			return &NeoResponse{expectedStatus, 600, fmt.Errorf("Unknown/badly formatted url: %v", req.requestUrl)}
		}
		elements[i] = batchElem
	}

	bodyData, err := json.Marshal(elements)
	if err != nil {
		return &NeoResponse{expectedStatus, 600, fmt.Errorf("Could not serialize batch element: %v", err.Error())}
	}
	bodyBuf := bytes.NewBuffer(bodyData)

	results := make([]*NeoBatchResultElement, len(n.requests))
	for i, req := range n.requests {
		resultElem := new(NeoBatchResultElement)
		resultElem.Body = req.result
		results[i] = resultElem
	}

	neoRequest, err := NewNeoHttpRequest("POST", n.service.builder.root.Batch.String(), bodyBuf)
	neoResponse := n.service.execute(neoRequest, err, 200, &results)

	for i, resultElem := range results {
		n.responses[i].StatusCode = resultElem.Status
		n.responses[i].ExpectedCode = n.requests[i].expectedStatus

		if resultElem.Status != n.requests[i].expectedStatus {
			neoResponse.StatusCode = 600
			neoResponse.NeoError = fmt.Errorf("The batch operation #%v has failed.", n.requests[i].batchId)
		}
	}

	return neoResponse
}
