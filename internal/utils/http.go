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
		if len(values) < 1 || key == "page" || key == "page_size" || key == "sort" {
			continue
		}
		if strings.Contains(key, "[") && strings.HasSuffix(key, "]") {
			// scrape field and operator from key
			// field[gt] --> field, gt
			field := key[:strings.Index(key, "[")]
			operatorKey := key[strings.Index(key, "[")+1 : len(key)-1]
			if operator, ok := types.OperatorMap[operatorKey]; ok {
				//
				value := values[0]
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
					Field:    field,
					Operator: operator,
					Value:    value,
				})
			} else {
				return nil, types.BadQueryParameter(fmt.Sprintf("Invalid operator in filter: %s", operatorKey))
			}
		} else {
			// default for not provided operator
			filters = append(filters, &types.QueryFilter{
				Field:    key,
				Operator: "=",
				Value:    values[0],
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
		SortField:    "created",
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
