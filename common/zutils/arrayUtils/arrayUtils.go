package arrayUtils

func ArrayIn[T comparable](list []T, found T) bool {
	for _, s := range list {
		if found == s {
			return true
		}
	}
	return false
}