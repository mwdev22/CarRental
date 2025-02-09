package types

type QueryOptions struct {
	Limit        int
	Offset       int
	SortField    string
	SortDiretion string
}

type FilterOperators map[string]string

var OperatorMap = FilterOperators{
	"eq":  "=",
	"neq": "!=",
	"gt":  ">",
	"lt":  "<",
	"gte": ">=",
	"lte": "<=",
	"ct":  "LIKE",
	"sw":  "LIKE",
	"ew":  "LIKE",
}

type QueryFilter struct {
	Field    string
	Operator string
	Value    interface{}
}
