package utils

func Repeat[T any](v T, n int) []T {
	arr := make([]T, n)

	for i := 0; i < n; i++ {
		arr[i] = v
	}

	return arr
}

func RepeatFunc[T any](v func() T, n int) []T {
	arr := make([]T, n)

	for i := 0; i < n; i++ {
		arr[i] = v()
	}

	return arr
}
