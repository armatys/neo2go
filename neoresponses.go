package neo2go

type NeoError struct {
	Message    string `json:"message"`
	Exception  string `json:"exception"`
	Stacktrace string `json:"stacktrace"`
}

func (n *NeoError) Error() string {
	return n.Message
}

type NeoErrors struct {
	Errors []NeoError `json:"errors"`
}

func (n *NeoErrors) Error() string {
	s := ""
	for _, n := range n.Errors {
		s = s + " " + n.Error()
	}
	return s
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

type NeoRoot struct {
	Management *UrlTemplate `json:"management"`
	Data       *UrlTemplate `json:"data"`
}

type NeoDataRoot struct {
	Extensions        *UrlTemplate `json:"extensions"`
	ExtensionsInfo    *UrlTemplate `json:"extensions_info"`
	Node              *UrlTemplate `json:"node"`
	NodeIndex         *UrlTemplate `json:"node_index"`
	NodeLabels        *UrlTemplate `json:"node_labels"`
	RelationshipIndex *UrlTemplate `json:"relationship_index"`
	RelationshipTypes *UrlTemplate `json:"relationship_types"`
	Batch             *UrlTemplate `json:"batch"`
	Cypher            *UrlTemplate `json:"cypher"`
	Indexes           *UrlTemplate `json:"indexes"`
	Constraints       *UrlTemplate `json:"constraints"`
	Transaction       *UrlTemplate `json:"transaction"`
	Neo4jVersion      string       `json:"neo4j_version"`
}

type CypherResponse struct {
	Columns []string
	Data    [][]interface{}
}
