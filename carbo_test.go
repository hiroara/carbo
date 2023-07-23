package carbo_test

import (
	"context"
	"fmt"
	"log"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/runner"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
)

type MyConfig struct {
	StringField string `yaml:"string_field"`
	IntField    int    `yaml:"int_field"`
}

// Build a flow and directly run it.
func Example_flow() {
	ss := source.FromSlice([]string{"a", "b", "c"})
	ds := task.Connect(
		ss.AsTask(),
		pipe.Map(func(ctx context.Context, s string) (string, error) {
			return s + s, nil
		}).AsTask(),
		1,
	)
	pr := task.Connect(
		ds,
		sink.ElementWise(func(ctx context.Context, s string) error {
			fmt.Println(s)
			return nil
		}).AsTask(),
		1,
	)

	err := flow.FromTask(pr).Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// aa
	// bb
	// cc
}

// Define a flow factory function to build a flow with a config struct, and run the flow.
func Example_flowFactory() {
	fac := func(cfg *MyConfig) (*flow.Flow, error) {
		ss := source.FromSlice([]string{cfg.StringField})
		pr := task.Connect(
			ss.AsTask(),
			sink.ElementWise(func(ctx context.Context, s string) error {
				fmt.Println(s)
				return nil
			}).AsTask(),
			1,
		)
		return flow.FromTask(pr), nil
	}

	flow.RunWithConfig(context.Background(), fac, "testdata/config.yaml")
	// Output:
	// value-from-string-field
}

// Define multiple flow factories, register them to a runner, and run a flow.
// This is useful to make an executable that takes a subcommand.
func Example_runner() {
	fac1 := func() (*flow.Flow, error) {
		ss := source.FromSlice([]string{"item1"})
		pr := task.Connect(
			ss.AsTask(),
			sink.ElementWise(func(ctx context.Context, s string) error {
				fmt.Println(s)
				return nil
			}).AsTask(),
			1,
		)
		return flow.FromTask(pr), nil
	}
	fac2 := func(cfg *MyConfig) (*flow.Flow, error) {
		ss := source.FromSlice([]int{cfg.IntField})
		pr := task.Connect(
			ss.AsTask(),
			sink.ElementWise(func(ctx context.Context, i int) error {
				fmt.Println(i)
				return nil
			}).AsTask(),
			1,
		)
		return flow.FromTask(pr), nil
	}

	r := runner.New()
	r.Define("flow1", flow.NewFactory(fac1))
	r.Define("flow2", flow.NewFactoryWithConfig(fac2, "testdata/config.yaml"))
	r.Run(context.Background(), "flow2")
	// Output:
	// 100
}
