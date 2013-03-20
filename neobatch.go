package neo2go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
)

var _ Grapher = (*NeoBatch)(nil)
var _ GraphIndexer = (*NeoBatch)(nil)

var batchIdRegExp *regexp.Regexp

func init() {
	batchIdRegExp = regexp.MustCompile(`^{([0-9]+)}$`)
}

type NeoBatchId uint32

type batchIdSetter interface {
	setBatchId(batchId NeoBatchId)
}

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

func (n *NeoBatch) queueRequestData(reqData *neoRequestData) *NeoResponse {
	batchId := n.nextBatchId()
	reqData.setBatchId(batchId)

	resp := NewLocalErrorResponse(reqData.expectedStatus, nil)
	n.responses = append(n.responses, resp)
	n.requests = append(n.requests, reqData)

	return resp
}

func (n *NeoBatch) queueRequestDataWithResult(reqData *neoRequestData, result batchIdSetter) *NeoResponse {
	batchId := n.nextBatchId()
	reqData.setBatchId(batchId)
	result.setBatchId(batchId)

	resp := NewLocalErrorResponse(reqData.expectedStatus, nil)
	n.responses = append(n.responses, resp)
	n.requests = append(n.requests, reqData)

	return resp
}

func (n *NeoBatch) CreateNode() (*NeoNode, *NeoResponse) {
	result, reqData := n.service.builder.CreateNode()
	resp := n.queueRequestDataWithResult(reqData, result)
	return result, resp
}

func (n *NeoBatch) CreateNodeWithProperties(properties interface{}) (*NeoNode, *NeoResponse) {
	result, reqData := n.service.builder.CreateNodeWithProperties(properties)
	resp := n.queueRequestDataWithResult(reqData, result)
	return result, resp
}

func (n *NeoBatch) DeleteNode(node *NeoNode) *NeoResponse {
	reqData := n.service.builder.DeleteNode(node)
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) GetNode(uri string) (*NeoNode, *NeoResponse) {
	result, reqData := n.service.builder.GetNode(uri)
	resp := n.queueRequestDataWithResult(reqData, result)
	return result, resp
}

// Untested below

func (n *NeoBatch) GetRelationship(uri string) (*NeoRelationship, *NeoResponse) {
	result, reqData := n.service.builder.GetRelationship(uri)
	resp := n.queueRequestDataWithResult(reqData, result)
	return result, resp
}

func (n *NeoBatch) CreateRelationshipWithType(source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse) {
	result, reqData := n.service.builder.CreateRelationshipWithType(source, target, relType)
	resp := n.queueRequestDataWithResult(reqData, result)
	return result, resp
}

func (n *NeoBatch) CreateRelationshipWithProperties(source *NeoNode, target *NeoNode, properties interface{}) (*NeoRelationship, *NeoResponse) {
	result, reqData := n.service.builder.CreateRelationshipWithProperties(source, target, properties)
	resp := n.queueRequestDataWithResult(reqData, result)
	return result, resp
}

func (n *NeoBatch) CreateRelationshipWithPropertiesAndType(source *NeoNode, target *NeoNode, properties interface{}, relType string) (*NeoRelationship, *NeoResponse) {
	result, reqData := n.service.builder.CreateRelationshipWithPropertiesAndType(source, target, properties, relType)
	resp := n.queueRequestDataWithResult(reqData, result)
	return result, resp
}

//======

func (n *NeoBatch) DeleteRelationship(rel *NeoRelationship) *NeoResponse {
	reqData := n.service.builder.DeleteRelationship(rel)
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) GetPropertiesForRelationship(rel *NeoRelationship, result interface{}) *NeoResponse {
	reqData := n.service.builder.GetPropertiesForRelationship(rel)
	reqData.result = result
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) ReplacePropertiesForRelationship(rel *NeoRelationship, properties interface{}) *NeoResponse {
	reqData := n.service.builder.ReplacePropertiesForRelationship(rel, properties)
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) GetPropertyForRelationship(rel *NeoRelationship, propertyKey string, result interface{}) *NeoResponse {
	reqData, err := n.service.builder.GetPropertyForRelationship(rel, propertyKey)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	reqData.result = result
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) SetPropertyForRelationship(rel *NeoRelationship, propertyKey string, propertyValue interface{}) *NeoResponse {
	reqData, err := n.service.builder.SetPropertyForRelationship(rel, propertyKey, propertyValue)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) GetRelationshipsForNode(node *NeoNode, direction NeoTraversalDirection) ([]*NeoRelationship, *NeoResponse) {
	result, reqData := n.service.builder.GetRelationshipsForNode(node, direction)
	return *result, n.queueRequestData(reqData)
}

func (n *NeoBatch) GetRelationshipsWithTypesForNode(node *NeoNode, direction NeoTraversalDirection, relTypes []string) ([]*NeoRelationship, *NeoResponse) {
	result, reqData, err := n.service.builder.GetRelationshipsWithTypesForNode(node, direction, relTypes)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, n.queueRequestData(reqData)
}

func (n *NeoBatch) GetRelationshipTypes() ([]string, *NeoResponse) {
	result, reqData := n.service.builder.GetRelationshipTypes()
	return *result, n.queueRequestData(reqData)
}

func (n *NeoBatch) SetPropertyForNode(node *NeoNode, propertyKey string, propertyValue interface{}) *NeoResponse {
	reqData, err := n.service.builder.SetPropertyForNode(node, propertyKey, propertyValue)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) ReplacePropertiesForNode(node *NeoNode, properties interface{}) *NeoResponse {
	reqData := n.service.builder.ReplacePropertiesForNode(node, properties)
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) GetPropertyForNode(node *NeoNode, propertyKey string, propertyValueResult interface{}) *NeoResponse {
	reqData, err := n.service.builder.GetPropertyForNode(node, propertyKey)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	reqData.result = propertyValueResult
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) GetPropertiesForNode(node *NeoNode, result interface{}) *NeoResponse {
	reqData := n.service.builder.GetPropertiesForNode(node)
	reqData.result = result
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) DeletePropertiesForNode(node *NeoNode) *NeoResponse {
	reqData := n.service.builder.DeletePropertiesForNode(node)
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) DeletePropertyWithKeyForNode(node *NeoNode, keyName string) *NeoResponse {
	reqData, err := n.service.builder.DeletePropertyWithKeyForNode(node, keyName)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) DeletePropertiesForRelationship(rel *NeoRelationship) *NeoResponse {
	reqData := n.service.builder.DeletePropertiesForRelationship(rel)
	return n.queueRequestData(reqData)
}

func (n *NeoBatch) DeletePropertyWithKeyForRelationship(rel *NeoRelationship, keyName string) *NeoResponse {
	reqData, err := n.service.builder.DeletePropertyWithKeyForRelationship(rel, keyName)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

// GraphIndexer

// 17.10.1 - Nodes
func (n *NeoBatch) CreateNodeIndex(name string) (*NeoIndex, *NeoResponse) {
	result, reqData := n.service.builder.CreateNodeIndex(name)
	return result, n.queueRequestData(reqData)
}

// 17.10.2
func (n *NeoBatch) CreateNodeIndexWithConfiguration(name string, config interface{}) (*NeoIndex, *NeoResponse) {
	result, reqData := n.service.builder.CreateNodeIndexWithConfiguration(name, config)
	return result, n.queueRequestData(reqData)
}

// 17.10.3
func (n *NeoBatch) DeleteIndex(index *NeoIndex) *NeoResponse {
	reqData, err := n.service.builder.DeleteIndex(index)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

// 17.10.4
func (n *NeoBatch) GetNodeIndexes() (map[string]*NeoIndex, *NeoResponse) {
	result, reqData := n.service.builder.GetNodeIndexes()
	return *result, n.queueRequestData(reqData)
}

// 17.10.5
func (n *NeoBatch) AddNodeToIndex(index *NeoIndex, node *NeoNode, key, value string) (*NeoNode, *NeoResponse) {
	result, reqData, err := n.service.builder.AddNodeToIndex(index, node, key, value)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, n.queueRequestData(reqData)
}

// 17.10.6
func (n *NeoBatch) DeleteAllIndexEntriesForNode(index *NeoIndex, node *NeoNode) *NeoResponse {
	reqData, err := n.service.builder.DeleteAllIndexEntriesForNode(index, node)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

// 17.10.7
func (n *NeoBatch) DeleteAllIndexEntriesForNodeAndKey(index *NeoIndex, node *NeoNode, key string) *NeoResponse {
	reqData, err := n.service.builder.DeleteAllIndexEntriesForNodeAndKey(index, node, key)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

// 17.10.8
func (n *NeoBatch) DeleteAllIndexEntriesForNodeKeyAndValue(index *NeoIndex, node *NeoNode, key string, value string) *NeoResponse {
	reqData, err := n.service.builder.DeleteAllIndexEntriesForNodeKeyAndValue(index, node, key, value)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

// 17.10.9
func (n *NeoBatch) FindNodeByExactMatch(index *NeoIndex, key, value string) ([]*NeoNode, *NeoResponse) {
	result, reqData, err := n.service.builder.FindNodeByExactMatch(index, key, value)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, n.queueRequestData(reqData)
}

// 17.10.10
func (n *NeoBatch) FindNodeByQuery(index *NeoIndex, query string) ([]*NeoNode, *NeoResponse) {
	result, reqData, err := n.service.builder.FindNodeByQuery(index, query)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, n.queueRequestData(reqData)
}

// 17.10.1 - Relationships
func (n *NeoBatch) CreateRelationshipIndex(name string) (*NeoIndex, *NeoResponse) {
	result, reqData := n.service.builder.CreateRelationshipIndex(name)
	return result, n.queueRequestData(reqData)
}

// 17.10.2
func (n *NeoBatch) CreateRelationshipIndexWithConfiguration(name string, config interface{}) (*NeoIndex, *NeoResponse) {
	result, reqData := n.service.builder.CreateRelationshipIndexWithConfiguration(name, config)
	return result, n.queueRequestData(reqData)
}

// 17.10.4
func (n *NeoBatch) GetRelationshipIndexes() (map[string]*NeoIndex, *NeoResponse) {
	result, reqData := n.service.builder.GetRelationshipIndexes()
	return *result, n.queueRequestData(reqData)
}

// 17.10.5
func (n *NeoBatch) AddRelationshipToIndex(index *NeoIndex, rel *NeoRelationship, key, value string) (*NeoRelationship, *NeoResponse) {
	result, reqData, err := n.service.builder.AddRelationshipToIndex(index, rel, key, value)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, n.queueRequestData(reqData)
}

// 17.10.6
func (n *NeoBatch) DeleteAllIndexEntriesForRelationship(index *NeoIndex, rel *NeoRelationship) *NeoResponse {
	reqData, err := n.service.builder.DeleteAllIndexEntriesForRelationship(index, rel)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

// 17.10.7
func (n *NeoBatch) DeleteAllIndexEntriesForRelationshipAndKey(index *NeoIndex, rel *NeoRelationship, key string) *NeoResponse {
	reqData, err := n.service.builder.DeleteAllIndexEntriesForRelationshipAndKey(index, rel, key)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

// 17.10.8
func (n *NeoBatch) DeleteAllIndexEntriesForRelationshipKeyAndValue(index *NeoIndex, rel *NeoRelationship, key string, value string) *NeoResponse {
	reqData, err := n.service.builder.DeleteAllIndexEntriesForRelationshipKeyAndValue(index, rel, key, value)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return n.queueRequestData(reqData)
}

// 17.10.9
func (n *NeoBatch) FindRelationshipByExactMatch(index *NeoIndex, key, value string) ([]*NeoRelationship, *NeoResponse) {
	result, reqData, err := n.service.builder.FindRelationshipByExactMatch(index, key, value)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, n.queueRequestData(reqData)
}

// 17.10.10
func (n *NeoBatch) FindRelationshipByQuery(index *NeoIndex, query string) ([]*NeoRelationship, *NeoResponse) {
	result, reqData, err := n.service.builder.FindRelationshipByQuery(index, query)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, n.queueRequestData(reqData)
}

// 17.11.1
func (n *NeoBatch) GetOrCreateUniqueNode(index *NeoIndex, key, value string) (*NeoNode, *NeoResponse) {
	result, reqData := n.service.builder.GetOrCreateUniqueNode(index, key, value)
	return result, n.queueRequestData(reqData)
}

func (n *NeoBatch) GetOrCreateUniqueNodeWithProperties(index *NeoIndex, key, value string, properties interface{}) (*NeoNode, *NeoResponse) {
	result, reqData := n.service.builder.GetOrCreateUniqueNodeWithProperties(index, key, value, properties)
	return result, n.queueRequestData(reqData)
}

// 17.11.3
func (n *NeoBatch) CreateUniqueNodeOrFail(index *NeoIndex, key, value string) (*NeoNode, *NeoResponse) {
	result, reqData := n.service.builder.CreateUniqueNodeOrFail(index, key, value)
	return result, n.queueRequestData(reqData)
}

func (n *NeoBatch) CreateUniqueNodeWithPropertiesOrFail(index *NeoIndex, key, value string, properties interface{}) (*NeoNode, *NeoResponse) {
	result, reqData := n.service.builder.CreateUniqueNodeWithPropertiesOrFail(index, key, value, properties)
	return result, n.queueRequestData(reqData)
}

// 17.11.5
// func (n *NeoBatch) GetOrCreateUniqueRelationship(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode) (*NeoRelationship, *NeoResponse) {
// 	result, reqData := n.service.builder.GetOrCreateUniqueRelationship(index, key, value, source, target)
// 	return result, n.queueRequestData(reqData)
// }

func (n *NeoBatch) GetOrCreateUniqueRelationship(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse) {
	result, reqData := n.service.builder.GetOrCreateUniqueRelationship(index, key, value, source, target, relType)
	return result, n.queueRequestData(reqData)
}

// func (n *NeoBatch) GetOrCreateUniqueRelationshipWithProperties(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, properties interface{}) (*NeoRelationship, *NeoResponse) {
// 	result, reqData := n.service.builder.GetOrCreateUniqueRelationshipWithProperties(index, key, value, source, target, properties)
// 	return result, n.queueRequestData(reqData)
// }

func (n *NeoBatch) GetOrCreateUniqueRelationshipWithProperties(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string, properties interface{}) (*NeoRelationship, *NeoResponse) {
	result, reqData := n.service.builder.GetOrCreateUniqueRelationshipWithProperties(index, key, value, source, target, relType, properties)
	return result, n.queueRequestData(reqData)
}

// 17.11.7
// func (n *NeoBatch) CreateUniqueRelationshipOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode) (*NeoRelationship, *NeoResponse) {
// 	result, reqData := n.service.builder.CreateUniqueRelationshipOrFail(index, key, value, source, target)
// 	return result, n.queueRequestData(reqData)
// }

func (n *NeoBatch) CreateUniqueRelationshipOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse) {
	result, reqData := n.service.builder.CreateUniqueRelationshipOrFail(index, key, value, source, target, relType)
	return result, n.queueRequestData(reqData)
}

// func (n *NeoBatch) CreateUniqueRelationshipWithPropertiesOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, properties interface{}) (*NeoRelationship, *NeoResponse) {
// 	result, reqData := n.service.builder.CreateUniqueRelationshipWithPropertiesOrFail(index, key, value, source, target, properties)
// 	return result, n.queueRequestData(reqData)
// }

func (n *NeoBatch) CreateUniqueRelationshipWithPropertiesOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string, properties interface{}) (*NeoRelationship, *NeoResponse) {
	result, reqData := n.service.builder.CreateUniqueRelationshipWithPropertiesOrFail(index, key, value, source, target, relType, properties)
	return result, n.queueRequestData(reqData)
}

func (n *NeoBatch) Commit() *NeoResponse {
	expectedStatus := 200
	if n.currentBatchId == 0 {
		return NewLocalErrorResponse(expectedStatus, fmt.Errorf("This batch does not contain any operations."))
	}

	elements := make([]*neoBatchElement, len(n.requests))
	baseUrlLength := len(n.service.builder.self.String()) - 1
	for i, reqData := range n.requests {
		batchElem := new(neoBatchElement)
		batchElem.Body = reqData.body
		batchElem.Id = reqData.batchId
		batchElem.Method = reqData.method

		if matches := batchIdRegExp.MatchString(reqData.requestUrl); matches {
			batchElem.To = reqData.requestUrl
		} else if len(reqData.requestUrl) >= baseUrlLength {
			batchElem.To = reqData.requestUrl[baseUrlLength:]
		} else {
			return NewLocalErrorResponse(expectedStatus, fmt.Errorf("Unknown/badly formatted url: %v", reqData.requestUrl))
		}

		elements[i] = batchElem
	}

	bodyData, err := json.Marshal(elements)
	if err != nil {
		return NewLocalErrorResponse(expectedStatus, fmt.Errorf("Could not serialize batch element: %v", err.Error()))
	}
	bodyBuf := bytes.NewBuffer(bodyData)

	results := make([]*NeoBatchResultElement, len(n.requests))
	for i, reqData := range n.requests {
		resultElem := new(NeoBatchResultElement)
		resultElem.Body = reqData.result
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
