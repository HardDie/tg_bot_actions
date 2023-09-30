package utils

func Allocate[T any](val T) *T {
	return &val
}
