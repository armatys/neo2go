package neo2go

type neoRequestBuilder struct {
	root *NeoServiceRoot
	self *UrlTemplate
}

type neoRequestData struct {
	batchId        NeoBatchId
	body           interface{}
	expectedStatus int
	method         string
	result         interface{}
	requestUrl     string
}

func (n *neoRequestBuilder) Connect() *neoRequestData {
	return &neoRequestData{expectedStatus: 200, method: "GET", result: n.root, requestUrl: n.self.String()}
}

func (n *neoRequestBuilder) Cypher(cql string, params map[string]interface{}) (*CypherResponse, *neoRequestData) {
	bodyMap := map[string]interface{}{
		"query":  cql,
		"params": params,
	}

	url := n.root.Cypher.String()
	cypherResp := new(CypherResponse)
	requestData := neoRequestData{body: &bodyMap, expectedStatus: 200, method: "POST", result: cypherResp, requestUrl: url}
	return cypherResp, &requestData
}

func (n *neoRequestBuilder) CreateNode() (*NeoNode, *neoRequestData) {
	return n.CreateNodeWithProperties(nil)
}

func (n *neoRequestBuilder) CreateNodeWithProperties(properties map[string]interface{}) (*NeoNode, *neoRequestData) {
	node := new(NeoNode)
	url := n.root.Node.String()
	requestData := neoRequestData{body: &properties, expectedStatus: 201, method: "POST", result: node, requestUrl: url}
	return node, &requestData
}

func (n *neoRequestBuilder) DeleteNode(node *NeoNode) *neoRequestData {
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: node.Self.String()}
}

func (n *neoRequestBuilder) GetNode(nodeUrl string) (*NeoNode, *neoRequestData) {
	node := new(NeoNode)
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: node, requestUrl: nodeUrl}
	return node, &requestData
}
