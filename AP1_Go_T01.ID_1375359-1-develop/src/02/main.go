package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Word struct {
	Word  string
	Count int
}

func readInput() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() == true {
		return scanner.Text(), nil
	}

	if scanner.Err() == nil {
		return "", io.EOF
	}

	return "", scanner.Err()
}

func getWorldsString(str string) []string {
	return strings.Fields(str)
}

func getWorldsMap(str []string) map[string]int {
	resultWorldsMap := map[string]int{}

	for _, value := range str {
		resultWorldsMap[value]++
	}

	return resultWorldsMap
}

func getSortWorldsMap(worldsMap map[string]int) []Word {
	worlds := []Word{}
	for word, count := range worldsMap {
		worlds = append(worlds, Word{Word: word, Count: count})
	}

	sort.Slice(worlds, func(i, j int) bool {
		if worlds[i].Count == worlds[j].Count {
			return worlds[i].Word < worlds[j].Word
		}

		return worlds[i].Count > worlds[j].Count
	})

	return worlds
}

func printKFirstWorlds(k int, worlds []Word) {
	firstKWorlds := worlds[:min(k, len(worlds))]

	for _, value := range firstKWorlds {
		fmt.Printf("%s ", value.Word)
	}

	fmt.Print("\n")
}

func main() {
	inputStr, errStr := readInput()
	if errStr != nil {
		fmt.Println("Invalid input:", errStr)
		return
	}

	inputKStr, errKStr := readInput()
	if errKStr != nil {
		fmt.Println("Invalid input:", errKStr)
		return
	}

	inputKInt, errKInt := strconv.Atoi(inputKStr)
	if errKInt != nil {
		fmt.Println("Invalid input:", errKInt)
		return
	}

	printKFirstWorlds(inputKInt, getSortWorldsMap(getWorldsMap(getWorldsString(inputStr))))
}
