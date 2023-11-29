package fc

import "sync"

type node struct {
	function func() error
	rollback func()
	prevNode *node
	nextNode *node
}

func newNode() *node {
	return &node{}
}

func (n *node) pushFront(function func() error, rollback func()) {
	newNd := &node{function: n.function, rollback: n.rollback}
	n.function = function
	n.rollback = rollback
	next := n.nextNode
	newNd.prevNode = n
	n.nextNode = newNd
	if next != nil {
		next.prevNode = newNd
		newNd.nextNode = next
	}
}

func (n *node) pushBack(function func() error, rollback func()) {
	end := n.end()
	if end.function == nil && end.rollback == nil {
		end.function = function
		end.rollback = rollback
	} else {
		newNd := &node{function: function, rollback: rollback}
		end.nextNode = newNd
		newNd.prevNode = end
	}
}

func (n *node) end() *node {
	end := n
	for {
		if end.nextNode == nil {
			break
		}
		end = end.nextNode
	}
	return end
}

func (n *node) next() *node {
	return n.nextNode
}

func (n *node) prev() *node {
	return n.prevNode
}

func (n *node) run() error {
	if n.function == nil {
		return nil
	}
	err := n.function()
	if err != nil {
		n.revert()
		return err
	}
	if n.next() == nil {
		return nil
	}
	return n.next().run()
}

func (n *node) revert() {
	if n.rollback != nil {
		n.rollback()
	}
	if n.prev() == nil {
		return
	}
	n.prev().revert()
}

func (n *node) runAsync() []error {
	var wg sync.WaitGroup
	var lock sync.Mutex
	var errs []error
	currNode := n
	wg.Add(1)
	for {
		if currNode.function != nil {
			wg.Add(1)
			go func(function func() error) {
				err := function()
				if err != nil {
					lock.Lock()
					errs = append(errs, err)
					lock.Unlock()
				}
				wg.Done()
			}(currNode.function)
		}
		if currNode.nextNode == nil {
			wg.Done()
			break
		}
		currNode = currNode.nextNode
	}
	wg.Wait()
	if len(errs) > 0 {
		currNode = n
		wg.Add(1)
		for {
			if currNode.rollback != nil {
				wg.Add(1)
				go func(rollback func()) {
					rollback()
					wg.Done()
				}(currNode.rollback)
			}
			if currNode.nextNode == nil {
				wg.Done()
				break
			}
			currNode = currNode.nextNode
		}
		wg.Wait()
	}
	return errs
}

func (n *node) start() {
	var wg sync.WaitGroup
	currNode := n
	wg.Add(1)
	for {
		if currNode.function != nil {
			wg.Add(1)
			go func(function func() error, rollback func()) {
				err := function()
				if err != nil {
					rollback()
				}
				wg.Done()
			}(currNode.function, currNode.rollback)
		}
		if currNode.nextNode == nil {
			wg.Done()
			break
		}
		currNode = currNode.nextNode
	}
	wg.Wait()
}
