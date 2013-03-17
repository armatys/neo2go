package neo2go

type GraphIndexer interface {
	// 17.10.1 - Nodes
	CreateNodeIndex(name string) (*NeoIndex, *NeoResponse)
	// 17.10.2
	CreateNodeIndexWithConfiguration(name string, config interface{}) (*NeoIndex, *NeoResponse)
	// 17.10.3
	DeleteIndex(*NeoIndex) *NeoResponse
	// 17.10.4
	GetNodeIndexes() (map[string]*NeoIndex, *NeoResponse)
	// 17.10.5
	AddNodeToIndex(index *NeoIndex, node *NeoNode, key, value string) (*NeoNode, *NeoResponse)
	// 17.10.6
	DeleteAllIndexEntriesForNode(*NeoIndex, *NeoNode) *NeoResponse
	// 17.10.7
	DeleteAllIndexEntriesForNodeAndKey(index *NeoIndex, node *NeoNode, key string) *NeoResponse
	// 17.10.8
	DeleteAllIndexEntriesForNodeKeyAndValue(index *NeoIndex, node *NeoNode, key string, value string) *NeoResponse
	// 17.10.9
	FindNodeByExactMatch(index *NeoIndex, key, value string) ([]*NeoNode, *NeoResponse)
	// 17.10.10
	FindNodeByQuery(index *NeoIndex, query string) ([]*NeoNode, *NeoResponse)

	// 17.10.1 - Relationships
	CreateRelationshipIndex(name string) (*NeoIndex, *NeoResponse)
	// 17.10.2
	CreateRelationshipIndexWithConfiguration(name string, config interface{}) (*NeoIndex, *NeoResponse)
	// 17.10.4
	GetRelationshipIndexes() ([]*NeoIndex, *NeoResponse)
	// 17.10.5
	AddRelationshipToIndex(index *NeoIndex, rel *NeoRelationship, key, value string) (*NeoRelationship, *NeoResponse)
	// 17.10.6
	DeleteAllIndexEntriesForRelationship(*NeoIndex, *NeoRelationship) *NeoResponse
	// 17.10.7
	DeleteAllIndexEntriesForRelationshipAndKey(index *NeoIndex, rel *NeoRelationship, key string) *NeoResponse
	// 17.10.8
	DeleteAllIndexEntriesForRelationshipKeyAndValue(index *NeoIndex, rel *NeoRelationship, key string, value string) *NeoResponse
	// 17.10.9
	FindRelationshipByExactMatch(index *NeoIndex, key, value string) ([]*NeoRelationship, *NeoResponse)
	// 17.10.10
	FindRelationshipByQuery(index *NeoIndex, query string) ([]*NeoRelationship, *NeoResponse)

	// 17.11.1
	GetOrCreateUniqueNode(index *NeoIndex, key, value string) (*NeoNode, *NeoResponse)
	GetOrCreateUniqueNodeWithProperties(index *NeoIndex, key, value string, properties interface{}) (*NeoNode, *NeoResponse)

	// 17.11.3
	CreateUniqueNodeOrFail(index *NeoIndex, key, value string) (*NeoNode, *NeoResponse)
	CreateUniqueNodeWithPropertiesOrFail(index *NeoIndex, key, value string, properties interface{}) (*NeoNode, *NeoResponse)

	// 17.11.5
	GetOrCreateUniqueRelationship(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode) (*NeoRelationship, *NeoResponse)
	GetOrCreateUniqueTypedRelationship(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse)
	GetOrCreateUniqueRelationshipWithProperties(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, properties interface{}) (*NeoRelationship, *NeoResponse)
	GetOrCreateUniqueTypedRelationshipWithProperties(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string, properties interface{}) (*NeoRelationship, *NeoResponse)

	// 17.11.7
	CreateUniqueRelationshipOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode) (*NeoRelationship, *NeoResponse)
	CreateUniqueTypedRelationshipOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string) (*NeoRelationship, *NeoResponse)
	CreateUniqueRelationshipWithPropertiesOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, properties interface{}) (*NeoRelationship, *NeoResponse)
	CreateUniqueTypedRelationshipWithPropertiesOrFail(index *NeoIndex, key, value string, source *NeoNode, target *NeoNode, relType string, properties interface{}) (*NeoRelationship, *NeoResponse)

	// 17.12.1
	// FindNodeByExactMatchingAutoIndex(query string) ([]*NeoNode, *NeoResponse)
	// 17.12.2
	// FindNodeByQueryingAutoIndex(query string) ([]*NeoNode, *NeoResponse)

	// 17.13.1 - autoindex configuration for Nodes
	// CreateNodeAutoIndexWithConfiguration(name string, config interface{}) (*NeoIndex, *NeoResponse)
	// 17.13.3
	// GetNodeAutoIndexStatus() (bool, *NeoResponse)
	// 17.13.4
	// SetNodeAutoIndexStatus(bool) *NeoResponse
	// 17.13.5
	// GetNodeAutoIndexProperties() ([]string, *NeoResponse)
	// 17.13.6
	// AddNodeAutoIndexProperty(propertyName string) *NeoResponse
	// 17.13.7
	// DeleteNodeAutoIndexProperty(propertyName string) *NeoResponse

	// 17.13.2 - autoindex configuration for Relationships
	// CreateRelationshipAutoIndexWithConfiguration(name string, config interface{}) (*NeoIndex, *NeoResponse)
	// 17.13.3
	// GetRelationshipAutoIndexStatus() (bool, *NeoResponse)
	// 17.13.4
	// SetRelationshipAutoIndexStatus(bool) *NeoResponse
	// 17.13.5
	// GetRelationshipAutoIndexProperties() ([]string, *NeoResponse)
	// 17.13.6
	// AddRelationshipAutoIndexProperty(propertyName string) *NeoResponse
	// 17.13.7
	// DeleteRelationshipAutoIndexProperty(propertyName string) *NeoResponse
}
