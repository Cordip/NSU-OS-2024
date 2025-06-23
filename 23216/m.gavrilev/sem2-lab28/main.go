package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sem2-lab28/myhttp"
	"sem2-lab28/paginator"
	"sync"
	"log"
)

const (
	networkBufferSize   = 4096
	errorResponseLimit  = 512
	dataChannelBuffer   = 100
	scannerMaxCapacity  = 1024 * 1024
)

func main() {
	log.SetPrefix("error: ")
	log.SetFlags(0)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <URL>\n", os.Args[0])
		os.Exit(1)
	}
	url := os.Args[1]

	fmt.Fprintf(os.Stderr, "[Fetching %s ...]\n", url)
	resp, err := myhttp.Get(url)
	if err != nil {
		log.Fatalf("failed to fetch URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != myhttp.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(io.LimitReader(resp.Body, errorResponseLimit))
		log.Fatalf("received non-200 status code: %d %s\nResponse body (partial): %s",
			resp.StatusCode, resp.Status, string(bodyBytes))
	}

	dataChan := make(chan []byte, dataChannelBuffer)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(){
		defer wg.Done()
		networkReader(ctx, resp.Body, dataChan)
	}()
	go func(){
		defer wg.Done()
		userInteractor(ctx, errChan, dataChan)
	}()

	if finalErr := <-errChan; finalErr != nil {
		cancel()
		if errors.Is(finalErr, pager.ErrInterrupted) {
			fmt.Fprintf(os.Stderr, "\n[Interrupted by user]\n")
		} else {
			log.Printf("critical error during interaction: %v", finalErr)
		}
	}

	wg.Wait()
	fmt.Fprintf(os.Stderr, "\n[Done]\n")
}

func networkReader(ctx context.Context, body io.ReadCloser, dataChan chan<- []byte) {
	defer close(dataChan)

	buffer := make([]byte, networkBufferSize)
	for {
		n, err := body.Read(buffer)
		if n > 0 {
			dataCopy := make([]byte, n)
			copy(dataCopy, buffer[:n])
			select {
			case dataChan <- dataCopy:
			case <-ctx.Done():
				return
			}
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("reading response body error: %v", err)
			}
			break
		}
	}
}

func userInteractor(ctx context.Context, errChan chan<- error, dataChan <-chan []byte) {
	defer close(errChan)
	p := pager.New()
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()

	go func() {
		defer pipeWriter.Close()
		for dataChunk := range dataChan {
			if _, writeErr := pipeWriter.Write(dataChunk); writeErr != nil {
				return
			}
		}
	}()

	scanner := bufio.NewScanner(pipeReader)
	buf := make([]byte, bufio.MaxScanTokenSize)
	scanner.Buffer(buf, scannerMaxCapacity)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err := p.Println(scanner.Text()); err != nil {
			errChan <- err
			return
		}
	}

	if err := scanner.Err(); err != nil && !errors.Is(err, io.ErrClosedPipe) {
		errChan <- fmt.Errorf("scanner error: %w", err)
	}
}