package neo2go

import (
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

	req, err := g.builder.Connect()
	return g.execute_(req, err, false)
}

func (g *GraphDatabaseService) Batch() *NeoBatch {
	batch := new(NeoBatch)
	batch.service = g
	return batch
}

func (g *GraphDatabaseService) Cypher(cql string, params map[string]interface{}) (*CypherResponse, *NeoResponse) {
	result, req, err := g.builder.Cypher(cql, params)
	return result, g.execute(req, err)
}

// Grapher interface

func (g *GraphDatabaseService) CreateNode() (*NeoNode, *NeoResponse) {
	result, req, err := g.builder.CreateNode()
	return result, g.execute(req, err)
}

func (g *GraphDatabaseService) CreateNodeWithProperties(properties map[string]interface{}) (*NeoNode, *NeoResponse) {
	result, req, err := g.builder.CreateNodeWithProperties(properties)
	return result, g.execute(req, err)
}

func (g *GraphDatabaseService) DeleteNode(node *NeoNode) *NeoResponse {
	req, err := g.builder.DeleteNode(node)
	return g.execute(req, err)
}

func (g *GraphDatabaseService) GetNode(uri string) (*NeoNode, *NeoResponse) {
	result, req, err := g.builder.GetNode(uri)
	return result, g.execute(req, err)
}

func (g *GraphDatabaseService) execute(neoRequest *NeoRequest, err error) *NeoResponse {
	return g.execute_(neoRequest, err, true)
}

// Execute given request. If passed err is not nil, returns immediately with that error
// embedded inside NeoResponse.
// If the returned NeoResponse.StatuCode contains a 6xx, it means there was a local error
// while processing the request or response.
func (g *GraphDatabaseService) execute_(neoRequest *NeoRequest, err error, connRequired bool) *NeoResponse {
	expectedStatusCode := neoRequest.expectedStatus

	if connRequired && len(g.builder.root.Neo4jVersion) == 0 {
		return &NeoResponse{expectedStatusCode, 600, fmt.Errorf("Cannot execute the request because the client is not connected.")}
	}

	if err != nil {
		return &NeoResponse{expectedStatusCode, 600, err}
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
		container = neoRequest.result
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
