package queryserver

import (
	"github.com/pganalyze/pg_query_go/v6/parser"
)

type Query struct {
	input   []byte
	Options QueryOptions
}

type QueryOptions struct {
	Deparse     bool
	SpiltColumn bool
	Normalize   bool
	Fingerprint bool
}

type Result struct {
	input []byte
	sql   SQL
}

type SQL struct {
	ParsedData   []byte
	NormalizeSQL []byte
	Fingerprint  string
}

func ParseSQL(sql *Query) (*Result, error) {
	// pg_query_go.Parse()
	protobufTree, err := parser.ParseToProtobuf(string(sql.input))
	if err != nil {
		return nil, err
	}
}
