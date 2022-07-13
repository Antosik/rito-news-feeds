package internal

import (
	"fmt"
)

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

func IsEqual[T interface{}](a, b []T, comparator func(a, b T) bool) bool {
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

func CompareAndGetDiff[T interface{}](
	old, new []T,
	keyfunc func(item T) string,
) (diff []string, isEqual bool) {
	om := make(map[string]uint8, len(old))

	for _, olditem := range old {
		om[keyfunc(olditem)] = om[keyfunc(olditem)] + 1
	}

	for _, newitem := range new {
		key := keyfunc(newitem)

		val, found := om[key]
		if found {
			om[key] = val - 1
		} else {
			diff = append(diff, fmt.Sprintf("+ %s", key))
		}
	}

	for key, value := range om {
		switch {
		case value < 0:
			diff = append(diff, fmt.Sprintf("+ %s", key))
		case value > 0:
			diff = append(diff, fmt.Sprintf("- %s", key))
		}
	}

	return diff, len(diff) == 0
}
