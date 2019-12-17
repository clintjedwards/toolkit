package listutil

// RemoveStringFromList removes an element from an array of string
// does not preserve list order
func RemoveStringFromList(list []string, value string) []string {
	for index, item := range list {
		if item == value {
			list[index] = list[len(list)-1]
			return list[:len(list)-1]
		}
	}

	return list
}

// FindListDifference returns list elements that are in list A
// but not found in B
func FindListDifference(a, b []string) []string {
	m := make(map[string]bool)
	diff := []string{}

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

// FindListUpdates is used to compare a new and old version of lists
// it will compare the old version to the new version and return
// which elements have been added or removed from the old list
func FindListUpdates(oldList, newList []string) (additions, removals []string) {

	removals = FindListDifference(oldList, newList)
	additions = FindListDifference(newList, oldList)

	return additions, removals
}
