package neo2go

// Contains methods for manipulating the graph.
// (nodes, relationships and properties)
type Grapher interface {
	// 17.5.1
	CreateNode() (*NeoNode, *NeoResponse)
	// 17.5.2
	CreateNodeWithProperties(map[string]interface{}) (*NeoNode, *NeoResponse)

	// 17.5.5
	DeleteNode(node *NeoNode) *NeoResponse
	// 17.5.3
	GetNode(id int64) (*NeoNode, *NeoResponse)

	// 17.6.1
	GetRelationship(id int64) (*NeoRelationship, *NeoResponse)

	// 17.6.2
	CreateRelationship(source *NeoNode, target *NeoNode) (*NeoRelationship, *NeoResponse)
	// 17.6.2
	CreateRelationshipWithType(source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse)
	// 17.6.3
	CreateRelationshipWithProperties(source *NeoNode, target *NeoNode, properties map[string]interface{}) (*NeoRelationship, *NeoResponse)
	// 17.6.3
	CreateRelationshipWithPropertiesAndType(source *NeoNode, target *NeoNode, properties map[string]interface{}, relType string) (*NeoRelationship, *NeoResponse)

	// 17.6.4
	DeleteRelationship(*NeoRelationship) *NeoResponse

	// 17.6.5
	GetRelationshipProperties(*NeoRelationship) (map[string]interface{}, *NeoResponse)
	// 17.6.6
	ReplaceRelationshipProperties(*NeoRelationship, map[string]interface{}) *NeoResponse
	// 17.6.7
	GetRelationshipProperty(relationship *NeoRelationship, propertyKey string) (interface{}, *NeoResponse)
	// 17.6.8
	SetRelationshipProperty(rel *NeoRelationship, propertyKey string, propertyValue interface{}) *NeoResponse

	// 17.6.9
	GetAllNodeRelationships(*NeoNode) ([]*NeoRelationship, *NeoResponse)
	// 17.6.10
	GetIncomingNodeRelationships(*NeoNode) ([]*NeoRelationship, *NeoResponse)
	// 17.6.11
	GetOoutgoingNodeRelationships(*NeoNode) ([]*NeoRelationship, *NeoResponse)

	// 17.6.12
	GetAllNodeRelationshipsWithTypes(node *NeoNode, relTypes []string) ([]*NeoRelationship, *NeoResponse)
	GetIncomingNodeRelationshipsWithTypes(node *NeoNode, relTypes []string) ([]*NeoRelationship, *NeoResponse)
	GetOoutgoingNodeRelationshipsWithTypes(node *NeoNode, relTypes []string) ([]*NeoRelationship, *NeoResponse)

	// 17.7.1
	GetRelationshipTypes() ([]string, *NeoResponse)

	// 17.8.1
	SetNodeProperty(node *NeoNode, properyyKey string, propertyValue interface{}) *NeoResponse
	// 17.8.2
	ReplaceNodeProperties(*NeoNode, map[string]interface{}) *NeoResponse
	// 17.8.3
	GetNodeProperties(*NeoNode) (map[string]interface{}, *NeoResponse)
	// 17.8.6
	DeleteNodeProperties(*NeoNode) *NeoResponse
	// 17.8.7
	DeleteNodePropertyForKeyName(node *NeoNode, keyName string) *NeoResponse

	// 17.9.1
	UpdateRelationshipProperty(*NeoRelationship, map[string]interface{}) *NeoResponse
	// 17.9.2
	DeleteRelationshipProperties(*NeoRelationship) *NeoResponse
	// 17.9.3
	DeleteRelationshipPropertyForKey(relationship *NeoRelationship, keyName string) *NeoResponse
}
