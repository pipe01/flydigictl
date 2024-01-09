package utils

func SplitSlice[T any](s []T, chunkSize int) [][]T {
	var chunks [][]T

	i := 0
	for ; i < len(s)-chunkSize; i += chunkSize {
		chunks = append(chunks, s[i:i+chunkSize])
	}
	if i < len(s) {
		chunks = append(chunks, s[i:])
	}

	return chunks
}

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
