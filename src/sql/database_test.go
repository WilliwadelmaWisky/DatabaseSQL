package sql

import (
	"testing"
)

func TestDatabaseGet(t *testing.T) {
	database := NewDatabase("", &Table{Name: "valid"})
	table, err := database.Get("valid")
	if table == nil || err != nil {
		t.Fatal("table should be found but was not")
	}

	_, err = database.Get("invalid")
	if err == nil {
		t.Fatal("error was not thrown but should have")
	}
}
