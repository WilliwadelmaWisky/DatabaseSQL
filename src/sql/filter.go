package sql

// Represents a single sql where expression
type Filter struct {
	ColumnName   string
	Operator     EqualityOperator
	CompareValue string
}

// Check if a value is included by the filter
func (filter *Filter) IsIncluded(value string, t ColumnType) bool {
	return filter.Operator.Compare(t, value, filter.CompareValue)
}

// Check if any filter returns true
func IsTrueForAny(filters []*Filter, match func(filter *Filter) bool) bool {
	for _, filter := range filters {
		if match(filter) {
			return true
		}
	}

	return false
}
