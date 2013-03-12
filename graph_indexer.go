package neo2go

type GraphIndexer interface {
	// 17.10.1 - Nodes
	CreateNodeIndex(name string) (*NeoIndex, NeoResponse)
	// 17.10.2
	CreateNodeIndexWithConfiguration(name string, config map[string]interface{}) (*NeoIndex, NeoResponse)
	// 17.10.3
	DeleteIndex(*NeoIndex) NeoResponse
	// 17.10.4
	GetNodeIndexes() ([]*NeoIndex, NeoResponse)
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

	// 17.10.1 - Relationships
	CreateRelationshipIndex(name string) (*NeoIndex, NeoResponse)
	// 17.10.2
	CreateRelationshipIndexWithConfiguration(name string, config map[string]interface{}) (*NeoIndex, NeoResponse)
	// 17.10.4
	GetRelationshipIndexes() ([]*NeoIndex, NeoResponse)
	// 17.10.5
	AddRelationshipToIndex(*NeoRelationship, *NeoIndex) (*NeoRelationship, NeoResponse)
	// 17.10.6
	RemoveAllIndexEntriesForRelationship(*NeoRelationship, *NeoIndex) NeoResponse
	// 17.10.7
	RemoveAllIndexEntriesForRelationshipAndKey(rel *NeoRelationship, key string, index *NeoIndex) NeoResponse
	// 17.10.8
	RemoveAllIndexEntriesForRelationshipKeyAndValue(rel *NeoRelationship, key string, value string, index *NeoIndex) NeoResponse
	// 17.10.9
	FindRelationshipByExactMatch(value string) ([]*NeoRelationship, NeoResponse)
	// 17.10.10
	FindRelationshipByQuery(query string) ([]*NeoRelationship, NeoResponse)

	// 17.11.1
	GetOrCreateUniqueNode(*NeoIndex) (*NeoNode, *NeoResponse)
	GetOrCreateUniqueNodeWithProperties(*NeoIndex, map[string]interface{}) (*NeoNode, *NeoResponse)

	// 17.11.3
	CreateUniqueNodeOrFail() (*NeoIndex, *NeoNode, *NeoResponse)
	CreateUniqueNodeWithPropertiesOrFail(*NeoIndex, map[string]interface{}) (*NeoNode, *NeoResponse)

	// 17.11.5
	GetOrCreateUniqueRelationship(index *NeoIndex, source *NeoNode, target *NeoNode) (*NeoRelationship, *NeoResponse)
	GetOrCreateUniqueTypedRelationship(index *NeoIndex, source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse)
	GetOrCreateUniqueRelationshipWithProperties(index *NeoIndex, source *NeoNode, target *NeoNode, properties map[string]interface{}) (*NeoRelationship, *NeoResponse)
	GetOrCreateUniqueTypedRelationshipWithProperties(index *NeoIndex, source *NeoNode, target *NeoNode, relType string, properties map[string]interface{}) (*NeoRelationship, *NeoResponse)

	// 17.11.7
	CreateUniqueRelationshipOrFail(index *NeoIndex, source *NeoNode, target *NeoNode) (*NeoRelationship, *NeoResponse)
	CreateUniqueTypedRelationshipOrFail(index *NeoIndex, source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse)
	CreateUniqueRelationshipWithPropertiesOrFail(index *NeoIndex, source *NeoNode, target *NeoNode, properties map[string]interface{}) (*NeoRelationship, *NeoResponse)
	CreateUniqueTypedRelationshipWithPropertiesOrFail(index *NeoIndex, source *NeoNode, target *NeoNode, relType string, properties map[string]interface{}) (*NeoRelationship, *NeoResponse)

	// 17.12.1
	FindNodeByExactMatchingAutoIndex(query string) ([]*NeoNode, NeoResponse)
	// 17.12.2
	FindNodeByQueryingAutoIndex(query string) ([]*NeoNode, NeoResponse)

	// 17.13.1 - autoindex configuration for Nodes
	CreateNodeAutoIndexWithConfiguration(name string, config map[string]interface{}) (*NeoIndex, NeoResponse)
	// 17.13.3
	GetNodeAutoIndexStatus() (bool, *NeoResponse)
	// 17.13.4
	SetNodeAutoIndexStatus(bool) *NeoResponse
	// 17.13.5
	GetNodeAutoIndexProperties() ([]string, *NeoResponse)
	// 17.13.6
	AddNodeAutoIndexProperty(propertyName string) *NeoResponse
	// 17.13.7
	DeleteNodeAutoIndexProperty(propertyName string) *NeoResponse

	// 17.13.2 - autoindex configuration for Relationships
	CreateRelationshipAutoIndexWithConfiguration(name string, config map[string]interface{}) (*NeoIndex, NeoResponse)
	// 17.13.3
	GetRelationshipAutoIndexStatus() (bool, *NeoResponse)
	// 17.13.4
	SetRelationshipAutoIndexStatus(bool) *NeoResponse
	// 17.13.5
	GetRelationshipAutoIndexProperties() ([]string, *NeoResponse)
	// 17.13.6
	AddRelationshipAutoIndexProperty(propertyName string) *NeoResponse
	// 17.13.7
	DeleteRelationshipAutoIndexProperty(propertyName string) *NeoResponse
}
