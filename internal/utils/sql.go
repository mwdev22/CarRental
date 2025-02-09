package utils

import (
	"fmt"

	"github.com/mwdev22/CarRental/internal/types"
)

func BuildBatchQuery(query string, filters []*types.QueryFilter, opts *types.QueryOptions) (string, []interface{}) {
	args := make([]interface{}, 0)
	// i use question marks because filters and their count are dynamic
	// could be 1, 2, 3 etc through loop, but its not necessary complexity i think
	for _, filter := range filters {
		query += fmt.Sprintf(" AND %s %s ?", filter.Field, filter.Operator)
		args = append(args, filter.Value)
	}

	if opts != nil {
		query += fmt.Sprintf(` ORDER BY %s %s`, opts.SortField, opts.SortDiretion)
		query += ` LIMIT ? OFFSET ?`
		args = append(args, opts.Limit, opts.Offset)
	} else {
		query += ` ORDER BY id DESC`
	}

	return query, args
}
