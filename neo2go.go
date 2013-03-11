package neo2go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

//var _ Grapher = (*GraphDatabaseService)(nil)
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
	result, req := g.builder.Cypher(cql, params)
	return result, g.executeFromRequestData(req)
}

// Grapher interface

func (g *GraphDatabaseService) CreateNode() (*NeoNode, *NeoResponse) {
	result, req := g.builder.CreateNode()
	return result, g.executeFromRequestData(req)
}

func (g *GraphDatabaseService) CreateNodeWithProperties(properties map[string]interface{}) (*NeoNode, *NeoResponse) {
	result, req := g.builder.CreateNodeWithProperties(properties)
	return result, g.executeFromRequestData(req)
}

func (g *GraphDatabaseService) DeleteNode(node *NeoNode) *NeoResponse {
	req := g.builder.DeleteNode(node)
	return g.executeFromRequestData(req)
}

func (g *GraphDatabaseService) GetNode(uri string) (*NeoNode, *NeoResponse) {
	result, req := g.builder.GetNode(uri)
	return result, g.executeFromRequestData(req)
}

func (g *GraphDatabaseService) httpRequestFromData(reqData *neoRequestData) (*NeoHttpRequest, error) {
	var bodyBuffer *bytes.Buffer = nil

	if reqData.body != nil {
		bodyData, err := json.Marshal(reqData.body)
		if err != nil {
			return nil, err
		}
		bodyBuffer = bytes.NewBuffer(bodyData)
	} else {
		bodyBuffer = nil
	}

	req, err := NewNeoHttpRequest(reqData.method, reqData.requestUrl, bodyBuffer)
	return req, err
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
