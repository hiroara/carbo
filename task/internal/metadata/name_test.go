package metadata_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiroara/carbo/task/internal/metadata"
)

func TestName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	assert.Equal(t, "<Anonymous Task>", metadata.GetName(ctx))

	ctx = metadata.WithName(ctx, "a-new-name")
	assert.Equal(t, "a-new-name", metadata.GetName(ctx))
}
