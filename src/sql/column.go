package sql

// Represents a single column in a table
type Column struct {
	Name   string     `json:"column"` // Column name
	Type   ColumnType `json:"type"`   // Column variable type
	Values []string   `json:"values"` // Column data
}
