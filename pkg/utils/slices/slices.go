package slices

func Filter[T comparable](filters, list []T) []T {
	filterMap := make(map[T]bool)
	for _, f := range filters {
		filterMap[f] = true
	}

	result := make([]T, 0)
	for _, elem := range list {
		if !filterMap[elem] {
			result = append(result, elem)
		}
	}

	return result
}

func Contains[T comparable](list []T, search T) bool {
	for _, element := range list {
		if element == search {
			return true
		}
	}

	return false
}

func NotIn[T comparable](list []T, search T) bool {
	for _, element := range list {
		if element == search {
			return false
		}
	}

	return true
}
