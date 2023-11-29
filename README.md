# FUNCTION CHAIN

## 介绍

FUNCTION CHAIN用于定义一系列可同步或异步执行的函数链，不同的函数执行链根据注册时的key值区分，每个函数都有其相应的回滚函数。

## 回滚

当函数链为同步函数链时，如果执行某个函数产生error，那么会中断此函数链，并且已经执行过的所有函数的回滚函数都会被执行。

当函数链为异步函数链时，执行某个函数产生error并不会中断函数链后续执行，当函数链中的所有函数都执行完毕后产生了error的函数的回滚函数会被依次执行。

例如定义了以下函数链（括号中为函数对应的回滚函数）:

A(Arb) - B(Brb) - C(Crb) - D(Drb) - E(Erb)

函数链的正常执行流程：A - B - C - D - E

当函数链为同步函数链时，执行C产生error，那么执行流程为：A - B - C - Crb - Brb - Arb

当函数链为异步函数链时，执行C产生error，那么执行流程为：A - B - C - D - E - Crb 当函数链为异步函数链时，执行C、E产生error，那么执行流程为：A - B - C - D - E - Crb - Erb

## 常用函数

##### Add(function func() error, rollback func())

添加同步函数及其回滚函数

##### AddAsync(function func() error, rollback func())

添加异步函数及其回滚函数

##### Start() error

异步函数链和同步函数链同时执行

##### Run() error

先执行异步函数链，执行完后执行同步函数链