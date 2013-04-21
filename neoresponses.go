package neo2go

type NeoError struct {
	Message    string   `json:"message"`
	Exception  string   `json:"exception"`
	Stacktrace []string `json:"stacktrace"`
}

func (n *NeoError) Error() string {
	return n.Message
}

type NeoResponse struct {
	ExpectedCode int
	StatusCode   int
	NeoError     error
	location     string
}

func NewLocalErrorResponse(expectedCode int, err error) *NeoResponse {
	return &NeoResponse{expectedCode, 600, err, ""}
}

func (n *NeoResponse) Ok() bool {
	if n.ExpectedCode == 200 {
		return n.StatusCode >= 200 && n.StatusCode < 300
	}
	return n.ExpectedCode == n.StatusCode
}

func (n *NeoResponse) Created() bool {
	return n.StatusCode == 201
}

type NeoServiceRoot struct {
	Node              *UrlTemplate `json:"node"`
	ReferenceNode     *UrlTemplate `json:"reference_node"`
	NodeIndex         *UrlTemplate `json:"node_index"`
	RelationshipIndex *UrlTemplate `json:"relationship_index"`
	ExtensionsInfo    *UrlTemplate `json:"extensions_info"`
	RelationshipTypes *UrlTemplate `json:"relationship_types"`
	Batch             *UrlTemplate `json:"batch"`
	Cypher            *UrlTemplate `json:"cypher"`
	Neo4jVersion      string       `json:"neo4j_version"`
}

type CypherResponse struct {
	Columns []string
	Data    [][]interface{}
}
