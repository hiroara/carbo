package behavior_test

import (
	"context"

	"github.com/hiroara/carbo/cache/internal/behavior"
)

type dummyEntry struct {
	value string
}

func (e *dummyEntry) Get(ctx context.Context) (*string, error) {
	if e.value == "" {
		return nil, nil
	}
	return &e.value, nil
}

func (e *dummyEntry) Set(ctx context.Context, v string) error {
	e.value = v
	return nil
}

func createBehavior(t behavior.BehaviorType) (behavior.Entry[string], behavior.Behavior[string, string]) {
	ent := behavior.Entry[string](&dummyEntry{})
	return ent, behavior.New[string](ent, t)
}
