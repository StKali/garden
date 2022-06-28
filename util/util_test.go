package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := RandInternalString(6, 72)
	var (
		hash string
		err  error
	)
	for i := 0; i < 100; i++ {
		hash, err = HashPassword(password)
		require.NoError(t, err)
		err = VerifyPassword(password, hash)
		require.NoError(t, err)
	}

	bigPassword := RandString(1024)
	hash, err = HashPassword(bigPassword)
	require.Error(t, err)
	require.Equal(t, "", hash)
}

func TestCheckError(t *testing.T) {
	CheckError("no err", nil)
}
