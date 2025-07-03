package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	numSorters                         = 10
	maxLength                          = 80
	minSortInterval                    = 3 * time.Second
	intervalRange                      = int64(12 * time.Second)
	sortObservationDelay               = 50 * time.Millisecond
)

type Node struct {
	data string
	next *Node
	mu   sync.RWMutex
}

type LinkedList struct {
	sentinel *Node
}

func NewLinkedList() *LinkedList {
	return &LinkedList{
		sentinel: &Node{},
	}
}

func (l *LinkedList) Add(data string) {
	var nodesToPush []*Node
	runes := []rune(data)
	for len(runes) > 0 {
		var part string
		if len(runes) > maxLength {
			part = string(runes[:maxLength])
			runes = runes[maxLength:]
		} else {
			part = string(runes)
			runes = nil
		}
		nodesToPush = append(nodesToPush, &Node{data: part})
	}

	l.sentinel.mu.Lock()
	defer l.sentinel.mu.Unlock()

	for i := len(nodesToPush) - 1; i >= 0; i-- {
		node := nodesToPush[i]
		node.next = l.sentinel.next
		l.sentinel.next = node
	}
}

func (l *LinkedList) Print() {
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

func (l *LinkedList) sort() (completedSuccessfully bool) {
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
			
			if sortObservationDelay > 0 {
				time.Sleep(sortObservationDelay)
			}
			
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

func sorter(ctx context.Context, id int, list *LinkedList, sleepTime time.Duration) {
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
				
				if completed := list.sort(); completed {
					break
				}
			}
		case <-ctx.Done():
			fmt.Printf("[Sorter %d] Shutdown...\n", id)
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	list := NewLinkedList()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	fmt.Printf("Starting %d sorter goroutines...\n", numSorters)
	for i := 0; i < numSorters; i++ {
		wg.Add(1)
		randomOffset := time.Duration(rand.Int63n(intervalRange))
		sleepDuration := minSortInterval + randomOffset

		go func (i int, sleepDuration time.Duration) {
			defer wg.Done()
			sorter(ctx, i, list, sleepDuration)
		}(i, sleepDuration)
	}
	fmt.Println("The sorters have been launched.")
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\nThe program is running, write something (press enter to display the list, Ctrl+D to exit)")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nEOF received (Ctrl+D), starting shutdown...")
				break
			}
			fmt.Println("Error reading input:", err)
			break
		}

		line := strings.TrimSpace(input)
		if line == "" {
			list.Print()
		} else {
			list.Add(line)
		}
	}

	fmt.Println("Sending a cancel signal to sorters...")
	cancel()

	fmt.Println("Wait for all sorter goroutines to complete...")
	wg.Wait()
	
	fmt.Println("All sorters have completed their work..")
	fmt.Println("\nFinal state of the list:")
	list.Print()
	fmt.Println("Program completed.")
}
