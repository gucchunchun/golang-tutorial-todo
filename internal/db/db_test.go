package db

import (
	"testing"

	"github.com/stretchr/testify/require"

	"golang/tutorial/todo/internal/config"
)

func TestDBConnect_UnreachableDB(t *testing.T) {
	t.Parallel()

	cfg := config.DBConfig{
		Host:     "127.0.0.1",
		Port:     "1",
		User:     "root",
		Password: "example",
		Name:     "testdb",
	}

	var c Client
	err := c.Connect(cfg)
	require.Error(t, err, "expected Connect to fail for unreachable DB")
	require.ErrorContains(t, err, "failed to ping DB")
	require.Nil(t, c.DB(), "DB should remain nil on failed connect")
}
