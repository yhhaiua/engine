package util

import (
	"math/rand"
	"sync"
	"time"
)

var (
	r  = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu sync.Mutex
)

// Int implements rand.Int on the grpcrand global source.
func Int() int {
	mu.Lock()
	defer mu.Unlock()
	return r.Int()
}

// Int63n implements rand.Int63n on the grpcrand global source.
func Int63n(n int64) int64 {
	if n <= 0 {
		return 0
	}
	mu.Lock()
	defer mu.Unlock()
	return r.Int63n(n)
}

// Intn implements rand.Intn on the grpcrand global source.
func Intn(n int) int {
	if n <= 0 {
		return 0
	}
	mu.Lock()
	defer mu.Unlock()
	return r.Intn(n)
}

// Float64 implements rand.Float64 on the grpcrand global source.
func Float64() float64 {
	mu.Lock()
	defer mu.Unlock()
	return r.Float64()
}

// Uint64 implements rand.Uint64 on the grpcrand global source.
func Uint64() uint64 {
	mu.Lock()
	defer mu.Unlock()
	return r.Uint64()
}

func RandInt(min, max int) int {
	if min >= max {
		return max
	}
	return Intn(max-min+1) + min
}

//Shuffle 随机数组
func Shuffle[T any](slice []T) {
	if len(slice) <= 1 {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
