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

func (n *neoRequestData) setBatchId(bid NeoBatchId) {
	n.batchId = bid
}

func (n *neoRequestBuilder) SelfReference() string {
	if n.self != nil {
		return n.self.String()
	}
	return ""
}

func (n *neoRequestBuilder) Connect() *neoRequestData {
	return &neoRequestData{expectedStatus: 200, method: "GET", result: n.root, requestUrl: n.SelfReference()}
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
	requestData := neoRequestData{body: properties, expectedStatus: 201, method: "POST", result: node, requestUrl: url}
	return node, &requestData
}

func (n *neoRequestBuilder) DeleteNode(node *NeoNode) *neoRequestData {
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: node.SelfReference()}
}

func (n *neoRequestBuilder) GetNode(nodeUrl string) (*NeoNode, *neoRequestData) {
	node := new(NeoNode)
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: node, requestUrl: nodeUrl}
	return node, &requestData
}

func (n *neoRequestBuilder) GetRelationship(uri string) (*NeoRelationship, *neoRequestData) {
	relationship := new(NeoRelationship)
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: relationship, requestUrl: uri}
	return relationship, &requestData
}

func (n *neoRequestBuilder) createRelationshipHelper(source *NeoNode, bodyMap map[string]interface{}) (*NeoRelationship, *neoRequestData) {
	relationship := new(NeoRelationship)
	url := source.CreateRelationship.String()
	requestData := neoRequestData{body: bodyMap, expectedStatus: 201, method: "POST", result: relationship, requestUrl: url}
	return relationship, &requestData
}

func (n *neoRequestBuilder) CreateRelationship(source *NeoNode, target *NeoNode) (*NeoRelationship, *neoRequestData) {
	bodyMap := map[string]interface{}{
		"to": target.SelfReference(),
	}
	return n.createRelationshipHelper(source, bodyMap)
}

func (n *neoRequestBuilder) CreateRelationshipWithType(source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *neoRequestData) {
	bodyMap := map[string]interface{}{
		"to":   target.SelfReference(),
		"type": relType,
	}
	return n.createRelationshipHelper(source, bodyMap)
}

func (n *neoRequestBuilder) CreateRelationshipWithProperties(source *NeoNode, target *NeoNode, properties map[string]interface{}) (*NeoRelationship, *neoRequestData) {
	bodyMap := map[string]interface{}{
		"to":   target.SelfReference(),
		"data": properties,
	}
	return n.createRelationshipHelper(source, bodyMap)
}

func (n *neoRequestBuilder) CreateRelationshipWithPropertiesAndType(source *NeoNode, target *NeoNode, properties map[string]interface{}, relType string) (*NeoRelationship, *neoRequestData) {
	bodyMap := map[string]interface{}{
		"to":   target.SelfReference(),
		"type": relType,
		"data": properties,
	}
	return n.createRelationshipHelper(source, bodyMap)
}
