package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func parseParams() (v uint, err error) {
	var K int
	flag.IntVar(&K, "K", 0, "K задаёт шаг тикера в секундах и имеет тип uint.")
	flag.Parse()

	if flag.NFlag() != 1 {
		return v, fmt.Errorf("Ошибка: нужно указать один флаг -K {шаг тика в секундах}")
	}

	if K <= 0 {
		return v, fmt.Errorf("Ошибка: Число K должно быть больше нуля")
	}

	return uint(K), nil
}

func ticker(ctx context.Context, wg *sync.WaitGroup, K uint) {
	sleepTime := time.Second * time.Duration(K)
	var currTime uint = K
	tick := 1

	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				time.Sleep(sleepTime)
				fmt.Printf("Tick %d since %d\n", tick, currTime)
				currTime += K
				tick++
			}

		}
	}()
}

func main() {
	K, err := parseParams()
	if err != nil {
		fmt.Println("--------")
		fmt.Println(err)
		flag.Usage()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	ticker(ctx, &wg, K)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	cancel()
	wg.Wait()

	fmt.Println("Termination")
}
