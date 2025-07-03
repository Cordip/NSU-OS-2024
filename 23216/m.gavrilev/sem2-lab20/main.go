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
	linkedList "sem2-lab20/list"
)

const (
	numSorters                         = 10
	minSortInterval                    = 3 * time.Second
	intervalRange                      = 12 * time.Second
)

func main() {
	rand.Seed(time.Now().UnixNano())

	list := linkedList.NewLinkedList()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	fmt.Printf("Starting %d sorter goroutines...\n", numSorters)
	for i := 0; i < numSorters; i++ {
		wg.Add(1)
		randomOffset := time.Duration(rand.Intn(int(intervalRange)))
		sleepDuration := minSortInterval + randomOffset

		go func (i int, sleepDuration time.Duration) {
			defer wg.Done()
			linkedList.Sorter(ctx, i, list, sleepDuration)
		}(i, sleepDuration)
	}
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
