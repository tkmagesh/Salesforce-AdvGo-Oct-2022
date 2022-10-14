package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}
	rootCtx := context.Background()
	valCtx := context.WithValue(rootCtx, "k1", "v1")
	cancelCtx, cancel := context.WithCancel(valCtx)
	wg.Add(1)
	go printDot(wg, cancelCtx)
	go func() {
		fmt.Println("Hit ENTER to stop...")
		fmt.Scanln()
		cancel()
	}()
	wg.Wait()
}

func printDot(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	wg.Add(1)
	overCtx := context.WithValue(ctx, "k1", "new-v1")
	val2Ctx := context.WithValue(overCtx, "k2", "v2")
	timeoutCtx, cancel := context.WithTimeout(val2Ctx, 5*time.Second)
	defer cancel()
	go printTick(wg, timeoutCtx)

LOOP:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stoping dots")
			break LOOP
		default:
			fmt.Print(".")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func printTick(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	fmt.Println("Data from context : ", ctx.Value("k1"))
	fmt.Println("Data from context : ", ctx.Value("k2"))
LOOP:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stoping ticks")
			break LOOP
		default:
			fmt.Print("Tick")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
