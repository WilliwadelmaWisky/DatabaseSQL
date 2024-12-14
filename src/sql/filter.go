package sql

type Filter struct {
	ColumnName   string
	Operator     EqualityOperator
	CompareValue string
}

func (filter *Filter) IsIncluded(value string, t ColumnType) bool {
	return filter.Operator.Compare(t, value, filter.CompareValue)
}

func IsTrueForAny(filters []*Filter, match func(filter *Filter) bool) bool {
	for _, filter := range filters {
		if match(filter) {
			return true
		}
	}

	return false
}
