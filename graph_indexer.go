package neo2go

type NeoIndex struct{}

type GraphIndexer interface {
	// 17.10.1
	CreateNodeIndex(name string) (*NeoIndex, NeoResponse)
	// 17.10.2
	CreateNodeIndexWithConfiguration(name string, config map[string]interface{}) (*NeoIndex, NeoResponse)
	// 17.10.3
	DeleteIndex(*NeoIndex) NeoResponse
	// 17.10.4
	GetIndexList() ([]*NeoIndex, NeoResponse)
	// 17.10.5
	AddNodeToIndex(*NeoNode, *NeoIndex) (*NeoNode, NeoResponse)
	// 17.10.6
	RemoveAllIndexEntriesForNode(*NeoNode, *NeoIndex) NeoResponse
	// 17.10.7
	RemoveAllIndexEntriesForNodeAndKey(node *NeoNode, key string, index *NeoIndex) NeoResponse
	// 17.10.8
	RemoveAllIndexEntriesForNodeKeyAndValue(node *NeoNode, key string, value string, index *NeoIndex) NeoResponse
	// 17.10.9
	FindNodeByExactMatch(value string) ([]*NeoNode, NeoResponse)
	// 17.10.10
	FindNodeByQuery(query string) ([]*NeoNode, NeoResponse)

	// 17.11.1
	GetOrCreateUniqueNode() (*NeoNode, *NeoResponse)
	GetOrCreateUniqueNodeWithProperties(map[string]interface{}) (*NeoNode, *NeoResponse)

	// 17.11.3
	CreateUniqueNodeOrFail() (*NeoNode, *NeoResponse)
	CreateUniqueNodeWithPropertiesOrFail(map[string]interface{}) (*NeoNode, *NeoResponse)

	// 17.11.5
	GetOrCreateUniqueRelationship(source *NeoNode, target *NeoNode) (*NeoRelationship, *NeoResponse)
	GetOrCreateUniqueRelationshipWithType(source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse)
	GetOrCreateUniqueRelationshipWithProperties(source *NeoNode, target *NeoNode, properties map[string]interface{}) (*NeoRelationship, *NeoResponse)
	GetOrCreateUniqueRelationshipWithPropertiesAndType(source *NeoNode, target *NeoNode, properties map[string]interface{}, relType string) (*NeoRelationship, *NeoResponse)

	// 17.11.7
	CreateUniqueRelationshipOrFail(source *NeoNode, target *NeoNode) (*NeoRelationship, *NeoResponse)
	CreateUniqueRelationshipWithTypeOrFail(source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse)
	CreateUniqueRelationshipWithPropertiesOrFail(source *NeoNode, target *NeoNode, properties map[string]interface{}) (*NeoRelationship, *NeoResponse)
	CreateUniqueRelationshipWithPropertiesAndTypeOrFail(source *NeoNode, target *NeoNode, properties map[string]interface{}, relType string) (*NeoRelationship, *NeoResponse)

	// 17.12.1
	FindNodeByExactMatchFromAutoIndex(query string) ([]*NeoNode, NeoResponse)
	// 17.12.2
	FindNodeByQueryFromAutoIndex(query string) ([]*NeoNode, NeoResponse)
}
