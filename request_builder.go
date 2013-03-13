package neo2go

import (
	"fmt"
	"net/url"
)

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

// Grapher

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

func (n *neoRequestBuilder) DeleteRelationship(rel *NeoRelationship) *neoRequestData {
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: rel.SelfReference()}
}

func (n *neoRequestBuilder) GetPropertiesForRelationship(rel *NeoRelationship) (map[string]interface{}, *neoRequestData) {
	var properties map[string]interface{}
	url := rel.Properties.String()
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: properties, requestUrl: url}
	return properties, &requestData
}

func (n *neoRequestBuilder) ReplacePropertiesForRelationship(rel *NeoRelationship, properties map[string]interface{}) *neoRequestData {
	url := rel.Properties.String()
	requestData := neoRequestData{body: properties, expectedStatus: 204, method: "PUT", requestUrl: url}
	return &requestData
}

func (n *neoRequestBuilder) GetPropertyForRelationship(rel *NeoRelationship, propertyKey string) (interface{}, *neoRequestData, error) {
	var propertyValue interface{}
	url, err := rel.Property.Render(map[string]interface{}{"key": propertyKey})
	if err != nil {
		return nil, nil, err
	}
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: propertyValue, requestUrl: url}
	return propertyValue, &requestData, nil
}

func (n *neoRequestBuilder) SetPropertyForRelationship(rel *NeoRelationship, propertyKey string, propertyValue interface{}) (*neoRequestData, error) {
	url, err := rel.Property.Render(map[string]interface{}{"key": propertyKey})
	if err != nil {
		return nil, err
	}
	requestData := neoRequestData{body: propertyValue, expectedStatus: 204, method: "PUT", requestUrl: url}
	return &requestData, nil
}

func getRelationshipsUrlForNodeAndDirection(node *NeoNode, direction NeoTraversalDirection) string {
	var template *UrlTemplate
	switch direction {
	case NeoTraversalAll:
		template = node.AllRelationships
	case NeoTraversalIn:
		template = node.IncomingRelationships
	case NeoTraversalOut:
		template = node.OutgoingRelationships
	}

	if template != nil {
		return template.String()
	}

	return ""
}

func (n *neoRequestBuilder) GetRelationshipsForNode(node *NeoNode, direction NeoTraversalDirection) ([]*NeoRelationship, *neoRequestData) {
	var relationships []*NeoRelationship
	url := getRelationshipsUrlForNodeAndDirection(node, direction)
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: relationships, requestUrl: url}
	return relationships, &requestData
}

func getTypedRelationshipsUrlForNodeAndDirection(node *NeoNode, direction NeoTraversalDirection, relTypes []string) (string, error) {
	var template *UrlTemplate
	switch direction {
	case NeoTraversalAll:
		template = node.AllTypedRelationships
	case NeoTraversalIn:
		template = node.IncomingTypedRelationships
	case NeoTraversalOut:
		template = node.OutgoingTypedRelationships
	}

	if template != nil {
		url, err := template.Render(map[string]interface{}{"types": relTypes})
		return url, err
	}

	return "", nil
}

func (n *neoRequestBuilder) GetRelationshipsWithTypesForNode(node *NeoNode, direction NeoTraversalDirection, relTypes []string) ([]*NeoRelationship, *neoRequestData, error) {
	var relationships []*NeoRelationship
	url, err := getTypedRelationshipsUrlForNodeAndDirection(node, direction, relTypes)
	if err != nil {
		return nil, nil, err
	}
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: relationships, requestUrl: url}
	return relationships, &requestData, nil
}

func (n *neoRequestBuilder) GetRelationshipTypes() ([]string, *neoRequestData) {
	var relTypes []string
	url := n.root.RelationshipTypes.String()
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: relTypes, requestUrl: url}
	return relTypes, &requestData
}

func (n *neoRequestBuilder) SetPropertyForNode(node *NeoNode, propertyKey string, propertyValue interface{}) (*neoRequestData, error) {
	url, err := node.Property.Render(map[string]interface{}{"key": propertyKey})
	if err != nil {
		return nil, err
	}
	requestData := neoRequestData{body: propertyValue, expectedStatus: 204, method: "PUT", requestUrl: url}
	return &requestData, nil
}

func (n *neoRequestBuilder) ReplacePropertiesForNode(node *NeoNode, properties map[string]interface{}) *neoRequestData {
	url := node.Properties.String()
	return &neoRequestData{body: properties, expectedStatus: 204, method: "PUT", requestUrl: url}
}

func (n *neoRequestBuilder) GetPropertiesForNode(node *NeoNode) (map[string]interface{}, *neoRequestData) {
	var properties map[string]interface{}
	url := node.Properties.String()
	reqData := neoRequestData{expectedStatus: 200, method: "GET", result: properties, requestUrl: url}
	return properties, &reqData
}

func (n *neoRequestBuilder) DeletePropertiesForNode(node *NeoNode) *neoRequestData {
	url := node.Properties.String()
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}
}

func (n *neoRequestBuilder) DeletePropertyWithKeyForNode(node *NeoNode, keyName string) (*neoRequestData, error) {
	url, err := node.Property.Render(map[string]interface{}{"key": keyName})
	if err != nil {
		return nil, err
	}
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) UpdatePropertiesForRelationship(rel *NeoRelationship, properties map[string]interface{}) *neoRequestData {
	url := rel.Properties.String()
	return &neoRequestData{body: properties, expectedStatus: 204, method: "PUT", requestUrl: url}
}

func (n *neoRequestBuilder) DeletePropertiesForRelationship(rel *NeoRelationship) *neoRequestData {
	url := rel.Properties.String()
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}
}

func (n *neoRequestBuilder) DeletePropertyWithKeyForRelationship(rel *NeoRelationship, keyName string) (*neoRequestData, error) {
	url, err := rel.Property.Render(map[string]interface{}{"key": keyName})
	if err != nil {
		return nil, err
	}
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

// GraphIndexer

func (n *neoRequestBuilder) CreateNodeIndex(name string) *neoRequestData {
	var index *NeoIndex
	url := n.root.NodeIndex.String()
	params := map[string]string{"name": name}
	return &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: index, requestUrl: url}
}

func (n *neoRequestBuilder) CreateNodeIndexWithConfiguration(name string, config map[string]interface{}) *neoRequestData {
	var index *NeoIndex
	url := n.root.NodeIndex.String()
	params := map[string]interface{}{"name": name, "config": config}
	return &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: index, requestUrl: url}
}

func (n *neoRequestBuilder) DeleteIndex(index *NeoIndex) (*neoRequestData, error) {
	url, err := index.Template.Render(nil)
	if err != nil {
		return nil, err
	}
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) GetNodeIndexes() *neoRequestData {
	url := n.root.NodeIndex.String()
	return &neoRequestData{expectedStatus: 200, method: "GET", requestUrl: url}
}

func (n *neoRequestBuilder) AddNodeToIndex(index *NeoIndex, node *NeoNode, key, value string) (*neoRequestData, error) {
	var resultNode *NeoNode
	url, err := index.Template.Render(nil)
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"key":   key,
		"value": value,
		"uri":   node.SelfReference(),
	}
	return &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: resultNode, requestUrl: url}, nil
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForNode(index *NeoIndex, node *NeoNode) (*neoRequestData, error) {
	indexUrl, err := index.Template.Render(nil)
	if err != nil {
		return nil, err
	}
	url := indexUrl + node.IdOrBatchId()
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForNodeAndKey(index *NeoIndex, node *NeoNode, key string) (*neoRequestData, error) {
	indexUrl, err := index.Template.Render(map[string]interface{}{"key": key})
	if err != nil {
		return nil, err
	}
	url := indexUrl + node.IdOrBatchId()
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForNodeKeyAndValue(index *NeoIndex, node *NeoNode, key string, value string) (*neoRequestData, error) {
	indexUrl, err := index.Template.Render(map[string]interface{}{"key": key, "value": value})
	if err != nil {
		return nil, err
	}
	url := indexUrl + node.IdOrBatchId()
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) FindNodeByExactMatch(index *NeoIndex, key, value string) (*neoRequestData, error) {
	var nodes []*NeoNode
	indexUrl, err := index.Template.Render(map[string]interface{}{"key": key, "value": value})
	if err != nil {
		return nil, err
	}
	return &neoRequestData{expectedStatus: 200, method: "GET", result: nodes, requestUrl: indexUrl}, nil
}

func (n *neoRequestBuilder) FindNodeByQuery(index *NeoIndex, query string) (*neoRequestData, error) {
	var nodes []*NeoNode
	indexUrl, err := index.Template.Render(nil)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?query=%s", indexUrl, url.QueryEscape(query))
	return &neoRequestData{expectedStatus: 200, method: "GET", result: nodes, requestUrl: url}, nil
}

func (n *neoRequestBuilder) CreateRelationshipIndex(name string) *neoRequestData {
	var index *NeoIndex
	url := n.root.RelationshipIndex.String()
	params := map[string]string{"name": name}
	return &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: index, requestUrl: url}
}

func (n *neoRequestBuilder) CreateRelationshipIndexWithConfiguration(name string, config map[string]interface{}) *neoRequestData {
	var index *NeoIndex
	url := n.root.RelationshipIndex.String()
	params := map[string]interface{}{"name": name, "config": config}
	return &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: index, requestUrl: url}
}

func (n *neoRequestBuilder) GetRelationshipIndexes() *neoRequestData {
	url := n.root.RelationshipIndex.String()
	return &neoRequestData{expectedStatus: 200, method: "GET", requestUrl: url}
}

func (n *neoRequestBuilder) AddRelationshipToIndex(index *NeoIndex, rel *NeoRelationship, key, value string) (*neoRequestData, error) {
	var resultRelationship *NeoRelationship
	url, err := index.Template.Render(nil)
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"key":   key,
		"value": value,
		"uri":   rel.SelfReference(),
	}
	return &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: resultRelationship, requestUrl: url}, nil
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForRelationship(index *NeoIndex, rel *NeoRelationship) (*neoRequestData, error) {
	indexUrl, err := index.Template.Render(nil)
	if err != nil {
		return nil, err
	}
	url := indexUrl + rel.IdOrBatchId()
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForRelationshipAndKey(index *NeoIndex, rel *NeoRelationship, key string) (*neoRequestData, error) {
	indexUrl, err := index.Template.Render(map[string]interface{}{"key": key})
	if err != nil {
		return nil, err
	}
	url := indexUrl + rel.IdOrBatchId()
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForRelationshipKeyAndValue(index *NeoIndex, rel *NeoRelationship, key string, value string) (*neoRequestData, error) {
	indexUrl, err := index.Template.Render(map[string]interface{}{"key": key, "value": value})
	if err != nil {
		return nil, err
	}
	url := indexUrl + rel.IdOrBatchId()
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) FindRelationshipByExactMatch(index *NeoIndex, key, value string) (*neoRequestData, error) {
	var rels []*NeoRelationship
	indexUrl, err := index.Template.Render(map[string]interface{}{"key": key, "value": value})
	if err != nil {
		return nil, err
	}
	return &neoRequestData{expectedStatus: 200, method: "GET", result: rels, requestUrl: indexUrl}, nil
}

func (n *neoRequestBuilder) FindRelationshipByQuery(index *NeoIndex, query string) (*neoRequestData, error) {
	var rels []*NeoRelationship
	indexUrl, err := index.Template.Render(nil)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?query=%s", indexUrl, url.QueryEscape(query))
	return &neoRequestData{expectedStatus: 200, method: "GET", result: rels, requestUrl: url}, nil
}

func (n *neoRequestBuilder) GetOrCreateUniqueNode(index *NeoIndex, key, value string) *neoRequestData {
	var nodeResult *NeoNode
	url := index.SelfReference() + "?uniqueness=get_or_create"
	params := map[string]interface{}{
		"key":   key,
		"value": value,
	}
	// TODO expectedStatus can be 200 or 201
	return &neoRequestData{body: params, expectedStatus: 200, method: "POST", result: nodeResult, requestUrl: url}
}

func (n *neoRequestBuilder) GetOrCreateUniqueNodeWithProperties(index *NeoIndex, key, value string, properties map[string]interface{}) *neoRequestData {
	var nodeResult *NeoNode
	url := index.SelfReference() + "?uniqueness=get_or_create"
	params := map[string]interface{}{
		"key":        key,
		"value":      value,
		"properties": properties,
	}
	// TODO expectedStatus can be 200 or 201
	return &neoRequestData{body: params, expectedStatus: 200, method: "POST", result: nodeResult, requestUrl: url}
}

func (n *neoRequestBuilder) CreateUniqueNodeOrFail(index *NeoIndex, key, value string) *neoRequestData {
	var nodeResult *NeoNode
	url := index.SelfReference() + "?uniqueness=create_or_fail"
	params := map[string]interface{}{
		"key":   key,
		"value": value,
	}
	return &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: nodeResult, requestUrl: url}
}

func (n *neoRequestBuilder) CreateUniqueNodeWithPropertiesOrFail(index *NeoIndex, key, value string, properties map[string]interface{}) *neoRequestData {
	var nodeResult *NeoNode
	url := index.SelfReference() + "?uniqueness=create_or_fail"
	params := map[string]interface{}{
		"key":        key,
		"value":      value,
		"properties": properties,
	}
	return &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: nodeResult, requestUrl: url}
}

// 17.11.5
