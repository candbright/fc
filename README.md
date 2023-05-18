# FUNCTION ROUTER

## 介绍

FUNCTION ROUTER用于定义一系列可同步或异步执行的函数链，不同的函数执行链根据注册时的key值区分，每个函数都有其相应的回滚函数。

## 回滚

当函数链为同步函数链时，如果执行某个函数产生error，那么会中断此函数链，并且已经执行过的所有函数的回滚函数都会被执行。

当函数链为异步函数链时，执行某个函数产生error并不会中断函数链后续执行，当函数链中的所有函数都执行完毕后产生了error的函数的回滚函数会被依次执行。

例如定义了以下函数链（括号中为函数对应的回滚函数）:

A(Arb) - B(Brb) - C(Crb) - D(Drb) - E(Erb)

函数链的正常执行流程：A - B - C - D - E

当函数链为同步函数链时，执行C产生error，那么执行流程为：A - B - C - Crb - Brb - Arb

当函数链为异步函数链时，执行C产生error，那么执行流程为：A - B - C - D - E - Crb 当函数链为异步函数链时，执行C、E产生error，那么执行流程为：A - B - C - D - E - Crb - Erb

## 常用函数

##### RegisterSync(key string, syncFunc func(instance T) error, rollback func(instance T))

注册同步函数及其回滚函数到相应的key中

##### RegisterAsync(key string, asyncFunc func(instance T) error, rollback func(instance T))

注册异步函数及其回滚函数到相应的key中

##### Sync(key string, instance T) error

执行key值对应的同步函数链，当某个函数执行产生error，则中断执行，并返回error

##### Async(key string, instance T)

执行key值对应的异步函数链