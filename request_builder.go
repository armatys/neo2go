package neo2go

import (
	"encoding/json"
)

type neoRequestBuilder struct {
	root *NeoServiceRoot
	self *UrlTemplate
}

func (n *neoRequestBuilder) Connect() (*NeoRequest, error) {
	return n.prepareRequest("GET", n.self.String(), nil, n.root, 200)
}

func (n *neoRequestBuilder) Cypher(cql string, params map[string]interface{}) (*CypherResponse, *NeoRequest, error) {
	bodyMap := map[string]interface{}{
		"query":  cql,
		"params": params,
	}

	url := n.root.Cypher.String()
	cypherResp := new(CypherResponse)
	req, err := n.prepareRequest("POST", url, &bodyMap, cypherResp, 200)
	return cypherResp, req, err
}

func (n *neoRequestBuilder) CreateNode() (*NeoNode, *NeoRequest, error) {
	return n.CreateNodeWithProperties(nil)
}

func (n *neoRequestBuilder) CreateNodeWithProperties(properties map[string]interface{}) (*NeoNode, *NeoRequest, error) {
	node := new(NeoNode)
	req, err := n.prepareRequest("POST", n.root.Node.String(), properties, node, 201)
	return node, req, err
}

func (n *neoRequestBuilder) DeleteNode(node *NeoNode) (*NeoRequest, error) {
	req, err := n.prepareRequest("DELETE", node.Self.String(), nil, nil, 204)
	return req, err
}

func (n *neoRequestBuilder) GetNode(uri string) (*NeoNode, *NeoRequest, error) {
	node := new(NeoNode)
	req, err := n.prepareRequest("GET", uri, nil, node, 200)
	return node, req, err
}

func (n *neoRequestBuilder) prepareRequest(method string, url string, body interface{}, result interface{}, expectedStatus int) (*NeoRequest, error) {
	var (
		bodyData []byte
		err      error
	)

	if body != nil {
		bodyData, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := NewNeoRequest(method, url, bodyData, result, expectedStatus)

	if err != nil {
		return nil, err
	}

	return req, nil
}
