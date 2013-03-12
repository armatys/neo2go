package neo2go

type GraphTraverser interface {
	// 17.14.1+
	TraverseNodes(*NeoTraversal) ([]*NeoNode, *NeoResponse)
	TraverseRelationships(*NeoTraversal) ([]*NeoRelationship, *NeoResponse)
	TraversePaths(*NeoTraversal) ([]*NeoPath, *NeoResponse)
	TraverseFullPaths(*NeoTraversal) ([]*NeoFullPath, *NeoResponse)

	// 17.14.5
	TraverseNodesWithPaging(*NeoTraversal) (*NeoPagedTraverser, []*NeoNode, *NeoResponse)
	TraverseRelationshipsWithPaging(*NeoTraversal) (*NeoPagedTraverser, []*NeoRelationship, *NeoResponse)
	TraversePathsWithPaging(*NeoTraversal) (*NeoPagedTraverser, []*NeoPath, *NeoResponse)
	TraverseFullPathsWithPaging(*NeoTraversal) (*NeoPagedTraverser, []*NeoFullPath, *NeoResponse)

	// 17.14.6+
	TraverseNodesGetNextPage(*NeoPagedTraverser) ([]*NeoNode, *NeoResponse)
	TraverseRelationshipsGetNextPage(*NeoPagedTraverser) ([]*NeoRelationship, *NeoResponse)
	TraversePathsGetNextPage(*NeoPagedTraverser) ([]*NeoPath, *NeoResponse)
	TraverseFullPathsGetNextPage(*NeoPagedTraverser) ([]*NeoFullPath, *NeoResponse)
}
