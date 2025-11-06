package main

import (
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func generateRandomArray(length int, min, max int32) []int32 {
	arr := make([]int32, length)

	for i := 0; i < length; i++ {
		arr[i] = min + r.Int31n(max-min+1)
	}

	return arr
}
