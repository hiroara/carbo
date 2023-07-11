package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/config"
	"github.com/hiroara/carbo/internal/testutils"
)

func TestParse(t *testing.T) {
	t.Parallel()

	var cfg testutils.Config
	err := config.Parse("../testdata/config.yaml", &cfg)
	require.NoError(t, err)

	assert.Equal(t, "thisisstring", cfg.StringField)
	assert.Equal(t, 100, cfg.IntField)
}
