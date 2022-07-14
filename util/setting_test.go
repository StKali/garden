package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSetting(t *testing.T) {
	s := GetSetting()
	require.True(t, s == &defaultSetting)
}
