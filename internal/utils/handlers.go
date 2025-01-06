package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mwdev22/CarRental/internal/types"
)

// scraping filters from query
func ParseQueryFilters(r *http.Request) ([]*types.QueryFilter, error) {
	query := r.URL.Query()
	filters := make([]*types.QueryFilter, 0)

	for key, values := range query {
		// skip pagination and sorting
		if key == "page" || key == "page_size" || key == "sort" {
			continue
		}

		// loop over the filters and parse them for sql query
		for _, v := range values {
			valParts := strings.Split(v, "_")
			var operatorKey string
			value := valParts[0]
			// if not operator provided, default is eq ---> =
			if len(valParts) != 2 {
				operatorKey = "eq"
			} else {
				operatorKey = valParts[1]
			}

			operator, ok := types.OperatorMap[operatorKey]
			if !ok {
				return nil, types.BadRequest(fmt.Sprintf("invalid operator: %s", operatorKey))
			}
			// if operator is LIKE, we need to add % to the value, for proper sql query
			if operator == "LIKE" {
				switch operatorKey {
				case "sw":
					value = value + "%"
				case "ew":
					value = "%" + value
				case "ct":
					value = "%" + value + "%"
				}
			}
			filters = append(filters, &types.QueryFilter{
				Field:    key,
				Operator: operator,
				Value:    value,
			})
		}
	}

	return filters, nil
}

func ParseQueryOptions(r *http.Request) (*types.QueryOptions, error) {
	query := r.URL.Query()

	// default options if not provided in query
	opts := &types.QueryOptions{
		Limit:        10,
		Offset:       0,
		SortField:    "created_at",
		SortDiretion: "asc",
	}

	if page, ok := query["page"]; ok {
		pageInt, err := strconv.Atoi(page[0])
		if err != nil {
			return nil, types.BadRequest("invalid page value")
		}
		// calculate the offset based on page number
		opts.Offset = (pageInt - 1) * opts.Limit
	}

	if pageSize, ok := query["page_size"]; ok {
		pageSizeInt, err := strconv.Atoi(pageSize[0])
		if err != nil {
			return nil, types.BadRequest("invalid pageSize value")
		}
		opts.Limit = pageSizeInt
	}

	if sort, ok := query["sort"]; ok {
		sortParts := strings.Split(sort[0], "-")
		if len(sortParts) != 2 {
			return nil, types.BadRequest("invalid sort value")
		}
		opts.SortField = sortParts[0]
		opts.SortDiretion = sortParts[1]
	}

	return opts, nil
}
