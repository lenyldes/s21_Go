package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func parseParams() (N int, M int, err error) {
	flag.IntVar(&N, "N", 0, "Количество горутин")
	flag.IntVar(&M, "M", 0, "Время сна, равное рандомному времени до M миллисекунд")
	flag.Parse()

	if flag.NFlag() != 2 {
		return 0, 0, fmt.Errorf("Ошибка: нужно указать оба флага -N и -M")
	}

	if N <= 0 || M <= 0 {
		return 0, 0, fmt.Errorf("Ошибка: N и M должно быть >= 1 (шт)")
	}

	return N, M, nil
}

func main() {
	N, M, err := parseParams()
	if err != nil {
		fmt.Println("--------")
		fmt.Println(err)
		flag.Usage()
		return
	}

	gorutins := make([]int, N)

	var wg sync.WaitGroup
	for i := range gorutins {
		mc := rand.Intn(M + 1)
		wg.Go(func() {
			time.Sleep(time.Duration(mc) * time.Millisecond)
		})
		gorutins[i] = mc
	}
	wg.Wait()

	sort.Slice(gorutins, func(i, j int) bool {
		return gorutins[i] > gorutins[j]
	})

	for i, v := range gorutins {
		fmt.Println(i, "\t", v)
	}
}
