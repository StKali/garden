package util

import (
	"math/rand"
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
	for i := 0; i < 100; i++ {
		min := rand.Intn(1024)
		max := min + rand.Intn(1024)
		str := RandInternalString(min, max)
		require.True(t, len(str) >= min && len(str) <= max)
	}
}
