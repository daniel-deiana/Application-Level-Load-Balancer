package helper

func RemoveByValue[T comparable](s []T, value T) []T {
	temp := make([]T, 0, len(s))
	for _, elem := range s {
		if elem != value {
			temp = append(temp, elem)
		}
	}
	return temp
}
