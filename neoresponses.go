package neo2go

type NeoError struct {
	Message    string   `json:"message"`
	Exception  string   `json:"exception"`
	Stacktrace []string `json:"stacktrace"`
}

type NeoResponse struct {
	StatusCode int
	NeoError   *NeoError
}

type NeoServiceRoot struct {
	Node              UrlTemplate `json:"node"`
	ReferenceNode     UrlTemplate `json:"reference_node,omitempty"`
	NodeIndex         UrlTemplate `json:"node_index"`
	RelationshipIndex UrlTemplate `json:"relationship_index"`
	ExtensionsInfo    UrlTemplate `json:"extensions_info"`
	RelationshipTypes UrlTemplate `json:"relationship_types"`
	Batch             UrlTemplate `json:"batch"`
	Cypher            UrlTemplate `json:"cypher"`
	Neo4jVersion      string      `json:"neo4j_version"`
}

type CypherResponse struct {
	Columns []string
	Data    []interface{}
}
