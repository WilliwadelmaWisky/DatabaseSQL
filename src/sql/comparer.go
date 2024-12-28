package sql

import (
	"strconv"
	"strings"
)

type EqualityOperator int

const (
	LESS             EqualityOperator = 1               // Less than operator <
	GREATER          EqualityOperator = 2               // Greater than operator >
	EQUAL            EqualityOperator = 4               // Equals to operator =
	LESS_OR_EQUAL    EqualityOperator = LESS | EQUAL    // Less than or equals to operator <=
	GREATER_OR_EQUAL EqualityOperator = GREATER | EQUAL // Greater than or equals to operator >=
)

// Get equality operator from string
func GetEqualityOperator(s string) EqualityOperator {
	switch s {
	case "<":
		return LESS
	case ">":
		return GREATER
	case "=":
		return EQUAL
	case "<=":
		return LESS_OR_EQUAL
	case ">=":
		return GREATER_OR_EQUAL
	}

	return -1
}

// Inverse of the equality operator, less than changes to greater than etc
func (operator EqualityOperator) Inverse() EqualityOperator {
	switch operator {
	case LESS:
		return GREATER
	case LESS_OR_EQUAL:
		return GREATER_OR_EQUAL
	case EQUAL:
		return EQUAL
	case GREATER:
		return LESS
	case GREATER_OR_EQUAL:
		return LESS_OR_EQUAL
	}

	return -1
}

// Compare values, based of type
func (operator EqualityOperator) Compare(t ColumnType, a string, b string) bool {
	switch t {
	case TYPE_INT:
		return operator.compareInt(a, b)
	case TYPE_VARCHAR:
		return operator.compareString(a, b)
	}

	return false
}

func Compare(t ColumnType, a string, b string) int {
	switch t {
	case TYPE_INT:
		inta, _ := strconv.Atoi(a)
		intb, _ := strconv.Atoi(b)
		return inta - intb
	case TYPE_VARCHAR:
		return strings.Compare(a, b)
	}

	return 0
}

func (operator EqualityOperator) compareInt(a string, b string) bool {
	intValue, _ := strconv.Atoi(a)
	intCompareValue, _ := strconv.Atoi(b)

	switch operator {
	case LESS:
		return intValue < intCompareValue
	case LESS_OR_EQUAL:
		return intValue <= intCompareValue
	case EQUAL:
		return intValue == intCompareValue
	case GREATER:
		return intValue > intCompareValue
	case GREATER_OR_EQUAL:
		return intValue >= intCompareValue
	}

	return false
}

func (operator EqualityOperator) compareString(a string, b string) bool {
	if operator == EQUAL {
		return a == b
	}

	return false
}
