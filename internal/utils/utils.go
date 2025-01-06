package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mwdev22/CarRental/internal/types"
)

func MakeLogger(filename string) (*log.Logger, error) {
	logFile, err := os.OpenFile(fmt.Sprintf("log/%s.log", filename), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	logger := log.New(logFile, "", log.LstdFlags)
	return logger, nil
}

func BuildBatchQuery(query string, filters []*types.QueryFilter, opts *types.QueryOptions) (string, []interface{}) {
	args := make([]interface{}, 0)
	// i use question marks because filters and their count are dynamic
	// could be 1, 2, 3 etc through loop, but its not necessary complexity i think
	for _, filter := range filters {
		query += fmt.Sprintf(" AND %s %s ?", filter.Field, filter.Operator)
		args = append(args, filter.Value)
	}

	if opts != nil {
		query += ` ORDER BY ? ?`
		query += ` LIMIT ? OFFSET ?`
		args = append(args, opts.SortField, opts.SortDiretion, opts.Limit, opts.Offset)
	} else {
		query += ` ORDER BY id DESC`
	}
	return query, args
}

func GenerateUniqueString(base string) string {
	return fmt.Sprintf("%s_%d", base, time.Now().UnixNano())
}
