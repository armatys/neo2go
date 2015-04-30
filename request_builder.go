package neo2go

import (
	"fmt"
	"net/url"
)

type neoRequestBuilder struct {
	root     *NeoRoot
	dataRoot *NeoDataRoot
	self     *UrlTemplate
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

func (n *neoRequestBuilder) getRoot() *neoRequestData {
	return &neoRequestData{expectedStatus: 200, method: "GET", result: n.root, requestUrl: n.SelfReference()}
}

func (n *neoRequestBuilder) getDataRoot() *neoRequestData {
	return &neoRequestData{expectedStatus: 200, method: "GET", result: n.dataRoot, requestUrl: n.root.Data.String()}
}

func (n *neoRequestBuilder) Cypher(cql string, params map[string]interface{}) (*CypherResponse, *neoRequestData) {
	if params == nil {
		params = map[string]interface{}{}
	}
	bodyMap := map[string]interface{}{
		"query":  cql,
		"params": params,
	}

	url := n.dataRoot.Cypher.String()
	cypherResp := new(CypherResponse)
	requestData := neoRequestData{body: &bodyMap, expectedStatus: 200, method: "POST", result: cypherResp, requestUrl: url}
	return cypherResp, &requestData
}

// Transactional Cypher

func (n *neoRequestBuilder) TransactionalCypher(cypherTrans *CypherTransaction, commit bool, requests ...*CypherTransactionRequest) (*CypherTransaction, *neoRequestData) {
	statememts := make([]map[string]interface{}, 0, len(requests))
	for _, req := range requests {
		stmt := map[string]interface{}{
			"statement":          req.cql,
			"parameters":         req.params,
			"resultDataContents": []string{"REST"},
		}
		statememts = append(statememts, stmt)
	}

	var url string
	var expectedStatus int = 200
	returnedCypherTrans := new(CypherTransaction)

	if cypherTrans != nil && commit {
		url = cypherTrans.Commit.String()
	} else if cypherTrans != nil {
		url = cypherTrans.Self.String()
		returnedCypherTrans.Self = cypherTrans.Self
	} else {
		url = n.dataRoot.Transaction.String()
		if commit {
			url = url + "/commit"
		} else {
			expectedStatus = 201
		}
	}

	bodyMap := map[string]interface{}{
		"statements": statememts,
	}

	requestData := neoRequestData{body: &bodyMap, expectedStatus: expectedStatus, method: "POST", result: returnedCypherTrans, requestUrl: url}
	return returnedCypherTrans, &requestData
}

// Grapher

func (n *neoRequestBuilder) CreateNode() (*NeoNode, *neoRequestData) {
	return n.CreateNodeWithProperties(nil)
}

func (n *neoRequestBuilder) CreateNodeWithProperties(properties interface{}) (*NeoNode, *neoRequestData) {
	node := new(NeoNode)
	url := n.dataRoot.Node.String()
	requestData := neoRequestData{body: properties, expectedStatus: 201, method: "POST", result: node, requestUrl: url}
	return node, &requestData
}

func (n *neoRequestBuilder) DeleteNode(node *NeoNode) *neoRequestData {
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: node.Self.String()}
}

func (n *neoRequestBuilder) AddLabel(node *NeoNode, label string) *neoRequestData {
	return &neoRequestData{body: label, expectedStatus: 204, method: "POST", requestUrl: node.Labels.String()}
}

func (n *neoRequestBuilder) AddLabels(node *NeoNode, labels []string) *neoRequestData {
	return &neoRequestData{body: labels, expectedStatus: 204, method: "POST", requestUrl: node.Labels.String()}
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

func (n *neoRequestBuilder) CreateRelationshipWithType(source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *neoRequestData) {
	bodyMap := map[string]interface{}{
		"to":   target.Self.String(),
		"type": relType,
	}
	return n.createRelationshipHelper(source, bodyMap)
}

func (n *neoRequestBuilder) CreateRelationshipWithPropertiesAndType(source *NeoNode, target *NeoNode, properties interface{}, relType string) (*NeoRelationship, *neoRequestData) {
	bodyMap := map[string]interface{}{
		"to":   target.Self.String(),
		"type": relType,
		"data": properties,
	}
	return n.createRelationshipHelper(source, bodyMap)
}

func (n *neoRequestBuilder) DeleteRelationship(rel *NeoRelationship) *neoRequestData {
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: rel.Self.String()}
}

func (n *neoRequestBuilder) GetPropertiesForRelationship(rel *NeoRelationship) *neoRequestData {
	url := rel.Properties.String()
	requestData := neoRequestData{expectedStatus: 200, method: "GET", requestUrl: url}
	return &requestData
}

func (n *neoRequestBuilder) ReplacePropertiesForRelationship(rel *NeoRelationship, properties interface{}) *neoRequestData {
	url := rel.Properties.String()
	requestData := neoRequestData{body: properties, expectedStatus: 204, method: "PUT", requestUrl: url}
	return &requestData
}

func (n *neoRequestBuilder) GetPropertyForRelationship(rel *NeoRelationship, propertyKey string) (*neoRequestData, error) {
	url, err := rel.Property.Render(map[string]interface{}{"key": propertyKey})
	if err != nil {
		return nil, err
	}
	requestData := neoRequestData{expectedStatus: 200, method: "GET", requestUrl: url}
	return &requestData, nil
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

func (n *neoRequestBuilder) GetRelationshipsForNode(node *NeoNode, direction NeoTraversalDirection) (*[]*NeoRelationship, *neoRequestData) {
	var relationships []*NeoRelationship
	url := getRelationshipsUrlForNodeAndDirection(node, direction)
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: &relationships, requestUrl: url}
	return &relationships, &requestData
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

func (n *neoRequestBuilder) GetRelationshipsWithTypesForNode(node *NeoNode, direction NeoTraversalDirection, relTypes []string) (*[]*NeoRelationship, *neoRequestData, error) {
	var relationships []*NeoRelationship
	url, err := getTypedRelationshipsUrlForNodeAndDirection(node, direction, relTypes)
	if err != nil {
		return nil, nil, err
	}
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: &relationships, requestUrl: url}
	return &relationships, &requestData, nil
}

func (n *neoRequestBuilder) GetRelationshipTypes() (*[]string, *neoRequestData) {
	var relTypes []string
	url := n.dataRoot.RelationshipTypes.String()
	requestData := neoRequestData{expectedStatus: 200, method: "GET", result: &relTypes, requestUrl: url}
	return &relTypes, &requestData
}

func (n *neoRequestBuilder) SetPropertyForNode(node *NeoNode, propertyKey string, propertyValue interface{}) (*neoRequestData, error) {
	url, err := node.Property.Render(map[string]interface{}{"key": propertyKey})
	if err != nil {
		return nil, err
	}
	requestData := neoRequestData{body: propertyValue, expectedStatus: 204, method: "PUT", requestUrl: url}
	return &requestData, nil
}

func (n *neoRequestBuilder) ReplacePropertiesForNode(node *NeoNode, properties interface{}) *neoRequestData {
	url := node.Properties.String()
	return &neoRequestData{body: properties, expectedStatus: 204, method: "PUT", requestUrl: url}
}

func (n *neoRequestBuilder) GetPropertyForNode(node *NeoNode, propertyKey string) (*neoRequestData, error) {
	url, err := node.Property.Render(map[string]interface{}{"key": propertyKey})
	if err != nil {
		return nil, err
	}
	reqData := neoRequestData{expectedStatus: 200, method: "GET", requestUrl: url}
	return &reqData, nil
}

func (n *neoRequestBuilder) GetPropertiesForNode(node *NeoNode) *neoRequestData {
	url := node.Properties.String()
	reqData := neoRequestData{expectedStatus: 200, method: "GET", requestUrl: url}
	return &reqData
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

func (n *neoRequestBuilder) createNodeIndexHelper(params map[string]interface{}) (*NeoIndex, *neoRequestData) {
	var index *NeoIndex = new(NeoIndex)
	url := n.dataRoot.NodeIndex.String()
	return index, &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: index, requestUrl: url}
}

func (n *neoRequestBuilder) CreateNodeIndex(name string) (*NeoIndex, *neoRequestData) {
	params := map[string]interface{}{"name": name}
	return n.createNodeIndexHelper(params)
}

func (n *neoRequestBuilder) CreateNodeIndexWithConfiguration(name string, config interface{}) (*NeoIndex, *neoRequestData) {
	params := map[string]interface{}{"name": name, "config": config}
	return n.createNodeIndexHelper(params)
}

func (n *neoRequestBuilder) DeleteIndex(index *NeoIndex) (*neoRequestData, error) {
	url, err := index.Template.Render(nil)
	if err != nil {
		return nil, err
	}
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) GetNodeIndexes() (*map[string]*NeoIndex, *neoRequestData) {
	var result map[string]*NeoIndex
	url := n.dataRoot.NodeIndex.String()
	return &result, &neoRequestData{expectedStatus: 200, method: "GET", result: &result, requestUrl: url}
}

func (n *neoRequestBuilder) AddNodeToIndex(index *NeoIndex, node *NeoNode, key, value string) (*NeoNode, *neoRequestData, error) {
	var resultNode *NeoNode = new(NeoNode)
	url, err := index.Template.Render(nil)
	if err != nil {
		return nil, nil, err
	}
	params := map[string]string{
		"key":   key,
		"value": value,
		"uri":   node.Self.String(),
	}
	return resultNode, &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: resultNode, requestUrl: url}, nil
}

func deleteAllIndexEntriesForNodeHelper(index *NeoIndex, node *NeoNode, params map[string]interface{}) (*neoRequestData, error) {
	indexUrl, err := index.Template.Render(params)
	if err != nil {
		return nil, err
	}
	url := indexUrl + node.IdOrBatchId()
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForNode(index *NeoIndex, node *NeoNode) (*neoRequestData, error) {
	return deleteAllIndexEntriesForNodeHelper(index, node, nil)
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForNodeAndKey(index *NeoIndex, node *NeoNode, key string) (*neoRequestData, error) {
	return deleteAllIndexEntriesForNodeHelper(index, node, map[string]interface{}{"key": key})
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForNodeKeyAndValue(index *NeoIndex, node *NeoNode, key string, value string) (*neoRequestData, error) {
	return deleteAllIndexEntriesForNodeHelper(index, node, map[string]interface{}{"key": key, "value": value})
}

func (n *neoRequestBuilder) FindNodeByExactMatch(index *NeoIndex, key, value string) (*[]*NeoNode, *neoRequestData, error) {
	var nodes []*NeoNode
	indexUrl, err := index.Template.Render(map[string]interface{}{"key": key, "value": value})
	if err != nil {
		return nil, nil, err
	}
	return &nodes, &neoRequestData{expectedStatus: 200, method: "GET", result: &nodes, requestUrl: indexUrl}, nil
}

func (n *neoRequestBuilder) FindNodeByQuery(index *NeoIndex, query string) (*[]*NeoNode, *neoRequestData, error) {
	var nodes []*NeoNode
	indexUrl, err := index.Template.Render(nil)
	if err != nil {
		return nil, nil, err
	}
	url := fmt.Sprintf("%s?query=%s", indexUrl, url.QueryEscape(query))
	return &nodes, &neoRequestData{expectedStatus: 200, method: "GET", result: &nodes, requestUrl: url}, nil
}

func (n *neoRequestBuilder) createRelationshipIndexHelper(params map[string]interface{}) (*NeoIndex, *neoRequestData) {
	var index *NeoIndex = new(NeoIndex)
	url := n.dataRoot.RelationshipIndex.String()
	return index, &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: index, requestUrl: url}
}

func (n *neoRequestBuilder) CreateRelationshipIndex(name string) (*NeoIndex, *neoRequestData) {
	params := map[string]interface{}{"name": name}
	return n.createRelationshipIndexHelper(params)
}

func (n *neoRequestBuilder) CreateRelationshipIndexWithConfiguration(name string, config interface{}) (*NeoIndex, *neoRequestData) {
	params := map[string]interface{}{"name": name, "config": config}
	return n.createRelationshipIndexHelper(params)
}

func (n *neoRequestBuilder) GetRelationshipIndexes() (*map[string]*NeoIndex, *neoRequestData) {
	var result map[string]*NeoIndex
	url := n.dataRoot.RelationshipIndex.String()
	return &result, &neoRequestData{expectedStatus: 200, method: "GET", result: &result, requestUrl: url}
}

func (n *neoRequestBuilder) AddRelationshipToIndex(index *NeoIndex, rel *NeoRelationship, key, value string) (*NeoRelationship, *neoRequestData, error) {
	var resultRelationship *NeoRelationship = new(NeoRelationship)
	url, err := index.Template.Render(nil)
	if err != nil {
		return nil, nil, err
	}
	params := map[string]string{
		"key":   key,
		"value": value,
		"uri":   rel.Self.String(),
	}
	return resultRelationship, &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: resultRelationship, requestUrl: url}, nil
}

func deleteAllIndexEntriesForRelationshipHelper(index *NeoIndex, rel *NeoRelationship, params map[string]interface{}) (*neoRequestData, error) {
	indexUrl, err := index.Template.Render(params)
	if err != nil {
		return nil, err
	}
	url := indexUrl + rel.IdOrBatchId()
	return &neoRequestData{expectedStatus: 204, method: "DELETE", requestUrl: url}, nil
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForRelationship(index *NeoIndex, rel *NeoRelationship) (*neoRequestData, error) {
	return deleteAllIndexEntriesForRelationshipHelper(index, rel, nil)
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForRelationshipAndKey(index *NeoIndex, rel *NeoRelationship, key string) (*neoRequestData, error) {
	return deleteAllIndexEntriesForRelationshipHelper(index, rel, map[string]interface{}{"key": key})
}

func (n *neoRequestBuilder) DeleteAllIndexEntriesForRelationshipKeyAndValue(index *NeoIndex, rel *NeoRelationship, key string, value string) (*neoRequestData, error) {
	return deleteAllIndexEntriesForRelationshipHelper(index, rel, map[string]interface{}{"key": key, "value": value})
}

func (n *neoRequestBuilder) FindRelationshipByExactMatch(index *NeoIndex, key, value string) (*[]*NeoRelationship, *neoRequestData, error) {
	var rels []*NeoRelationship
	indexUrl, err := index.Template.Render(map[string]interface{}{"key": key, "value": value})
	if err != nil {
		return nil, nil, err
	}
	return &rels, &neoRequestData{expectedStatus: 200, method: "GET", result: &rels, requestUrl: indexUrl}, nil
}

func (n *neoRequestBuilder) FindRelationshipByQuery(index *NeoIndex, query string) (*[]*NeoRelationship, *neoRequestData, error) {
	var rels []*NeoRelationship
	indexUrl, err := index.Template.Render(nil)
	if err != nil {
		return nil, nil, err
	}
	url := fmt.Sprintf("%s?query=%s", indexUrl, url.QueryEscape(query))
	return &rels, &neoRequestData{expectedStatus: 200, method: "GET", result: &rels, requestUrl: url}, nil
}

func getOrCreateUniqueNodeHelper(index *NeoIndex, params map[string]interface{}) (*NeoNode, *neoRequestData, error) {
	var nodeResult *NeoNode = new(NeoNode)
	url, err := index.Template.Render(nil)
	if err != nil {
		return nil, nil, err
	}
	url = url + "?uniqueness=get_or_create"
	return nodeResult, &neoRequestData{body: params, expectedStatus: 200, method: "POST", result: nodeResult, requestUrl: url}, nil
}

func (n *neoRequestBuilder) GetOrCreateUniqueNode(index *NeoIndex, key, value string) (*NeoNode, *neoRequestData, error) {
	params := map[string]interface{}{
		"key":   key,
		"value": value,
	}
	return getOrCreateUniqueNodeHelper(index, params)
}

func (n *neoRequestBuilder) GetOrCreateUniqueNodeWithProperties(index *NeoIndex, key, value string, properties interface{}) (*NeoNode, *neoRequestData, error) {
	params := map[string]interface{}{
		"key":        key,
		"value":      value,
		"properties": properties,
	}
	return getOrCreateUniqueNodeHelper(index, params)
}

func createUniqueNodeOrFailHelper(index *NeoIndex, params map[string]interface{}) (*NeoNode, *neoRequestData, error) {
	var nodeResult *NeoNode = new(NeoNode)
	url, err := index.Template.Render(nil)
	if err != nil {
		return nil, nil, err
	}
	url = url + "?uniqueness=create_or_fail"
	return nodeResult, &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: nodeResult, requestUrl: url}, nil
}

func (n *neoRequestBuilder) CreateUniqueNodeOrFail(index *NeoIndex, key, value string) (*NeoNode, *neoRequestData, error) {
	params := map[string]interface{}{
		"key":   key,
		"value": value,
	}
	return createUniqueNodeOrFailHelper(index, params)
}

func (n *neoRequestBuilder) CreateUniqueNodeWithPropertiesOrFail(index *NeoIndex, key, value string, properties interface{}) (*NeoNode, *neoRequestData, error) {
	params := map[string]interface{}{
		"key":        key,
		"value":      value,
		"properties": properties,
	}
	return createUniqueNodeOrFailHelper(index, params)
}

func getOrCreateUniqueRelationshipHelper(index *NeoIndex, params map[string]interface{}) (*NeoRelationship, *neoRequestData, error) {
	var relResult *NeoRelationship = new(NeoRelationship)
	url, err := index.Template.Render(nil)
	if err != nil {
		return nil, nil, err
	}
	url = url + "?uniqueness=get_or_create"
	return relResult, &neoRequestData{body: params, expectedStatus: 200, method: "POST", result: relResult, requestUrl: url}, nil
}

func (n *neoRequestBuilder) GetOrCreateUniqueRelationship(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *neoRequestData, error) {
	params := map[string]interface{}{
		"key":   key,
		"value": value,
		"start": source.Self.String(),
		"end":   target.Self.String(),
		"type":  relType,
	}
	return getOrCreateUniqueRelationshipHelper(index, params)
}

func (n *neoRequestBuilder) GetOrCreateUniqueRelationshipWithProperties(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string, properties interface{}) (*NeoRelationship, *neoRequestData, error) {
	params := map[string]interface{}{
		"key":   key,
		"value": value,
		"start": source.Self.String(),
		"end":   target.Self.String(),
		"data":  properties,
		"type":  relType,
	}
	return getOrCreateUniqueRelationshipHelper(index, params)
}

func createUniqueRelationshipOrFailHelper(index *NeoIndex, params map[string]interface{}) (*NeoRelationship, *neoRequestData, error) {
	var relResult *NeoRelationship = new(NeoRelationship)
	url, err := index.Template.Render(nil)
	if err != nil {
		return nil, nil, err
	}
	url = url + "?uniqueness=create_or_fail"
	return relResult, &neoRequestData{body: params, expectedStatus: 201, method: "POST", result: relResult, requestUrl: url}, nil
}

func (n *neoRequestBuilder) CreateUniqueRelationshipOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *neoRequestData, error) {
	params := map[string]interface{}{
		"key":   key,
		"value": value,
		"start": source.Self.String(),
		"end":   target.Self.String(),
		"type":  relType,
	}
	return createUniqueRelationshipOrFailHelper(index, params)
}

func (n *neoRequestBuilder) CreateUniqueRelationshipWithPropertiesOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string, properties interface{}) (*NeoRelationship, *neoRequestData, error) {
	params := map[string]interface{}{
		"key":   key,
		"value": value,
		"start": source.Self.String(),
		"end":   target.Self.String(),
		"data":  properties,
		"type":  relType,
	}
	return createUniqueRelationshipOrFailHelper(index, params)
}

// GraphTraverser

func traverseHelper(traversal *NeoTraversal, start *NeoNode, params map[string]interface{}, result interface{}) (*neoRequestData, error) {
	url, err := start.Traverse.Render(params)
	if err != nil {
		return nil, err
	}
	return &neoRequestData{body: traversal, expectedStatus: 200, method: "POST", result: result, requestUrl: url}, nil
}

func (n *neoRequestBuilder) TraverseByNodes(traversal *NeoTraversal, start *NeoNode) (*[]*NeoNode, *neoRequestData, error) {
	var result []*NeoNode
	reqData, err := traverseHelper(traversal, start, map[string]interface{}{"returnType": "node"}, &result)
	return &result, reqData, err
}

func (n *neoRequestBuilder) TraverseByRelationships(traversal *NeoTraversal, start *NeoNode) (*[]*NeoRelationship, *neoRequestData, error) {
	var result []*NeoRelationship
	reqData, err := traverseHelper(traversal, start, map[string]interface{}{"returnType": "relationship"}, &result)
	return &result, reqData, err
}

func (n *neoRequestBuilder) TraverseByPaths(traversal *NeoTraversal, start *NeoNode) (*[]*NeoPath, *neoRequestData, error) {
	var result []*NeoPath
	reqData, err := traverseHelper(traversal, start, map[string]interface{}{"returnType": "path"}, &result)
	return &result, reqData, err
}

func (n *neoRequestBuilder) TraverseByFullPaths(traversal *NeoTraversal, start *NeoNode) (*[]*NeoFullPath, *neoRequestData, error) {
	var result []*NeoFullPath
	reqData, err := traverseHelper(traversal, start, map[string]interface{}{"returnType": "fullpath"}, &result)
	return &result, reqData, err
}

func pagedTraverseHelper(traversal *NeoTraversal, start *NeoNode, params map[string]interface{}, result interface{}) (*neoRequestData, error) {
	if traversal.LeaseTime > 0 {
		params["leaseTime"] = fmt.Sprintf("%d", traversal.LeaseTime)
	}
	if traversal.PageSize > 0 {
		params["pageSize"] = fmt.Sprintf("%d", traversal.PageSize)
	}
	url, err := start.PagedTraverse.Render(params)
	if err != nil {
		return nil, err
	}
	return &neoRequestData{body: traversal, expectedStatus: 201, method: "POST", result: result, requestUrl: url}, nil
}

func (n *neoRequestBuilder) TraverseByNodesWithPaging(traversal *NeoTraversal, start *NeoNode) (*[]*NeoNode, *neoRequestData, error) {
	var result []*NeoNode
	reqData, err := pagedTraverseHelper(traversal, start, map[string]interface{}{"returnType": "node"}, &result)
	return &result, reqData, err
}

func (n *neoRequestBuilder) TraverseByRelationshipsWithPaging(traversal *NeoTraversal, start *NeoNode) (*[]*NeoRelationship, *neoRequestData, error) {
	var result []*NeoRelationship
	reqData, err := pagedTraverseHelper(traversal, start, map[string]interface{}{"returnType": "relationship"}, &result)
	return &result, reqData, err
}

func (n *neoRequestBuilder) TraverseByPathsWithPaging(traversal *NeoTraversal, start *NeoNode) (*[]*NeoPath, *neoRequestData, error) {
	var result []*NeoPath
	reqData, err := pagedTraverseHelper(traversal, start, map[string]interface{}{"returnType": "path"}, &result)
	return &result, reqData, err
}

func (n *neoRequestBuilder) TraverseByFullPathsWithPaging(traversal *NeoTraversal, start *NeoNode) (*[]*NeoFullPath, *neoRequestData, error) {
	var result []*NeoFullPath
	reqData, err := pagedTraverseHelper(traversal, start, map[string]interface{}{"returnType": "fullpath"}, &result)
	return &result, reqData, err
}

func (n *neoRequestBuilder) TraverseByNodesGetNextPage(traverser *NeoPagedTraverser) (*[]*NeoNode, *neoRequestData) {
	var result []*NeoNode
	return &result, &neoRequestData{expectedStatus: 200, method: "GET", result: &result, requestUrl: traverser.location}
}

func (n *neoRequestBuilder) TraverseByRelationshipsGetNextPage(traverser *NeoPagedTraverser) (*[]*NeoRelationship, *neoRequestData) {
	var result []*NeoRelationship
	return &result, &neoRequestData{expectedStatus: 200, method: "GET", result: &result, requestUrl: traverser.location}
}

func (n *neoRequestBuilder) TraverseByPathsGetNextPage(traverser *NeoPagedTraverser) (*[]*NeoPath, *neoRequestData) {
	var result []*NeoPath
	return &result, &neoRequestData{expectedStatus: 200, method: "GET", result: &result, requestUrl: traverser.location}
}

func (n *neoRequestBuilder) TraverseByFullPathsGetNextPage(traverser *NeoPagedTraverser) (*[]*NeoFullPath, *neoRequestData) {
	var result []*NeoFullPath
	return &result, &neoRequestData{expectedStatus: 200, method: "GET", result: &result, requestUrl: traverser.location}
}

// GraphPathFinder

func (n *neoRequestBuilder) FindPathFromNode(start *NeoNode, target *NeoNode, spec *NeoPathFinderSpec) (*NeoPath, *neoRequestData) {
	var result *NeoPath = new(NeoPath)
	url := start.Self.String()
	if url[len(url)-1] != '/' {
		url += "/path"
	} else {
		url += "path"
	}
	spec.To = target.Self.String()
	return result, &neoRequestData{body: spec, expectedStatus: 200, method: "POST", result: result, requestUrl: url}
}

func (n *neoRequestBuilder) FindPathsFromNode(start *NeoNode, target *NeoNode, spec *NeoPathFinderSpec) (*[]*NeoPath, *neoRequestData) {
	var result []*NeoPath
	url := start.Self.String()
	if url[len(url)-1] != '/' {
		url += "/paths"
	} else {
		url += "paths"
	}
	spec.To = target.Self.String()
	return &result, &neoRequestData{body: spec, expectedStatus: 200, method: "POST", result: &result, requestUrl: url}
}
