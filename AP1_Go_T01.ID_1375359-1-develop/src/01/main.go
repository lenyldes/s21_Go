package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func isValidOperator(str string) bool {
	return str == "+" || str == "-" || str == "*" || str == "/"
}

func readOperand() float64 {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if scanner.Err() != nil {
			fmt.Println("Invalid input")
			fmt.Println("Repeat input:")
			continue
		}

		tmpInput := scanner.Text()

		value, err := strconv.ParseFloat(tmpInput, 64)
		if err != nil {
			fmt.Println("Invalid input")
			fmt.Println("Repeat input:")
			continue
		}

		return value
	}
}

func readOperation() string {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if scanner.Err() != nil {
			fmt.Println("Invalid input")
			fmt.Println("Repeat input:")
			continue
		}

		operator := scanner.Text()

		if !isValidOperator(operator) {
			fmt.Println("Invalid input")
			fmt.Println("Repeat input:")
			continue
		}

		return operator
	}
}

func readInput() (leftOperand float64, operator string, rightOperand float64) {
	fmt.Println("Input left operand:")
	leftOperand = readOperand()

	fmt.Println("Input operation (+ - * /):")
	operator = readOperation()

	fmt.Println("Input right operand:")
	rightOperand = readOperand()

	return leftOperand, operator, rightOperand
}

func evaluation(leftOperand float64, operator string, rightOperand float64) (result float64, err string) {
	switch operator {
	case "+":
		return leftOperand + rightOperand, err
	case "-":
		return leftOperand - rightOperand, err
	case "*":
		return leftOperand * rightOperand, err
	case "/":
		epsilon := 1e-6
		if math.Abs(rightOperand) <= epsilon {
			return 0, "division by zero"
		}
		return leftOperand / rightOperand, err
	default:
		return 0, "invalid operator"
	}
}

func main() {
	fmt.Println("Hi! This is a simple console calculator in Go")
	fmt.Println("I work with Float64 data type and support the following operations: + - / *")

	leftOperand, operator, rightOperand := readInput()
	result, err := evaluation(leftOperand, operator, rightOperand)

	if err != "" {
		fmt.Println("Calculation error:", err)
	} else {
		fmt.Printf("Result: %.3f\n", result)
	}
}
