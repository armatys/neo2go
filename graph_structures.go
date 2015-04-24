package neo2go

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
)

// http://docs.neo4j.org/chunked/milestone/graphdb-neo4j-properties.html
/*
Type    Description       Value range

boolean                   true/false

byte    8-bit integer     -128 to 127, inclusive

short   16-bit integer    -32768 to 32767, inclusive

int     32-bit integer    -2147483648 to 2147483647, inclusive

long    64-bit integer    -9223372036854775808 to 9223372036854775807, inclusive

float   32-bit IEEE 754 floating-point number

double  64-bit IEEE 754 floating-point number

char    16-bit unsigned   u0000 to uffff (0 to 65535)
        integers representing
        Unicode characters

String  sequence of Unicode characters
*/

func setTemplateIfNil(tmpl **UrlTemplate, val string) {
	if *tmpl == nil {
		newTempl := NewUrlTemplate(val)
		*tmpl = newTempl
	}
}

type NeoNode struct {
	AllRelationships           *UrlTemplate           `json:"all_relationships"`
	AllTypedRelationships      *UrlTemplate           `json:"all_typed_relationships"`
	CreateRelationship         *UrlTemplate           `json:"create_relationship"`
	Data                       json.RawMessage        `json:"data"`
	Extensions                 map[string]interface{} `json:"extensions"`
	IncomingRelationships      *UrlTemplate           `json:"incoming_relationships"`
	IncomingTypedRelationships *UrlTemplate           `json:"incoming_typed_relationships"`
	Labels                     *UrlTemplate           `json:"labels"`
	OutgoingRelationships      *UrlTemplate           `json:"outgoing_relationships"`
	OutgoingTypedRelationships *UrlTemplate           `json:"outgoing_typed_relationships"`
	PagedTraverse              *UrlTemplate           `json:"paged_traverse"`
	Properties                 *UrlTemplate           `json:"properties"`
	Property                   *UrlTemplate           `json:"property"`
	Self                       *UrlTemplate           `json:"self"`
	Traverse                   *UrlTemplate           `json:"traverse"`
	batchId                    NeoBatchId
}

func (n *NeoNode) ParseData(result interface{}) error {
	return json.Unmarshal([]byte(n.Data), result)
}

func (n *NeoNode) Id() int64 {
	if n.Self != nil {
		selfUri := n.Self.String()
		_, file := path.Split(selfUri)
		id, err := strconv.ParseInt(file, 10, 64)
		if err != nil {
			return 0
		}
		return id
	}
	return 0
}

func (n *NeoNode) IdOrBatchId() string {
	if n.Self != nil {
		return fmt.Sprintf("%d", n.Id())
	} else if n.batchId > 0 {
		return fmt.Sprintf("{%d}", n.batchId)
	}
	return ""
}

func (n *NeoNode) setBatchId(bid NeoBatchId) {
	n.batchId = bid

	setTemplateIfNil(&n.AllRelationships, fmt.Sprintf(`{%v}/relationships/all`, bid))
	setTemplateIfNil(&n.AllTypedRelationships, fmt.Sprintf(`{%v}/relationships/all/{-list|&|types}`, bid))
	setTemplateIfNil(&n.CreateRelationship, fmt.Sprintf(`{%v}/relationships`, bid))
	setTemplateIfNil(&n.IncomingRelationships, fmt.Sprintf(`{%v}/relationships/in`, bid))
	setTemplateIfNil(&n.IncomingTypedRelationships, fmt.Sprintf(`{%v}/relationships/in/{-list|&|types}`, bid))
	setTemplateIfNil(&n.OutgoingRelationships, fmt.Sprintf(`{%v}/relationships/out`, bid))
	setTemplateIfNil(&n.OutgoingTypedRelationships, fmt.Sprintf(`{%v}/relationships/out/{-list|&|types}`, bid))
	setTemplateIfNil(&n.PagedTraverse, fmt.Sprintf(`{%v}/paged/traverse/{returnType}{?pageSize,leaseTime}`, bid))
	setTemplateIfNil(&n.Properties, fmt.Sprintf(`{%v}/properties`, bid))
	setTemplateIfNil(&n.Property, fmt.Sprintf(`{%v}/properties/{key}`, bid))
	setTemplateIfNil(&n.Self, fmt.Sprintf(`{%v}`, bid))
	setTemplateIfNil(&n.Traverse, fmt.Sprintf(`{%v}/traverse/{returnType}`, bid))
}

func (n *NeoNode) String() string {
	return fmt.Sprintf("<Node id:%d>", n.Id())
}

type NeoRelationship struct {
	Data       map[string]interface{} `json:"data"`
	Extensions map[string]interface{} `json:"extensions"`
	Start      *UrlTemplate           `json:"start"`
	Property   *UrlTemplate           `json:"property"`
	Self       *UrlTemplate           `json:"self"`
	Properties *UrlTemplate           `json:"properties"`
	Type       string                 `json:"type"`
	End        *UrlTemplate           `json:"end"`
	batchId    NeoBatchId
}

func (n *NeoRelationship) ParseData(result interface{}) error {
	bytes, err := json.Marshal(n.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, result)
}

func (n *NeoRelationship) Id() int64 {
	if n.Self != nil {
		selfUri := n.Self.String()
		_, file := path.Split(selfUri)
		id, err := strconv.ParseInt(file, 10, 64)
		if err != nil {
			return 0
		}
		return id
	}
	return 0
}

func (n *NeoRelationship) IdOrBatchId() string {
	if n.Self != nil {
		return fmt.Sprintf("%d", n.Id())
	} else if n.batchId > 0 {
		return fmt.Sprintf("{%d}", n.batchId)
	}
	return ""
}

func (n *NeoRelationship) setBatchId(bid NeoBatchId) {
	n.batchId = bid

	setTemplateIfNil(&n.Property, fmt.Sprintf(`{%v}/properties/{key}`, bid))
	setTemplateIfNil(&n.Self, fmt.Sprintf(`{%v}`, bid))
	setTemplateIfNil(&n.Properties, fmt.Sprintf(`{%v}/properties`, bid))
}

type NeoIndex struct {
	Provider string
	Template *UrlTemplate
	Type     string
	batchId  NeoBatchId
}

func (n *NeoIndex) setBatchId(bid NeoBatchId) {
	n.batchId = bid

	setTemplateIfNil(&n.Template, fmt.Sprintf(`{%v}{key}/{value}`, bid))
}

type NeoCodeSnippet struct {
	Body     string `json:"body,omitempty"`
	Language string `json:"language"`
	Name     string `json:"name,omitempty"`
}

type NeoTraversalOrder uint8

const (
	NeoTraversalBreadthFirst NeoTraversalOrder = iota
	NeoTraversalDepthFirst
)

func (n NeoTraversalOrder) MarshalJSON() ([]byte, error) {
	if n == NeoTraversalBreadthFirst {
		return []byte(`"breadth_first"`), nil
	} else if n == NeoTraversalDepthFirst {
		return []byte(`"depth_first"`), nil
	}
	return nil, fmt.Errorf("Could not marshal NeoTraversalOrder to JSON.")
}

type NeoTraversalDirection uint8

const (
	NeoTraversalAll NeoTraversalDirection = iota
	NeoTraversalIn
	NeoTraversalOut
)

func (n NeoTraversalDirection) MarshalJSON() ([]byte, error) {
	if n == NeoTraversalAll {
		return []byte(`"all"`), nil
	} else if n == NeoTraversalIn {
		return []byte(`"in"`), nil
	} else if n == NeoTraversalOut {
		return []byte(`"out"`), nil
	}
	return nil, fmt.Errorf("Could not marshal NeoTraversalDirection to JSON.")
}

type NeoTraversalUniqueness uint8

const (
	NeoTraversalNodeGlobal NeoTraversalUniqueness = iota
	NeoTraversalNone
	NeoTraversalRelationhipGlobal
	NeoTraversalNodePath
	NeoTraversalRelationshipPath
)

func (n NeoTraversalUniqueness) MarshalJSON() ([]byte, error) {
	if n == NeoTraversalNodeGlobal {
		return []byte(`"node_global"`), nil
	} else if n == NeoTraversalNone {
		return []byte(`"none"`), nil
	} else if n == NeoTraversalRelationhipGlobal {
		return []byte(`"relationship_global"`), nil
	} else if n == NeoTraversalNodePath {
		return []byte(`"node_path"`), nil
	} else if n == NeoTraversalRelationshipPath {
		return []byte(`"relationship_path"`), nil
	}
	return nil, fmt.Errorf("Could not marshal NeoTraversalUniqueness to JSON.")
}

type NeoTraversalRelationship struct {
	Direction NeoTraversalDirection `json:"direction"`
	Type      string                `json:"type"`
}

func NewNeoReturnFilterAll() *NeoCodeSnippet {
	return &NeoCodeSnippet{Name: "all", Language: "builtin"}
}

func NewNeoReturnFilterAllButStartNode() *NeoCodeSnippet {
	return &NeoCodeSnippet{Name: "all_but_start_node", Language: "builtin"}
}

type NeoTraversal struct {
	LeaseTime      uint32                      `json:"-"`
	PageSize       uint32                      `json:"-"`
	Order          NeoTraversalOrder           `json:"order"`
	Relationships  []*NeoTraversalRelationship `json:"relationships,omitempty"`
	Uniqueness     NeoTraversalUniqueness      `json:"uniqueness"`
	PruneEvaluator *NeoCodeSnippet             `json:"prune_evaluator,omitempty"`
	ReturnFilter   *NeoCodeSnippet             `json:"return_filter,omitempty"`
	MaxDepth       uint32                      `json:"max_depth"`
}

type NeoPath struct {
	Weight        float64
	Start         string
	Nodes         []string
	Length        uint32
	Relationships []string
	End           string
}

type NeoFullPath struct {
	Start         *NeoNode
	Nodes         []*NeoNode
	Length        uint32
	Relationships []*NeoRelationship
	End           *NeoNode
}

type NeoPagedTraverser struct {
	location string
}

type NeoGraphAlgorithm uint8

const (
	NeoShortestPath NeoGraphAlgorithm = iota
	NeoAllSimplePaths
	NeoAllPaths
	NeoDijkstra
)

func (n NeoGraphAlgorithm) MarshalJSON() ([]byte, error) {
	if n == NeoShortestPath {
		return []byte(`"shortestPath"`), nil
	} else if n == NeoAllSimplePaths {
		return []byte(`"allSimplePaths"`), nil
	} else if n == NeoAllPaths {
		return []byte(`"allPaths"`), nil
	} else if n == NeoDijkstra {
		return []byte(`"dijkstra"`), nil
	}
	return nil, fmt.Errorf("Could not marshal NeoGraphAlgorithm to JSON.")
}

type NeoPathFinderSpec struct {
	CostProperty  string                    `json:"cost_property,omitempty"`
	DefaultCost   float64                   `json:"default_cost,omitempty"`
	MaxDepth      uint32                    `json:"max_depth"`
	Relationships *NeoTraversalRelationship `json:"relationships,omitempty"`
	Algorithm     NeoGraphAlgorithm         `json:"algorithm"`
	To            string                    `json:"to"`
}

func NewNeoPathFinderSpecWithRelationships(rels *NeoTraversalRelationship) *NeoPathFinderSpec {
	spec := new(NeoPathFinderSpec)
	spec.MaxDepth = 1
	spec.Algorithm = NeoShortestPath
	spec.Relationships = rels
	return spec
}
