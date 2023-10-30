package fc

import (
	"errors"
	"sync"
)

type Chain struct {
	syncNode  *node
	asyncNode *node
}

func New() *Chain {
	return &Chain{
		syncNode:  newNode(),
		asyncNode: newNode(),
	}
}

func (chain *Chain) Add(function func() error, rollback func()) *Chain {
	chain.syncNode.pushBack(function, rollback)
	return chain
}

func (chain *Chain) AddAsync(function func() error, rollback func()) *Chain {
	chain.asyncNode.pushBack(function, rollback)
	return chain
}

func (chain *Chain) Start() error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		chain.asyncNode.start()
		wg.Done()
	}()
	err := chain.syncNode.run()
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (chain *Chain) Run() error {
	errs := chain.asyncNode.runAsync()
	err := chain.syncNode.run()
	if err != nil {
		return err
	} else if len(errs) > 0 {
		chain.syncNode.end().revert()
		return errors.Join(errs...)
	} else {
		return nil
	}
}
