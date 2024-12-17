package sql

import "fmt"

type Column struct {
	Name   string     `json:"column"`
	Type   ColumnType `json:"type"`
	Values []string   `json:"values"`
}

func (Type ColumnType) GetDefaultValue() (string, error) {
	switch Type {
	case TYPE_INT:
		return "0", nil
	case TYPE_VARCHAR:
		return "", nil
	}

	return "", fmt.Errorf("column type does not have default value: %s", Type.ToString())
}

func (Type ColumnType) ToString() string {
	switch Type {
	case TYPE_INT:
		return "INT"
	case TYPE_VARCHAR:
		return "VARCHAR"
	}

	return "NULL"
}
