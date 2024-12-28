package sql

import (
	"fmt"
	"strings"
)

type SortDirection int

const (
	DIRECTION_ASCENDING  SortDirection = 1  // Ascending order
	DIRECTION_DESCENDING SortDirection = -1 // Descending order
)

type Sorter struct {
	ColumnName string        // Column name to sort by
	Direction  SortDirection // Order of the sorting
}

func GetSortDirection(s string) (SortDirection, error) {
	switch strings.ToUpper(s) {
	case "ASC":
		return DIRECTION_ASCENDING, nil
	case "DESC":
		return DIRECTION_DESCENDING, nil
	}

	return -1, fmt.Errorf("invalid sort direction %s", s)
}
