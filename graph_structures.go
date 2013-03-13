package neo2go

import (
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
	AllRelationships           *UrlTemplate
	AllTypedRelationships      *UrlTemplate
	CreateRelationship         *UrlTemplate
	Data                       map[string]interface{}
	Extensions                 map[string]interface{}
	IncomingRelationships      *UrlTemplate
	IncomingTypedRelationships *UrlTemplate
	OutgoingRelationships      *UrlTemplate
	OutgoingTypedRelationships *UrlTemplate
	PagedTraverse              *UrlTemplate
	Properties                 *UrlTemplate
	Property                   *UrlTemplate
	Self                       *UrlTemplate
	Traverse                   *UrlTemplate
	batchId                    NeoBatchId
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

type NeoRelationship struct {
	Data       map[string]interface{}
	Extensions []*UrlTemplate
	Start      *UrlTemplate
	Property   *UrlTemplate
	Self       *UrlTemplate
	Properties *UrlTemplate
	Type       string
	End        *UrlTemplate
	batchId    NeoBatchId
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

	setTemplateIfNil(&n.Template, fmt.Sprintf(`{%v}`, bid))
}

type NeoCodeSnippet struct {
	Body     string
	Language string
}

type NeoTraversalOrder uint8

const (
	NeoTraversalBreadthFirst NeoTraversalOrder = iota
	NeoTraversalDepthFirst
)

type NeoTraversalDirection uint8

const (
	NeoTraversalAll NeoTraversalDirection = iota
	NeoTraversalIn
	NeoTraversalOut
)

type NeoTraversalUniqueness uint8

const (
	NeoTraversalNodeGlobal NeoTraversalUniqueness = iota
	NeoTraversalNone
	NeoTraversalRelationhipGlobal
	NeoTraversalNodePath
	NeoTraversalRelationshipPath
)

type NeoTraversalRelationship struct {
	Direction NeoTraversalDirection
	Type      string
}

type NeoTraversal struct {
	LeaseTime      uint32
	PageSize       uint32
	Order          NeoTraversalOrder
	Relationships  []*NeoTraversalRelationship
	Uniqueness     NeoTraversalUniqueness
	PruneEvaluator *NeoCodeSnippet
	ReturnFilter   *NeoCodeSnippet
	MaxDepth       uint32
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
	LeaseTime uint32
	PageSize  uint32
	url       string
}

type NeoGraphAlgorithm uint8

const (
	NeoShortestPath NeoGraphAlgorithm = iota
	NeoAllSimplePaths
	NeoAllPaths
	NeoDijkstra
)

type NeoPathFinderSpec struct {
	CostProperty  string
	DefaultCost   float64
	MaxDepth      uint32
	Relationships *NeoTraversalRelationship
	Algorithm     NeoGraphAlgorithm
}
