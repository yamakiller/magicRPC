package test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func run1() {

}

func run(wait *sync.WaitGroup) {
	go func() {
		defer func() {
			wait.Done()
		}()

		ch := make(chan bool)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Microsecond)
		defer func() {
			cancel()
		}()

		go func() {
			fmt.Fprintf(os.Stderr, "start\n")
			time.Sleep(time.Duration(6 * time.Millisecond))
			fmt.Fprintf(os.Stderr, "end\n")
			ch <- true
		}()

		select {
		case <-ch:
			fmt.Fprintf(os.Stderr, "ch\n")
			return
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "timeout\n")
			fmt.Fprintf(os.Stderr, "exit\n")
			return
		}
	}()
	//wait.Wait()
}

func TestContext(t *testing.T) {
	wait := sync.WaitGroup{}
	wait.Add(1)
	run(&wait)
	/*go func() {
		defer func() {
			wait.Done()
		}()
		ch := make(chan bool)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Microsecond)
		go func() {
			fmt.Fprintf(os.Stderr, "start\n")
			time.Sleep(time.Duration(6 * time.Millisecond))
			fmt.Fprintf(os.Stderr, "end\n")
			ch <- true
		}()

		select {
		case <-ch:
			fmt.Fprintf(os.Stderr, "ch\n")
			cancel()
			return
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "timeout\n")
			cancel()
			fmt.Fprintf(os.Stderr, "exit\n")
			return
		}
	}()*/

	time.Sleep(time.Duration(10 * time.Millisecond))
}
