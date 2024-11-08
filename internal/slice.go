package internal

func SplitSliceToChunks[T interface{}](arr []T, count int) [][]T {
	var (
		size   = int(float64(len(arr)) / float64(count))
		exceed = len(arr) % count
		result = make([][]T, count)
	)

	pointer := 0
	for index := range count {
		start := pointer
		pointer += size

		if index < exceed {
			pointer++
		}

		if pointer > len(arr) {
			pointer = len(arr)
		}

		result[index] = make([]T, pointer-start)
		copy(result[index], arr[start:pointer])
	}

	return result
}

func IsEqual[T interface{}](first, second []T, comparator func(a, b T) bool) bool {
	if len(first) != len(second) {
		return false
	}

	for i := range first {
		if !comparator(first[i], second[i]) {
			return false
		}
	}

	return true
}

func CompareAndGetDiff[T interface{}](
	olditems, newitems []T,
	keyfunc func(item T) string,
) ([]string, bool) {
	changesmap := make(map[string]int8, len(olditems))
	diff := []string{}

	for _, olditem := range olditems {
		changesmap[keyfunc(olditem)] = changesmap[keyfunc(olditem)] + 1
	}

	for _, newitem := range newitems {
		key := keyfunc(newitem)

		val, found := changesmap[key]
		if found {
			changesmap[key] = val - 1
		} else {
			diff = append(diff, "+ "+key)
		}
	}

	for key, value := range changesmap {
		switch {
		case value < 0:
			diff = append(diff, "+ "+key)
		case value > 0:
			diff = append(diff, "- "+key)
		}
	}

	return diff, len(diff) == 0
}
