package main

import (
	"errors"
	"fmt"
)

type qdata struct {
	key  string
	data interface{}
}

type llnode struct {
	data qdata
	next *llnode
}

type queue struct {
	head *llnode
	tail *llnode
}

type LRUCache struct {
	/* cache struct
	cache struct*/
	Q              queue
	capacity, size int
	prevnodeLookup map[string]*llnode
}

func (q *queue) push(data qdata) *llnode {
	// push to the tail
	var prevnode *llnode
	node := llnode{data, nil}
	if q.head == nil {
		q.head, q.tail = &node, &node
	} else {
		q.tail.next = &node
		prevnode = q.tail
		q.tail = q.tail.next
	}
	return prevnode
}

func (q *queue) pop() (*llnode, error) {
	// pop from the head
	if q.head == nil {
		return nil, errors.New("Empty queue")
	}
	node := q.head
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	return node, nil
}

func (q *queue) movetobottom(prev *llnode) error {
	var node *llnode
	if prev == nil {
		node = q.head
	} else {
		node = prev.next
	}
	if node == nil {
		return errors.New("Couldn't move node to the top")
	}
	if node == q.tail {
		return nil
	}
	if node != q.head {
		prev.next = node.next
	} else {
		q.head = node.next
	}
	q.tail.next = node
	q.tail = node
	node.next = nil
	return nil
}

func (head *llnode) iterate() chan interface{} {
	c := make(chan interface{})
	go func() {
		p := head
		for p != nil {
			c <- p.data
			p = p.next
		}
		close(c)
	}()
	return c
}

func (q *queue) iterate() chan interface{} {
	return q.head.iterate()
}

func (q *queue) String() string {
	repr := ""
	for data := range q.iterate() {
		repr += fmt.Sprint(data, " -> ")
	}
	return repr + "nil"
}

func (cache *LRUCache) init(capacity int) {
	cache.capacity = capacity
	cache.prevnodeLookup = make(map[string]*llnode)
}

func (cache *LRUCache) read(key string) (interface{}, error) {
	if cache.prevnodeLookup == nil {
		return nil, errors.New("cache not initialsed")
	}
	prevnode, ok := cache.prevnodeLookup[key]
	if !ok {
		return nil, nil
	}
	node := prevnode.next
	cache.prevnodeLookup[key] = nil
	return node.data, nil
}

func (cache *LRUCache) touch(key string) error {
	prevnode, ok := cache.prevnodeLookup[key]
	if !ok {
		return errors.New("No data for key " + key)
	}
	if cache.Q.tail == cache.Q.head {
		// only one node
		cache.prevnodeLookup[key] = nil
	} else {
		// more than one node,
		// this node will become new tail
		// so the current tail should be it's prevnode
		cache.prevnodeLookup[key] = cache.Q.tail
	}
	cache.Q.movetobottom(prevnode)
	// Whatever the head be, it's prev is nil
	cache.prevnodeLookup[cache.Q.head.data.key] = nil
	return nil
}

func (cache *LRUCache) write(key string, data interface{}) (qdata, error) {
	var evicted qdata
	if cache.prevnodeLookup == nil {
		return evicted, errors.New("cache not initialised")
	}
	prevnode, ok := cache.prevnodeLookup[key]
	if !ok {
		prevnode = cache.Q.push(qdata{key, data})
		cache.prevnodeLookup[key] = prevnode
		cache.size++
	}
	cache.touch(key)
	if cache.size > cache.capacity {
		node, _ := cache.Q.pop()
		delete(cache.prevnodeLookup, node.data.key)
		cache.size--
		evicted = node.data
	}
	return evicted, nil
}
