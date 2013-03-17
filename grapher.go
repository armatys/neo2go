package neo2go

// Contains methods for manipulating the graph.
// (nodes, relationships and properties)
type Grapher interface {
	// 17.5.1
	CreateNode() (*NeoNode, *NeoResponse)
	// 17.5.2
	CreateNodeWithProperties(interface{}) (*NeoNode, *NeoResponse)

	// 17.5.5
	DeleteNode(node *NeoNode) *NeoResponse
	// 17.5.3
	GetNode(uri string) (*NeoNode, *NeoResponse)

	// 17.6.1
	GetRelationship(uri string) (*NeoRelationship, *NeoResponse)

	// 17.6.2
	CreateRelationshipWithType(source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse)
	// 17.6.3
	CreateRelationshipWithProperties(source *NeoNode, target *NeoNode, properties interface{}) (*NeoRelationship, *NeoResponse)
	// 17.6.3
	CreateRelationshipWithPropertiesAndType(source *NeoNode, target *NeoNode, properties interface{}, relType string) (*NeoRelationship, *NeoResponse)

	// 17.6.4
	DeleteRelationship(rel *NeoRelationship) *NeoResponse

	// 17.6.5
	GetPropertiesForRelationship(rel *NeoRelationship, result interface{}) *NeoResponse
	// 17.6.6
	ReplacePropertiesForRelationship(rel *NeoRelationship, properties interface{}) *NeoResponse
	// 17.6.7
	GetPropertyForRelationship(rel *NeoRelationship, propertyKey string, result interface{}) *NeoResponse
	// 17.6.8
	SetPropertyForRelationship(rel *NeoRelationship, propertyKey string, propertyValue interface{}) *NeoResponse

	// 17.6.9 - 17.6.11
	GetRelationshipsForNode(node *NeoNode, direction NeoTraversalDirection) ([]*NeoRelationship, *NeoResponse)

	// 17.6.12
	GetRelationshipsWithTypesForNode(node *NeoNode, direction NeoTraversalDirection, relTypes []string) ([]*NeoRelationship, *NeoResponse)

	// 17.7.1
	GetRelationshipTypes() ([]string, *NeoResponse)

	// ?
	GetPropertyForNode(node *NeoNode, propertyKey string, result interface{}) *NeoResponse
	// 17.8.1
	SetPropertyForNode(node *NeoNode, propertyKey string, propertyValue interface{}) *NeoResponse
	// 17.8.2
	ReplacePropertiesForNode(node *NeoNode, properties interface{}) *NeoResponse
	// 17.8.3
	GetPropertiesForNode(node *NeoNode, result interface{}) *NeoResponse
	// 17.8.6
	DeletePropertiesForNode(node *NeoNode) *NeoResponse
	// 17.8.7
	DeletePropertyWithKeyForNode(node *NeoNode, keyName string) *NeoResponse

	// 17.9.1
	UpdatePropertiesForRelationship(rel *NeoRelationship, properties interface{}) *NeoResponse
	// 17.9.2
	DeletePropertiesForRelationship(rel *NeoRelationship) *NeoResponse
	// 17.9.3
	DeletePropertyWithKeyForRelationship(rel *NeoRelationship, keyName string) *NeoResponse
}
