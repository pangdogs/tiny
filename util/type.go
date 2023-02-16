package util

// Zero 创建零值
func Zero[T any]() T {
	var zero T
	return zero
}
