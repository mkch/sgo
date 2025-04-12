# sgo

[English](README.md)

sgo: Go语言结构化并发基础设施。

sgo 通过封装Go语言并发原语，为常用并发操作提供更高级的抽象。使用sgo能够显著降低并发编程的复杂度，从而提高并发编程的效率和鲁棒性。

## 动机

sgo的涉及理念来源于[Notes on structured concurrency, or: Go statement considered harmful](https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/)。

文中指出，直接使用`go`关键字进行并发编程，其原始性堪比早期FLOW-MATIC语言中使用`goto`语句控制流程。这种低层次的并发控制会破坏代码抽象能力，不利于模块化开发和局部逻辑推理。正如结构化编程用顺序、选择、循环结构取代`goto`，结构化并发也需要通过高层抽象来取代对`go`的使用。

## 结构化并发基本要件

根据实践经验，Go语言的并发需求主要可分为三类典型场景：

1. 并行执行与等待

    例如程序启动时并发执行多个初始化任务，全部完成后继续后续流程。

2. 并行计算与结果收集

    例如将大数组分片并发求和，最后汇总结果。

3. 竞速执行与择优选取

    例如并发查询多个搜索引擎，返回首个有效结果。

以上三种需求都有其标准实现范式，但要正确和完整地实现这些范式往往需要丰富的经验、周密的心思和全面的测试。每次都“重新制造轮子”将为代码实现者带来了显著的心智负担。按照"结构化并发"的理念，将这些范式封装为统一的抽象结构，可以系统性地解决这些问题。以上3种范式被sgo封装为：

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

## 关于“结构化并发”的思考

结构化并发的本质在于将经过实践验证的并发模式抽象为可复用的编程构件。这一理念与结构化编程中将`goto`语句转化为标准控制结构的思路一脉相承。其核心价值不在于具体实现技术，而在于其所倡导的编程范式革新。
