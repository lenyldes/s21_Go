package main

import (
	"flag"
	"fmt"
)

func parseParams() (K int, N int, err error) {
	flag.IntVar(&K, "K", 0, "Где K: Числа от K до N")
	flag.IntVar(&N, "N", 0, "Где N: Числа от K до N")
	flag.Parse()

	if flag.NFlag() != 2 {
		return 0, 0, fmt.Errorf("Ошибка: нужно указать оба флага -K и -N")
	}

	if N < K {
		return 0, 0, fmt.Errorf("Ошибка: Число N должно быть больше или равно K (N>=K)")
	}

	return K, N, nil
}

func generate(K int, N int) <-chan int {
	ch := make(chan int)

	go func() {
		for i := K; i <= N; i++ {
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func squere(in_ch <-chan int) <-chan int {
	out_ch := make(chan int)

	go func() {
		for val := range in_ch {
			out_ch <- val * val
		}
		close(out_ch)
	}()

	return out_ch
}

func main() {
	K, N, err := parseParams()
	if err != nil {
		fmt.Println("--------")
		fmt.Println(err)
		flag.Usage()
		return
	}

	ch1 := generate(K, N)
	ch2 := squere(ch1)

	for val := range ch2 {
		fmt.Println(val)
	}
}
