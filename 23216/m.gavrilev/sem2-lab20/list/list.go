package list

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type node struct {
	data string
	next *node
	mu   sync.RWMutex
}

type linkedList struct {
	sentinel *node
}

func NewLinkedList() *linkedList {
	return &linkedList{
		sentinel: &node{},
	}
}

func (l *linkedList) Add(data string) {
	l.sentinel.mu.Lock()
	defer l.sentinel.mu.Unlock()

	newNode := &node{
		data: data,
		next: l.sentinel.next,
	}
	l.sentinel.next = newNode
}

func (l *linkedList) Print() {
	fmt.Println("--- Current state of the list ---")
	l.sentinel.mu.RLock()
	current := l.sentinel.next
	l.sentinel.mu.RUnlock()

	if current == nil {
		fmt.Println("The list is empty.")
		fmt.Println("--------------------------------")
		return
	}
	for i := 0; current != nil; i++ {
		current.mu.RLock()
		nodeData := current.data
		nextNode := current.next
		current.mu.RUnlock()
		fmt.Printf("[%d] -> %s\n", i, nodeData)
		current = nextNode
	}
	fmt.Println("--------------------------------")
}

func (l *linkedList) Sort() (completedSuccessfully bool) {
	prev := l.sentinel

	prev.mu.RLock()
	current := prev.next
	prev.mu.RUnlock()

	for current != nil && current.next != nil {
		next := current.next

		lockedPrev, lockedCurrent, lockedNext := prev, current, next

		lockedPrev.mu.Lock()
		lockedCurrent.mu.Lock()
		lockedNext.mu.Lock()

		contextLost := (lockedPrev.next != lockedCurrent) || (lockedCurrent.next != lockedNext)
		if contextLost {
			lockedNext.mu.Unlock()
			lockedCurrent.mu.Unlock()
			lockedPrev.mu.Unlock()
			return false
		}

		if lockedCurrent.data > lockedNext.data {
			lockedCurrent.next = lockedNext.next
			lockedNext.next = lockedCurrent
			lockedPrev.next = lockedNext

			prev = lockedNext
		} else {
			prev = lockedCurrent
			current = lockedNext
		}

		lockedNext.mu.Unlock()
		lockedCurrent.mu.Unlock()
		lockedPrev.mu.Unlock()
	}
	return true
}

func Sorter(ctx context.Context, id int, list *linkedList, sleepTime time.Duration) {
	ticker := time.NewTicker(sleepTime)
	defer ticker.Stop()
	fmt.Printf("[Sorter %d] Launched and sleeping for %d (seconds).\n", id, sleepTime/time.Second)

	for {
		select {
		case <-ticker.C:
			for {
				if ctx.Err() != nil {
					break
				}

				if completed := list.Sort(); completed {
					break
				}
			}
		case <-ctx.Done():
			fmt.Printf("[Sorter %d] Shutdown...\n", id)
			return
		}
	}
}
