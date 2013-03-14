package neo2go

type GraphPathFinder interface {
	// 17.15.1+
	FindPathFromNode(start *NeoNode, target *NeoNode, spec *NeoPathFinderSpec) (*NeoPath, *NeoResponse)
	FindPathsFromNode(start *NeoNode, target *NeoNode, spec *NeoPathFinderSpec) ([]*NeoPath, *NeoResponse)
}
