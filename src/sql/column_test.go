package sql

import (
	"strings"
	"testing"
)

func TestColumnTypeGetDefaultValue(t *testing.T) {
	val, err := TYPE_INT.GetDefaultValue()
	if val != "0" || err != nil {
		t.Fatal("wrong default value configured, expected=0")
	}

	val, err = TYPE_VARCHAR.GetDefaultValue()
	if val != "" || err != nil {
		t.Fatal("wrong default value configured, expected empty string")
	}
}

func TestColumnTypeToString(t *testing.T) {
	val := TYPE_INT.ToString()
	if strings.ToUpper(val) != "INT" {
		t.Fatal("wrong to_string value configured, expected int")
	}

	val = TYPE_VARCHAR.ToString()
	if strings.ToUpper(val) != "VARCHAR" {
		t.Fatal("wrong to_string value configured, expected varchar")
	}
}
