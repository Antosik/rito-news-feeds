package internal

func SplitSliceToChunks[T interface{}](arr []T, count int) [][]T {
	var (
		size   = int(float64(len(arr)) / float64(count))
		exceed = len(arr) % count
		result = make([][]T, count)
	)

	pointer := 0
	for i := 0; i < count; i++ {
		start := pointer
		pointer = pointer + size
		if i < exceed {
			pointer = pointer + 1
		}
		if pointer > len(arr) {
			pointer = len(arr)
		}

		result[i] = make([]T, pointer-start)
		copy(result[i], arr[start:pointer])

	}

	return result
}

func IsEqual[T comparable](a, b []T, comparator func(a, b T) bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !comparator(a[i], b[i]) {
			return false
		}
	}
	return true
}
