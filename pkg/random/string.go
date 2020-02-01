package random

import (
	"math/rand"
	"time"
)

const (
	DictAlpha       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DictNumber      = "0123456789"
	DictAlphaNumber = DictAlpha + DictNumber
)

func String(length uint, dict string) (str string) {
	rand.Seed(time.Now().UnixNano())
	for i := uint(0); i < length; i++ {
		index := rand.Uint64() % uint64(len(dict))
		str += dict[index : index+1]
	}
	return
}
