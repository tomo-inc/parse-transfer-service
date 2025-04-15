package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig(t *testing.T) {
	config := GetConfig("")
	require.True(t, config.ListenHost != "")
	require.True(t, config.AlertConfig.LarkBotId != "")
	require.True(t, config.EVMEndpoints["100"].Endpoint != "")
	require.False(t, config.EVMEndpoints["100"].SupportDebug)
	require.True(t, config.AlertConfig.Interval != 0)

	config = GetConfig("./config2.yaml")
	require.True(t, config.EVMEndpoints["5600"].Endpoint != "")
}
