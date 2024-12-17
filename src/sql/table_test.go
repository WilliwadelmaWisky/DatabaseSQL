package sql

import (
	"testing"
)

func TestTableInsert(t *testing.T) {
	table := &Table{Columns: []*Column{{Name: "col1", Type: TYPE_INT}, {Name: "col2", Type: TYPE_VARCHAR}}}
	err := table.Insert([]RowData{{ColName: "col1", Value: "1"}, {ColName: "col2", Value: "val2"}})
	if err != nil {
		t.Fatal("insert returned an error but should not have")
	}
}

func TestTableGet(t *testing.T) {
	table := &Table{Columns: []*Column{{Name: "col1", Type: TYPE_INT, Values: []string{"1", "2", "3"}}}}
	_, err := table.Get([]string{"col1"}, []*Filter{})
	if err != nil {
		t.Fatal("get returned an error but should not have")
	}
}

func TestTableUpdate(t *testing.T) {
	table := &Table{Columns: []*Column{{Name: "col1", Type: TYPE_INT, Values: []string{"1", "2", "3"}}}}
	err := table.Update([]RowData{{ColName: "col1", Value: "5"}}, []*Filter{})
	if err != nil {
		t.Fatal("update returned an error but should not have")
	}
}

func TestTableDelete(t *testing.T) {
	table := &Table{Columns: []*Column{{Name: "col1", Type: TYPE_INT, Values: []string{"1", "2", "3"}}}}
	err := table.Delete([]*Filter{})
	if err != nil {
		t.Fatal("delete returned an error but should not have")
	}
}
