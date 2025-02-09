package utils

import (
	"testing"

	"github.com/mwdev22/CarRental/internal/types"
)

func TestBuildBatchQuery(t *testing.T) {
	query := "SELECT * FROM users WHERE 1=1"
	filters := []*types.QueryFilter{
		{
			Field:    "id",
			Operator: "=",
			Value:    1,
		},
		{
			Field:    "email",
			Operator: "=",
			Value:    "",
		},
	}

	opts := &types.QueryOptions{
		SortField:    "id",
		SortDiretion: "ASC",
		Limit:        10,
		Offset:       0,
	}

	expectedQuery := "SELECT * FROM users WHERE 1=1 AND id = ? AND email = ? ORDER BY id ASC LIMIT ? OFFSET ?"

	result, args := BuildBatchQuery(query, filters, opts)

	if result != expectedQuery {
		t.Errorf("expected %s, got %s", expectedQuery, result)
	}

	if len(args) != 4 {
		t.Errorf("expected 4 args, got %d", len(args))
	}

}
