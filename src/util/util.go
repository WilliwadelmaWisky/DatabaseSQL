package util

func Map[TIn, TOut any](arr []TIn, proc func(TIn) TOut) []TOut {
	out := make([]TOut, len(arr))
	for i, v := range arr {
		out[i] = proc(v)
	}

	return out
}

func Filter[T any](arr []T, match func(T) bool) []T {
	out := make([]T, 0)
	for _, v := range arr {
		if match(v) {
			out = append(out, v)
		}
	}

	return out
}
