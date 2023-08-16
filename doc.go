/*
Carbo provides building blocks to compose data processing pipeline with concurrency.

Carbo's core component is task.Task, which represents a process that takes inputs and outputs.

You can build an entire build pipeline with a number of tasks.

Althouogh it is possible to define a task directly with task.FromFn, it is recommended to use sub-types of task.Task: source, pipe and sink.

# Sub-types of task.Task

A data pipeline is built with three sub-types of task.Task: source, pipe and sink.

The basic form of a data pipeline should look like: sink -> pipe -> ... -> pipe -> sink

# Source

A source is a special type of task.Task that takes an empty input that will be closed immediately after a data pipeline starts.

This is used as an entry point of a data pipeline.
For example, source.FromSlice takes a slice as an argument, and then produces each element of the slice as its outputs.

# Pipe

A pipe is a similar component to task.Task, except the passed output channel is closed automatically after its process is finished.

This is used to convert inputs into outputs.
For example, pipe.Map processes inputs one by one, and produces corresponding outputs.
This can be used in case an input and an output has one-to-one correspondence.

# Sink

A sink is a special type of task.Task that takes an empty output. This is just like an opposite of a source.

This is used as an last component of a data pipeline.
For example, sink.ToSlice takes a pointer to a slice, and appends elements from its inputs.

# Connecting tasks

Each task can be connected with another one using task.Connect if the type of an upstream task's output and the type of a downstream task's input match.

task.Connect also returns task.Task so it can be chained. An entire data pipeline can be built by connecting multiple tasks in this way.

# Running a data pipeline

Carbo has a Flow component which is a wrapper of a task that takes an empty input and an empty output.
This kind of task is typically built as a chain of tasks that starts from a source and ends with a sink.
*/
package carbo
