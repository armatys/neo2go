package neo2go

import (
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
	selfUri := n.Self.String()
	_, file := path.Split(selfUri)
	id, err := strconv.ParseInt(file, 10, 64)
	if err != nil {
		return 0
	}
	return id
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
