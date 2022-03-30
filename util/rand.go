package util

import (
	"math/rand"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {

	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func RandInternalString(min, max int) string {
	var n int
	if max < min {
		n = min - max
	} else {
		n = max - min
	}
	n = min + rand.Intn(n)
	return RandString(n)
}

var emailSuffixs = []string{
	"@google.com",
	"@yahoo.com",
	"@mit.edu",
	"@163.com",
	"@outlook.com",
}

func RandEmail() string {
	prefix := RandInternalString(1, 32)
	return prefix + emailSuffixs[len(prefix)%len(emailSuffixs)]
}
