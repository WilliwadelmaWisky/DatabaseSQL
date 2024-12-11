package util

//
//
// # Arguments
//  - in
//  - proc
//  - returns
func Map[TIn, TOut any](in []TIn, proc func(TIn) TOut) []TOut {
	out := make([]TOut, len(in))
	for i, v := range in {
		out[i] = proc(v)
	}

	return out
}

//
//
// # Arguments
//  - in
//  - match
//  - returns
func Filter[T any](in []T, match func(T) bool) []T {
	out := make([]T, 0)
	for _, v := range in {
		if match(v) {
			out = append(out, v)
		}
	}

	return out
}

//
//
// # Arguments
//  - in
//  - value
//  - returns
func Contains[T comparable](in []T, value T) bool {
	for _, v := range in {
		if v == value {
			return true
		}
	}

	return false
}
