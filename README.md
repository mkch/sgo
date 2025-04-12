# sgo

[简体中文](README.zh-CN.md)

sgo: Structured Concurrency Infrastructure for Go.

sgo provides higher-level abstractions by encapsulating Go's concurrency primitives for common concurrent operations. Using sgo can significantly reduce the complexity of concurrent programming, thereby improving its efficiency and robustness.

## Motivation

The design philosophy of sgo originates from [Notes on structured concurrency, or: Go statement considered harmful](https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/).

The article points out that using the `go` keyword directly for concurrent programming is as primitive and harmful as using `goto` statements for flow control in early FLOW-MATIC languages. This low-level concurrency control breaks code abstraction and hinders modular development and local reasoning. Just as structured programming replaced `goto` with sequence, selection, and loop constructs, structured concurrency requires high-level abstractions to replace direct use of `go`.

## Basic Elements of Structured Concurrency

Based on practical experience, Go's concurrency requirements can be categorized into three typical scenarios:

1. Parallel execution and waiting

    For example, executing multiple initialization tasks concurrently during program startup, then proceeding after all tasks complete.

2. Parallel computation and result collection

    For example, splitting a large array into slices for concurrent summation, then aggregating the results.

3. Race execution and optimal selection

    For example, querying multiple search engines concurrently and returning the first valid result.

These three requirements each have standard implementation patterns, but correctly and completely implementing them often requires:

- Rich experience
- Meticulous attention
- Comprehensive testing

Reinventing these patterns each time places significant cognitive burden on developers. Following the "structured concurrency" philosophy, encapsulating these patterns into unified abstract structures can systematically solve these problems. sgo encapsulates the above three patterns as:

1. sgo.Group

    ```go
    sgo.NewGroup().
        Go(setupStep1).
        Go(setupStep2).
        Wait()
    ```

2. sgo.Collector

    ```go
    result := sgo.NewCollector[int]().
        Go(func() int { return sumSlice(s[:99]) }).
        Go(func() int { return sumSlice(s[99:]) }).
        Collect()
    sum := result[0] + result[1]
    ```

3. sgo.Racer

    ```go
    result := sgo.NewRacer[string](context.Background()).
        Go(func(ctx context.Context) string {
            return searchGoogle(ctx, keyword);
        }).
        Go(func(ctx context.Context) string {
           return searchBing(ctx, keyword);
        }).
        Collect()
    ```

## Reflections on "Structured Concurrency"

The essence of structured concurrency lies in abstracting proven concurrent patterns into reusable programming components. This philosophy shares the same lineage as structured programming's transformation of `goto` statements into standard control structures. Its core value lies not in the implementation technology itself, but in the programming paradigm innovation it advocates.
