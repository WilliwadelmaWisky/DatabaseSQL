package sql

import (
	"fmt"
	"strings"
)

type ColumnType int

const (
	TYPE_INT ColumnType = iota
	TYPE_VARCHAR
)

func GetType(s string) (ColumnType, error) {
	switch strings.ToUpper(s) {
	case "INT":
		return TYPE_INT, nil
	case "VARCHAR":
		return TYPE_VARCHAR, nil
	}

	return -1, fmt.Errorf("invalid column type: %s", s)
}
