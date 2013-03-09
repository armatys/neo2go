package neo2go

type GraphIndexer interface {
	/*
		CreateNodeIndex(name string) (*NeoIndex, NeoResponse)
		CreateNodeIndexWithConfiguration(name string, config map[string]interface{}) (*NeoIndex, NeoResponse)
		DeleteIndex(*NeoIndex) NeoResponse
		GetIndexList() ([]*NeoIndex, NeoResponse)
		AddNodeToIndex(*NeoNode, *NeoIndex) (*NeoNode, NeoResponse)
		RemoveAllIndexEntriesForNode(*NeoNode, *NeoIndex) NeoResponse
		RemoveAllIndexEntriesForNodeAndKey(node *NeoNode, key string, index *NeoIndex) NeoResponse
		RemoveAllIndexEntriesForNodeKeyAndValue(node *NeoNode, key string, value string, index *NeoIndex) NeoResponse
		FindNodeByExactMatch(value string) ([]*NeoNode, NeoResponse)
		FindNodeByQuery(query string) ([]*NeoNode, NeoResponse)
	*/
}
