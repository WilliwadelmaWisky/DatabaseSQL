package sql

type SortDirection int

const (
	ASC SortDirection = iota
	DESC
)

type Sorter struct {
	ColumnName string
	Direction  SortDirection
}
