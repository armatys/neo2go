package neo2go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

var _ Grapher = (*GraphDatabaseService)(nil)

//var _ GraphIndexer = (*GraphDatabaseService)(nil)

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
	err := g.builder.self.parse()
	if err != nil {
		return &NeoResponse{0, 600, err}
	}

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

func (g *GraphDatabaseService) CreateNodeWithProperties(properties map[string]interface{}) (*NeoNode, *NeoResponse) {
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

// Untested below

func (g *GraphDatabaseService) GetRelationship(uri string) (*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.GetRelationship(uri)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) CreateRelationship(source *NeoNode, target *NeoNode) (*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.CreateRelationship(source, target)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) CreateRelationshipWithType(source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.CreateRelationshipWithType(source, target, relType)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) CreateRelationshipWithProperties(source *NeoNode, target *NeoNode, properties map[string]interface{}) (*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.CreateRelationshipWithProperties(source, target, properties)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) CreateRelationshipWithPropertiesAndType(source *NeoNode, target *NeoNode, properties map[string]interface{}, relType string) (*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.CreateRelationshipWithPropertiesAndType(source, target, properties, relType)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeleteRelationship(rel *NeoRelationship) *NeoResponse {
	reqData := g.builder.DeleteRelationship(rel)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetPropertiesForRelationship(rel *NeoRelationship) (map[string]interface{}, *NeoResponse) {
	result, reqData := g.builder.GetPropertiesForRelationship(rel)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) ReplacePropertiesForRelationship(rel *NeoRelationship, properties map[string]interface{}) *NeoResponse {
	reqData := g.builder.ReplacePropertiesForRelationship(rel, properties)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetPropertyForRelationship(rel *NeoRelationship, propertyKey string) (interface{}, *NeoResponse) {
	result, reqData, err := g.builder.GetPropertyForRelationship(rel, propertyKey)
	if err != nil {
		return result, &NeoResponse{reqData.expectedStatus, 600, err}
	}
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) SetPropertyForRelationship(rel *NeoRelationship, propertyKey string, propertyValue interface{}) *NeoResponse {
	reqData, err := g.builder.SetPropertyForRelationship(rel, propertyKey, propertyValue)
	if err != nil {
		return &NeoResponse{reqData.expectedStatus, 600, err}
	}
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetRelationshipsForNode(node *NeoNode, direction NeoTraversalDirection) ([]*NeoRelationship, *NeoResponse) {
	result, reqData := g.builder.GetRelationshipsForNode(node, direction)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetRelationshipsWithTypesForNode(node *NeoNode, direction NeoTraversalDirection, relTypes []string) ([]*NeoRelationship, *NeoResponse) {
	result, reqData, err := g.builder.GetRelationshipsWithTypesForNode(node, direction, relTypes)
	if err != nil {
		return result, &NeoResponse{reqData.expectedStatus, 600, err}
	}
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetRelationshipTypes() ([]string, *NeoResponse) {
	result, reqData := g.builder.GetRelationshipTypes()
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) SetPropertyForNode(node *NeoNode, propertyKey string, propertyValue interface{}) *NeoResponse {
	reqData, err := g.builder.SetPropertyForNode(node, propertyKey, propertyValue)
	if err != nil {
		return &NeoResponse{reqData.expectedStatus, 600, err}
	}
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) ReplacePropertiesForNode(node *NeoNode, properties map[string]interface{}) *NeoResponse {
	reqData := g.builder.ReplacePropertiesForNode(node, properties)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) GetPropertiesForNode(node *NeoNode) (map[string]interface{}, *NeoResponse) {
	result, reqData := g.builder.GetPropertiesForNode(node)
	return result, g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeletePropertiesForNode(node *NeoNode) *NeoResponse {
	reqData := g.builder.DeletePropertiesForNode(node)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeletePropertyWithKeyForNode(node *NeoNode, keyName string) *NeoResponse {
	reqData, err := g.builder.DeletePropertyWithKeyForNode(node, keyName)
	if err != nil {
		return &NeoResponse{reqData.expectedStatus, 600, err}
	}
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) UpdatePropertiesForRelationship(rel *NeoRelationship, properties map[string]interface{}) *NeoResponse {
	reqData := g.builder.UpdatePropertiesForRelationship(rel, properties)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeletePropertiesForRelationship(rel *NeoRelationship) *NeoResponse {
	reqData := g.builder.DeletePropertiesForRelationship(rel)
	return g.executeFromRequestData(reqData)
}

func (g *GraphDatabaseService) DeletePropertyWithKeyForRelationship(rel *NeoRelationship, keyName string) *NeoResponse {
	reqData, err := g.builder.DeletePropertyWithKeyForRelationship(rel, keyName)
	if err != nil {
		return &NeoResponse{reqData.expectedStatus, 600, err}
	}
	return g.executeFromRequestData(reqData)
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
		return &NeoResponse{expectedStatusCode, 600, fmt.Errorf("Cannot execute the request because the client is not connected.")}
	}

	if neoRequestErr != nil {
		return &NeoResponse{expectedStatusCode, 600, neoRequestErr}
	}

	resp, err := g.client.Do(neoRequest.Request)
	if err != nil {
		return &NeoResponse{expectedStatusCode, 600, err}
	}

	defer resp.Body.Close()
	var container interface{}

	neoResponse := new(NeoResponse)
	neoResponse.ExpectedCode = expectedStatusCode
	neoResponse.StatusCode = resp.StatusCode
	if resp.StatusCode >= 400 {
		neoErr := &NeoError{}
		container = neoErr
		neoResponse.NeoError = neoErr
	} else {
		container = result
	}

	if container != nil {
		ctype := resp.Header.Get("content-type")
		if matched, err := regexp.MatchString("^application/json.*", ctype); matched {
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(container)
			if err != nil {
				return &NeoResponse{expectedStatusCode, 600, err}
			}
		} else if len(ctype) == 0 {
			return &NeoResponse{expectedStatusCode, 600, &NeoError{"Server did not return a content-type for this response.", "", nil}}
		} else {
			errorMessage := fmt.Sprintf("Server has returned a response with unsupported content-type (%s)", ctype)
			return &NeoResponse{expectedStatusCode, 600, &NeoError{errorMessage, "", nil}}
		}
	}

	return neoResponse
}
