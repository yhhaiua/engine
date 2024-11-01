package util

type Compare interface {
	~int | ~int64 | ~int32 | ~int8 | ~uint | ~float32 | ~float64
}

func Min[T Compare](x, y T) T {
	if x < y {
		return x
	}
	return y
}
func Max[T Compare](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Abs[T Compare](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
