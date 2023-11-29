package fc

import (
	"errors"
	"testing"
)

func TestNewNode(t *testing.T) {
	n := newNode()
	n.pushFront(func() error {
		t.Log("first front")
		return nil
	}, func() {
		t.Log("first front rollback")
	})
	n.pushBack(func() error {
		t.Log("first")
		return nil
	}, func() {
		t.Log("first rollback")
	})
	n.pushBack(func() error {
		t.Log("second")
		return nil
	}, func() {
		t.Log("second rollback")
	})
	n.pushBack(func() error {
		t.Log("third")
		return errors.New("panic")
	}, func() {
		t.Log("third rollback")
	})
	t.Log(n.run())
}

func TestChain_Run(t *testing.T) {
	err := data(t).Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestChain_Start(t *testing.T) {
	err := data(t).Start()
	if err != nil {
		t.Fatal(err)
	}
}

func data(t *testing.T) *Chain {
	return New().Add(func() error {
		t.Log("1")
		return nil
	}, func() {
		t.Log("1 cb")
	}).Add(func() error {
		t.Log("2")
		return nil
	}, func() {
		t.Log("2 cb")
	}).Add(func() error {
		t.Log("3")
		return nil
	}, func() {
		t.Log("3 cb")
	}).AddAsync(func() error {
		t.Log("1 async")
		return nil
	}, func() {
		t.Log("1 async cb")
	}).AddAsync(func() error {
		t.Log("2 async")
		return nil
	}, func() {
		t.Log("2 async cb")
	}).AddAsync(func() error {
		t.Log("3 async")
		return nil
	}, func() {
		t.Log("3 async cb")
	}).AddAsync(func() error {
		t.Log("4 async")
		return errors.New("4 async panic")
	}, func() {
		t.Log("4 async cb")
	}).AddAsync(func() error {
		t.Log("5 async")
		return errors.New("5 async panic")
	}, func() {
		t.Log("5 async cb")
	})
}
