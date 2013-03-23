package neo2go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

const (
	Version = 1
)

var _ Grapher = (*GraphDatabaseService)(nil)
var _ GraphIndexer = (*GraphDatabaseService)(nil)
var _ GraphPathFinder = (*GraphDatabaseService)(nil)
var _ GraphTraverser = (*GraphDatabaseService)(nil)

var jsonContentTypeRegExp *regexp.Regexp

func init() {
	jsonContentTypeRegExp = regexp.MustCompile(`^application/json.*`)
}

type GraphDatabaseService struct {
	client  *http.Client
	builder *neoRequestBuilder
}

func NewGraphDatabaseService() *GraphDatabaseService {
	service := GraphDatabaseService{
		client:  &http.Client{},
		builder: &neoRequestBuilder{new(NeoServiceRoot), &UrlTemplate{}},
	}
	return &service
}

func (g *GraphDatabaseService) Connect(url string) *NeoResponse {
	g.builder.self.template = url
	g.builder.self.parse()

	reqData := g.builder.Connect()
	req, err := g.httpRequestFromData(reqData)
	return g.execute_(req, err, reqData.expectedStatus, reqData.result, false)
}

func (g *GraphDatabaseService) Batch() *NeoBatch {
	batch := new(NeoBatch)
	batch.service = g
	return batch
}

func (g *GraphDatabaseService) Cypher(cql string, params map[string]interface{}) (*CypherResponse, *NeoResponse) {
	result, reqData := g.builder.Cypher(cql, params)
	return result, g.executeFromRequestData(reqData)
}

// Grapher interface

func (g *GraphDatabaseService) CreateNode() (*NeoNode, *NeoResponse) {
	result, reqData := g.builder.CreateNode()
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) CreateNodeWithProperties(properties interface{}) (*NeoNode, *NeoResponse) {
	result, reqData := g.builder.CreateNodeWithProperties(properties)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeleteNode(node *NeoNode) *NeoResponse {
	reqData := g.builder.DeleteNode(node)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetNode(uri string) (*NeoNode, *NeoResponse) {
	result, reqData := g.builder.GetNode(uri)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetRelationship(uri string) (*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.GetRelationship(uri)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) CreateRelationshipWithType(source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.CreateRelationshipWithType(source, target, relType)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) CreateRelationshipWithProperties(source *NeoNode, target *NeoNode, properties interface{}) (*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.CreateRelationshipWithProperties(source, target, properties)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) CreateRelationshipWithPropertiesAndType(source *NeoNode, target *NeoNode, properties interface{}, relType string) (*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.CreateRelationshipWithPropertiesAndType(source, target, properties, relType)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeleteRelationship(rel *NeoRelationship) *NeoResponse {
	reqData := g.builder.DeleteRelationship(rel)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetPropertiesForRelationship(rel *NeoRelationship, result interface{}) *NeoResponse {
	reqData := g.builder.GetPropertiesForRelationship(rel)
	reqData.result = result
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) ReplacePropertiesForRelationship(rel *NeoRelationship, properties interface{}) *NeoResponse {
	reqData := g.builder.ReplacePropertiesForRelationship(rel, properties)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetPropertyForRelationship(rel *NeoRelationship, propertyKey string, result interface{}) *NeoResponse {
	reqData, err := g.builder.GetPropertyForRelationship(rel, propertyKey)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	reqData.result = result
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) SetPropertyForRelationship(rel *NeoRelationship, propertyKey string, propertyValue interface{}) *NeoResponse {
	reqData, err := g.builder.SetPropertyForRelationship(rel, propertyKey, propertyValue)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetRelationshipsForNode(node *NeoNode, direction NeoTraversalDirection) ([]*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.GetRelationshipsForNode(node, direction)
	return *result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetRelationshipsWithTypesForNode(node *NeoNode, direction NeoTraversalDirection, relTypes []string) ([]*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.GetRelationshipsWithTypesForNode(node, direction, relTypes)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetRelationshipTypes() ([]string, *NeoResponse) {
	result, reqData := g.builder.GetRelationshipTypes()
	return *result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) SetPropertyForNode(node *NeoNode, propertyKey string, propertyValue interface{}) *NeoResponse {
	reqData, err := g.builder.SetPropertyForNode(node, propertyKey, propertyValue)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) ReplacePropertiesForNode(node *NeoNode, properties interface{}) *NeoResponse {
	reqData := g.builder.ReplacePropertiesForNode(node, properties)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetPropertyForNode(node *NeoNode, propertyKey string, propertyValueResult interface{}) *NeoResponse {
	reqData, err := g.builder.GetPropertyForNode(node, propertyKey)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	reqData.result = propertyValueResult
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetPropertiesForNode(node *NeoNode, result interface{}) *NeoResponse {
	reqData := g.builder.GetPropertiesForNode(node)
	reqData.result = result
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeletePropertiesForNode(node *NeoNode) *NeoResponse {
	reqData := g.builder.DeletePropertiesForNode(node)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeletePropertyWithKeyForNode(node *NeoNode, keyName string) *NeoResponse {
	reqData, err := g.builder.DeletePropertyWithKeyForNode(node, keyName)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeletePropertiesForRelationship(rel *NeoRelationship) *NeoResponse {
	reqData := g.builder.DeletePropertiesForRelationship(rel)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeletePropertyWithKeyForRelationship(rel *NeoRelationship, keyName string) *NeoResponse {
	reqData, err := g.builder.DeletePropertyWithKeyForRelationship(rel, keyName)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

// GraphIndexer

// 17.10.1 - Nodes
func (g *GraphDatabaseService) CreateNodeIndex(name string) (*NeoIndex, *NeoResponse) {
	result, reqData := g.builder.CreateNodeIndex(name)
	return result, g.executeFromRequestData(reqData)
}

// 17.10.2
func (g *GraphDatabaseService) CreateNodeIndexWithConfiguration(name string, config interface{}) (*NeoIndex, *NeoResponse) {
	result, reqData := g.builder.CreateNodeIndexWithConfiguration(name, config)
	return result, g.executeFromRequestData(reqData)
}

// 17.10.3
func (g *GraphDatabaseService) DeleteIndex(index *NeoIndex) *NeoResponse {
	reqData, err := g.builder.DeleteIndex(index)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

// 17.10.4
func (g *GraphDatabaseService) GetNodeIndexes() (map[string]*NeoIndex, *NeoResponse) {
	result, reqData := g.builder.GetNodeIndexes()
	return *result, g.executeFromRequestData(reqData)
}

// 17.10.5
func (g *GraphDatabaseService) AddNodeToIndex(index *NeoIndex, node *NeoNode, key, value string) (*NeoNode, *NeoResponse) {
	result, reqData, err := g.builder.AddNodeToIndex(index, node, key, value)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, g.executeFromRequestData(reqData)
}

// 17.10.6
func (g *GraphDatabaseService) DeleteAllIndexEntriesForNode(index *NeoIndex, node *NeoNode) *NeoResponse {
	reqData, err := g.builder.DeleteAllIndexEntriesForNode(index, node)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

// 17.10.7
func (g *GraphDatabaseService) DeleteAllIndexEntriesForNodeAndKey(index *NeoIndex, node *NeoNode, key string) *NeoResponse {
	reqData, err := g.builder.DeleteAllIndexEntriesForNodeAndKey(index, node, key)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

// 17.10.8
func (g *GraphDatabaseService) DeleteAllIndexEntriesForNodeKeyAndValue(index *NeoIndex, node *NeoNode, key string, value string) *NeoResponse {
	reqData, err := g.builder.DeleteAllIndexEntriesForNodeKeyAndValue(index, node, key, value)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

// 17.10.9
func (g *GraphDatabaseService) FindNodeByExactMatch(index *NeoIndex, key, value string) ([]*NeoNode, *NeoResponse) {
	result, reqData, err := g.builder.FindNodeByExactMatch(index, key, value)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, g.executeFromRequestData(reqData)
}

// 17.10.10
func (g *GraphDatabaseService) FindNodeByQuery(index *NeoIndex, query string) ([]*NeoNode, *NeoResponse) {
	result, reqData, err := g.builder.FindNodeByQuery(index, query)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, g.executeFromRequestData(reqData)
}

// 17.10.1 - Relationships
func (g *GraphDatabaseService) CreateRelationshipIndex(name string) (*NeoIndex, *NeoResponse) {
	result, reqData := g.builder.CreateRelationshipIndex(name)
	return result, g.executeFromRequestData(reqData)
}

// 17.10.2
func (g *GraphDatabaseService) CreateRelationshipIndexWithConfiguration(name string, config interface{}) (*NeoIndex, *NeoResponse) {
	result, reqData := g.builder.CreateRelationshipIndexWithConfiguration(name, config)
	return result, g.executeFromRequestData(reqData)
}

// 17.10.4
func (g *GraphDatabaseService) GetRelationshipIndexes() (map[string]*NeoIndex, *NeoResponse) {
	result, reqData := g.builder.GetRelationshipIndexes()
	return *result, g.executeFromRequestData(reqData)
}

// 17.10.5
func (g *GraphDatabaseService) AddRelationshipToIndex(index *NeoIndex, rel *NeoRelationship, key, value string) (*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.AddRelationshipToIndex(index, rel, key, value)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, g.executeFromRequestData(reqData)
}

// 17.10.6
func (g *GraphDatabaseService) DeleteAllIndexEntriesForRelationship(index *NeoIndex, rel *NeoRelationship) *NeoResponse {
	reqData, err := g.builder.DeleteAllIndexEntriesForRelationship(index, rel)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

// 17.10.7
func (g *GraphDatabaseService) DeleteAllIndexEntriesForRelationshipAndKey(index *NeoIndex, rel *NeoRelationship, key string) *NeoResponse {
	reqData, err := g.builder.DeleteAllIndexEntriesForRelationshipAndKey(index, rel, key)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

// 17.10.8
func (g *GraphDatabaseService) DeleteAllIndexEntriesForRelationshipKeyAndValue(index *NeoIndex, rel *NeoRelationship, key string, value string) *NeoResponse {
	reqData, err := g.builder.DeleteAllIndexEntriesForRelationshipKeyAndValue(index, rel, key, value)
	if err != nil {
		return NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return g.executeFromRequestData(reqData)
}

// 17.10.9
func (g *GraphDatabaseService) FindRelationshipByExactMatch(index *NeoIndex, key, value string) ([]*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.FindRelationshipByExactMatch(index, key, value)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, g.executeFromRequestData(reqData)
}

// 17.10.10
func (g *GraphDatabaseService) FindRelationshipByQuery(index *NeoIndex, query string) ([]*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.FindRelationshipByQuery(index, query)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, g.executeFromRequestData(reqData)
}

// 17.11.1
func (g *GraphDatabaseService) GetOrCreateUniqueNode(index *NeoIndex, key, value string) (*NeoNode, *NeoResponse) {
	result, reqData, err := g.builder.GetOrCreateUniqueNode(index, key, value)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetOrCreateUniqueNodeWithProperties(index *NeoIndex, key, value string, properties interface{}) (*NeoNode, *NeoResponse) {
	result, reqData, err := g.builder.GetOrCreateUniqueNodeWithProperties(index, key, value, properties)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, g.executeFromRequestData(reqData)
}

// 17.11.3
func (g *GraphDatabaseService) CreateUniqueNodeOrFail(index *NeoIndex, key, value string) (*NeoNode, *NeoResponse) {
	result, reqData, err := g.builder.CreateUniqueNodeOrFail(index, key, value)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) CreateUniqueNodeWithPropertiesOrFail(index *NeoIndex, key, value string, properties interface{}) (*NeoNode, *NeoResponse) {
	result, reqData, err := g.builder.CreateUniqueNodeWithPropertiesOrFail(index, key, value, properties)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, g.executeFromRequestData(reqData)
}

// 17.11.5
func (g *GraphDatabaseService) GetOrCreateUniqueRelationship(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.GetOrCreateUniqueRelationship(index, key, value, source, target, relType)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetOrCreateUniqueRelationshipWithProperties(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string, properties interface{}) (*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.GetOrCreateUniqueRelationshipWithProperties(index, key, value, source, target, relType, properties)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, g.executeFromRequestData(reqData)
}

// 17.11.7
func (g *GraphDatabaseService) CreateUniqueRelationshipOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.CreateUniqueRelationshipOrFail(index, key, value, source, target, relType)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) CreateUniqueRelationshipWithPropertiesOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string, properties interface{}) (*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.CreateUniqueRelationshipWithPropertiesOrFail(index, key, value, source, target, relType, properties)
	if err != nil {
		return result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return result, g.executeFromRequestData(reqData)
}

// GraphTraverser

// 17.14.1+
func (g *GraphDatabaseService) TraverseByNodes(traversal *NeoTraversal, start *NeoNode) ([]*NeoNode, *NeoResponse) {
	result, reqData, err := g.builder.TraverseByNodes(traversal, start)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) TraverseByRelationships(traversal *NeoTraversal, start *NeoNode) ([]*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.TraverseByRelationships(traversal, start)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) TraverseByPaths(traversal *NeoTraversal, start *NeoNode) ([]*NeoPath, *NeoResponse) {
	result, reqData, err := g.builder.TraverseByPaths(traversal, start)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) TraverseByFullPaths(traversal *NeoTraversal, start *NeoNode) ([]*NeoFullPath, *NeoResponse) {
	result, reqData, err := g.builder.TraverseByFullPaths(traversal, start)
	if err != nil {
		return *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	return *result, g.executeFromRequestData(reqData)
}

// 17.14.5
func (g *GraphDatabaseService) TraverseByNodesWithPaging(traversal *NeoTraversal, start *NeoNode) (*NeoPagedTraverser, []*NeoNode, *NeoResponse) {
	result, reqData, err := g.builder.TraverseByNodesWithPaging(traversal, start)
	if err != nil {
		return nil, *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	response := g.executeFromRequestData(reqData)
	if len(response.location) > 0 {
		pagedTraverser := &NeoPagedTraverser{response.location}
		return pagedTraverser, *result, response
	}
	return nil, nil, NewLocalErrorResponse(reqData.expectedStatus, fmt.Errorf("The server did not return traverser's location."))
}

func (g *GraphDatabaseService) TraverseByRelationshipsWithPaging(traversal *NeoTraversal, start *NeoNode) (*NeoPagedTraverser, []*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.TraverseByRelationshipsWithPaging(traversal, start)
	if err != nil {
		return nil, *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	response := g.executeFromRequestData(reqData)
	if len(response.location) > 0 {
		pagedTraverser := &NeoPagedTraverser{response.location}
		return pagedTraverser, *result, response
	}
	return nil, nil, NewLocalErrorResponse(reqData.expectedStatus, fmt.Errorf("The server did not return traverser's location."))
}

func (g *GraphDatabaseService) TraverseByPathsWithPaging(traversal *NeoTraversal, start *NeoNode) (*NeoPagedTraverser, []*NeoPath, *NeoResponse) {
	result, reqData, err := g.builder.TraverseByPathsWithPaging(traversal, start)
	if err != nil {
		return nil, *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	response := g.executeFromRequestData(reqData)
	if len(response.location) > 0 {
		pagedTraverser := &NeoPagedTraverser{response.location}
		return pagedTraverser, *result, response
	}
	return nil, nil, NewLocalErrorResponse(reqData.expectedStatus, fmt.Errorf("The server did not return traverser's location."))
}

func (g *GraphDatabaseService) TraverseByFullPathsWithPaging(traversal *NeoTraversal, start *NeoNode) (*NeoPagedTraverser, []*NeoFullPath, *NeoResponse) {
	result, reqData, err := g.builder.TraverseByFullPathsWithPaging(traversal, start)
	if err != nil {
		return nil, *result, NewLocalErrorResponse(reqData.expectedStatus, err)
	}
	response := g.executeFromRequestData(reqData)
	if len(response.location) > 0 {
		pagedTraverser := &NeoPagedTraverser{response.location}
		return pagedTraverser, *result, response
	}
	return nil, nil, NewLocalErrorResponse(reqData.expectedStatus, fmt.Errorf("The server did not return traverser's location."))
}

// 17.14.6+
func (g *GraphDatabaseService) TraverseByNodesGetNextPage(traverser *NeoPagedTraverser) ([]*NeoNode, *NeoResponse) {
	result, reqData := g.builder.TraverseByNodesGetNextPage(traverser)
	return *result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) TraverseByRelationshipsGetNextPage(traverser *NeoPagedTraverser) ([]*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.TraverseByRelationshipsGetNextPage(traverser)
	return *result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) TraverseByPathsGetNextPage(traverser *NeoPagedTraverser) ([]*NeoPath, *NeoResponse) {
	result, reqData := g.builder.TraverseByPathsGetNextPage(traverser)
	return *result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) TraverseByFullPathsGetNextPage(traverser *NeoPagedTraverser) ([]*NeoFullPath, *NeoResponse) {
	result, reqData := g.builder.TraverseByFullPathsGetNextPage(traverser)
	return *result, g.executeFromRequestData(reqData)
}

// GraphPathFinder

func (g *GraphDatabaseService) FindPathFromNode(start *NeoNode, target *NeoNode, spec *NeoPathFinderSpec) (*NeoPath, *NeoResponse) {
	result, reqData := g.builder.FindPathFromNode(start, target, spec)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) FindPathsFromNode(start *NeoNode, target *NeoNode, spec *NeoPathFinderSpec) ([]*NeoPath, *NeoResponse) {
	result, reqData := g.builder.FindPathsFromNode(start, target, spec)
	return *result, g.executeFromRequestData(reqData)
}

// Utility methods

func (g *GraphDatabaseService) httpRequestFromData(reqData *neoRequestData) (*NeoHttpRequest, error) {
	var bodyBuffer *bytes.Buffer = nil

	if reqData.body != nil {
		bodyData, err := json.Marshal(reqData.body)
		if err != nil {
			return nil, err
		}
		bodyBuffer = bytes.NewBuffer(bodyData)
		return NewNeoHttpRequest(reqData.method, reqData.requestUrl, bodyBuffer)
	}

	return NewNeoHttpRequest(reqData.method, reqData.requestUrl, nil)
}

func (g *GraphDatabaseService) executeFromRequestData(reqData *neoRequestData) *NeoResponse {
	req, err := g.httpRequestFromData(reqData)
	return g.execute(req, err, reqData.expectedStatus, reqData.result)
}

func (g *GraphDatabaseService) execute(neoRequest *NeoHttpRequest, neoRequestErr error, expectedStatus int, result interface{}) *NeoResponse {
	return g.execute_(neoRequest, neoRequestErr, expectedStatus, result, true)
}

// Execute given request. If passed err is not nil, returns immediately with that error
// embedded inside NeoResponse.
// If the returned NeoResponse.StatuCode contains a 6xx, it means there was a local error
// while processing the request or response.
func (g *GraphDatabaseService) execute_(neoRequest *NeoHttpRequest, neoRequestErr error, expectedStatusCode int, result interface{}, connRequired bool) *NeoResponse {
	if connRequired && len(g.builder.root.Neo4jVersion) == 0 {
		return NewLocalErrorResponse(expectedStatusCode, fmt.Errorf("Cannot execute the request because the client is not connected."))
	}

	if neoRequestErr != nil {
		return NewLocalErrorResponse(expectedStatusCode, neoRequestErr)
	}

	resp, err := g.client.Do(neoRequest.Request)
	if err != nil {
		return NewLocalErrorResponse(expectedStatusCode, err)
	}

	defer resp.Body.Close()
	var container interface{}

	neoResponse := new(NeoResponse)
	neoResponse.ExpectedCode = expectedStatusCode
	neoResponse.StatusCode = resp.StatusCode
	locationUrl, err := resp.Location()
	if err == nil {
		neoResponse.location = locationUrl.String()
	}
	if resp.StatusCode >= 400 {
		neoErr := &NeoError{}
		container = neoErr
		neoResponse.NeoError = neoErr
	} else {
		container = result
	}

	if container != nil {
		ctype := resp.Header.Get("content-type")
		matched := jsonContentTypeRegExp.MatchString(ctype)

		if matched {
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(container)
			if err != nil {
				return NewLocalErrorResponse(expectedStatusCode, err)
			}
		} else if len(ctype) == 0 {
			//return NewLocalErrorResponse(expectedStatusCode, fmt.Errorf("Server did not return a content-type for this response."))
		} else {
			err := fmt.Errorf("Server has returned a response with unsupported content-type (%s)", ctype)
			return NewLocalErrorResponse(expectedStatusCode, err)
		}
	}

	return neoResponse
}
