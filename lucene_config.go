package neo2go

import (
	"encoding/json"
)

type LuceneAnalyzerType string

const (
	LuceneAnalyzerExact    LuceneAnalyzerType = "exact"
	LuceneAnalyzerFullText LuceneAnalyzerType = "fulltext"
)

type LuceneIndexConfig struct {
	provider    string             `json:"provider"`
	Type        LuceneAnalyzerType `json:"type"`
	ToLowerCase bool               `json:"to_lower_case"`
	Analyzer    string             `json:"analyzer,omitempty"`
}

func NewLuceneIndexConfig() *LuceneIndexConfig {
	return &LuceneIndexConfig{"", LuceneAnalyzerExact, true, ""}
}

func (l *LuceneIndexConfig) MarshalJSON() ([]byte, error) {
	l.provider = "lucene"
	return json.Marshal(l)
}
