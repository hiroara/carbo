# Carbo

Carbo is a package that contains a framework for Gophers to provide a way to build a data pipeline.

There are multiple headaches when building an application containing a data pipeline, like controlling concurrency, scalability, back pressure, etc.
This framework is built to deal with these issues without much effort.

## Why Carbo?

As far as I know, there are many great frameworks to control tasks and those dependencies in complex workflows, and such a tool also provides a way to monitor tasks' situations.

But, there are cases where a programmer thinks configuring a cluster and running a large number of tasks there can be overkill.

Carbo would fit such a case. It is a pure Golang implementation that helps run small tasks in a process with easy control of concurrencies.

Additionally, Carbo also provides an easy way to feed data from one process to another with gRPC. This way provides enough scalability in many cases.

## Exposing / pulling data through gRPC

As described above, Carbo provides an easy way to feed data from one process to another with gRPC.

In this way, for example, you can separate a data pipeline into a CPU-intensive part and an IO-intensive part as different processes, and run it with a different concurrency limit.

Additionally, this means that Carbo doesn't necessarily force you to stick to this framework itself or even Golang, thanks to the programming language-agnostic RPC protocol, gRPC.

For example, you can pull data from a Golang process that uses Carbo with grpcurl for debugging.

Or, you can also write another program, for example, in Python, to pull data via gRPC.
This is convenient, for example, when you want to write a fast data pipeline in Golang and feed the output into Python to build a machine-learning model with scikit-learn.
