package neo2go

import (
	"encoding/json"
)

type NeoRequestBuilder struct {
	root *NeoServiceRoot
	self *UrlTemplate
}

func (n *NeoRequestBuilder) Connect() (*NeoRequest, error) {
	return n.prepareRequest("GET", n.self.String(), nil, n.root)
}

func (n *NeoRequestBuilder) Cypher(cql string, params []*NeoProperty) (*CypherResponse, *NeoRequest, error) {
	bodyMap := map[string]interface{}{
		"query":  cql,
		"params": params,
	}

	url := n.root.Cypher.String()
	cypherResp := new(CypherResponse)
	req, err := n.prepareRequest("POST", url, &bodyMap, cypherResp)
	return cypherResp, req, err
}

func (n *NeoRequestBuilder) CreateNode() (*NeoNode, *NeoRequest, error) {
	return n.CreateNodeWithProperties(nil)
}

func (n *NeoRequestBuilder) CreateNodeWithProperties(properties []*NeoProperty) (*NeoNode, *NeoRequest, error) {
	node := new(NeoNode)
	req, err := n.prepareRequest("POST", n.root.Node.String(), properties, node)
	return node, req, err
}

func (n *NeoRequestBuilder) DeleteNode(node *NeoNode) (*NeoRequest, error) {
	req, err := n.prepareRequest("DELETE", node.Self.String(), nil, nil)
	return req, err
}

func (n *NeoRequestBuilder) GetNode(uri string) (*NeoNode, *NeoRequest, error) {
	node := new(NeoNode)
	req, err := n.prepareRequest("GET", uri, nil, node)
	return node, req, err
}

func (n *NeoRequestBuilder) prepareRequest(method string, url string, body interface{}, result interface{}) (*NeoRequest, error) {
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

	req, err := NewNeoRequest(method, url, bodyData, result)

	if err != nil {
		return nil, err
	}

	return req, nil
}
