package sql

// Represents a single column in a table
type Column struct {
	Name   string     `json:"column"`
	Type   ColumnType `json:"type"`
	Values []string   `json:"values"`
}
