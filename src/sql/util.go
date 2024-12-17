package sql

// Contruct a new array by proc
func Map[TIn, TOut any](in []TIn, proc func(TIn) TOut) []TOut {
	out := make([]TOut, len(in))
	for i, v := range in {
		out[i] = proc(v)
	}

	return out
}

//
// func Filter[T any](in []T, match func(T) bool) []T {
// 	out := make([]T, 0)
// 	for _, v := range in {
// 		if match(v) {
// 			out = append(out, v)
// 		}
// 	}

// 	return out
// }
