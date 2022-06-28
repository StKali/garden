package util

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandString(t *testing.T) {
	for i := 0; i < 100; i++ {
		require.Equal(t, len(RandString(i)), i)
	}
}

func TestRandEmail(t *testing.T) {
	for i := 0; i < 100; i++ {
		email := RandEmail()
		require.Contains(t, email, "@")
	}
}

func TestRandIntervalString(t *testing.T) {
	for i := 0; i < 10; i++ {
		min := rand.Intn(1024)
		max := min + rand.Intn(1024)
		str1 := RandInternalString(min, max)
		str2 := RandInternalString(max, min)
		require.True(t, len(str1) >= min && len(str1) <= max)
		require.True(t, len(str2) >= min && len(str2) <= max)
	}
}

func TestRandIP(t *testing.T) {

	for i := 0; i < 100; i++ {
		ip := RandIP()
		seg := strings.Split(ip, ".")
		require.Equal(t, 4, len(seg))
	}
}
