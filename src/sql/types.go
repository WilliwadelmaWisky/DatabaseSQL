package sql

import (
	"fmt"
	"strings"
)

// Enum to represent all the possible datatypes in the database, values named with TYPE prefix
type ColumnType int

const (
	// Represents an integer value, a whole number
	TYPE_INT ColumnType = iota
	// Represents a string value
	TYPE_VARCHAR
)

// Get a datatype based of a string
func GetType(s string) (ColumnType, error) {
	switch strings.ToUpper(s) {
	case "INT":
		return TYPE_INT, nil
	case "VARCHAR":
		return TYPE_VARCHAR, nil
	}

	return -1, fmt.Errorf("invalid column type: %s", s)
}

// Get a default value of a datatype
func (Type ColumnType) GetDefaultValue() (string, error) {
	switch Type {
	case TYPE_INT:
		return "0", nil
	case TYPE_VARCHAR:
		return "", nil
	}

	return "", fmt.Errorf("column type does not have default value: %s", Type.ToString())
}

// Get a string value of a datatype
func (Type ColumnType) ToString() string {
	switch Type {
	case TYPE_INT:
		return "INT"
	case TYPE_VARCHAR:
		return "VARCHAR"
	}

	return "NULL"
}
