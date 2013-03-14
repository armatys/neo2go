package neo2go

type GraphTraverser interface {
	// 17.14.1+
	TraverseByNodes(traversal *NeoTraversal, start *NeoNode) ([]*NeoNode, *NeoResponse)
	TraverseByRelationships(traversal *NeoTraversal, start *NeoNode) ([]*NeoRelationship, *NeoResponse)
	TraverseByPaths(traversal *NeoTraversal, start *NeoNode) ([]*NeoPath, *NeoResponse)
	TraverseByFullPaths(traversal *NeoTraversal, start *NeoNode) ([]*NeoFullPath, *NeoResponse)

	// 17.14.5
	TraverseByNodesWithPaging(traversal *NeoTraversal, start *NeoNode) (*NeoPagedTraverser, []*NeoNode, *NeoResponse)
	TraverseByRelationshipsWithPaging(traversal *NeoTraversal, start *NeoNode) (*NeoPagedTraverser, []*NeoRelationship, *NeoResponse)
	TraverseByPathsWithPaging(traversal *NeoTraversal, start *NeoNode) (*NeoPagedTraverser, []*NeoPath, *NeoResponse)
	TraverseByFullPathsWithPaging(traversal *NeoTraversal, start *NeoNode) (*NeoPagedTraverser, []*NeoFullPath, *NeoResponse)

	// 17.14.6+
	TraverseByNodesGetNextPage(*NeoPagedTraverser) ([]*NeoNode, *NeoResponse)
	TraverseByRelationshipsGetNextPage(*NeoPagedTraverser) ([]*NeoRelationship, *NeoResponse)
	TraverseByPathsGetNextPage(*NeoPagedTraverser) ([]*NeoPath, *NeoResponse)
	TraverseByFullPathsGetNextPage(*NeoPagedTraverser) ([]*NeoFullPath, *NeoResponse)
}
