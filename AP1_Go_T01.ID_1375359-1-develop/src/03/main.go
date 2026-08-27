package intersection

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

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

func getIntSliceFromString(str string) []int {
	strSlice := strings.Fields(str)

	resultIntSlice := []int{}

	for _, value := range strSlice {
		tmpInt, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Invalid input")
			os.Exit(1)
		}
		resultIntSlice = append(resultIntSlice, tmpInt)
	}

	return resultIntSlice
}

func findIntersectionFastWithMap(a []int, b []int) []int {
	resultIntSlice := []int{}
	mapFromB := map[int]int{} // value/count

	for _, value := range b {
		mapFromB[value]++
	}

	for _, value := range a {
		if mapFromB[value] > 0 {
			resultIntSlice = append(resultIntSlice, value)
			mapFromB[value]--
		}
	}

	return resultIntSlice
}

func printResult(a []int) {
	if len(a) == 0 {
		println("Empty intersection")
		return
	}

	for i := 0; i < len(a); i++ {
		fmt.Printf("%d ", a[i])
	}

	fmt.Println()
}

func main() {
	inputStr1, errStr1 := readInput()
	if errStr1 != nil {
		fmt.Println("Invalid input:", errStr1)
		return
	}

	intSlice1 := getIntSliceFromString(inputStr1)

	inputStr2, errStr2 := readInput()
	if errStr2 != nil {
		fmt.Println("Invalid input:", errStr2)
		return
	}

	intSlice2 := getIntSliceFromString(inputStr2)

	printResult(findIntersectionFastWithMap(intSlice1, intSlice2))
}
