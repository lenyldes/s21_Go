package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type pationVisit struct {
	doctorSpec string    // специализация врача
	visitDate  time.Time // дата посещения
}

func parseDateFromString(timeStr string) (time.Time, error) {
	layout := "2006-01-02" // YYYY-MM-DD

	parsedDate, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDate, nil
}

func readInput() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() == true {
		return strings.TrimSpace(scanner.Text()), nil
	}

	if scanner.Err() == nil {
		return "", io.EOF
	}

	return "", scanner.Err()
}

func readName() string {
	fmt.Println("Введите ФИО (например: Иванов Иван Иванович):")
	for {
		name, errNameStr := readInput()
		if errNameStr != nil {
			fmt.Println("Ошибка при чтении ФИО:", errNameStr)
			fmt.Println("Попробуйте снова:")
			continue
		}
		return name
	}
}

func readDoctor() string {
	fmt.Println("Введи специализацию (например: ортопед):")
	for {
		dockor, errDockorStr := readInput()
		if errDockorStr != nil {
			fmt.Println("Ошибка при чтении специализации врача:", errDockorStr)
			fmt.Println("Попробуйте снова:")
			continue
		}
		return dockor
	}
}

func readDate() time.Time {
	fmt.Println("Введи дату в формате YYYY-MM-YY (например: 2024-04-13):")
	for {
		dateStr, errDateStr := readInput()
		if errDateStr != nil {
			fmt.Println("Ошибка при чтении даты:", errDateStr)
			fmt.Println("Попробуйте снова:")
			continue
		}

		date, errDate := parseDateFromString(dateStr)
		if errDate != nil {
			fmt.Println("Ошибка при парсинге даты:", errDate)
			fmt.Println("Попробуйте снова:")
			continue
		}

		return date
	}
}

func save(cardIndex map[string][]pationVisit) {
	nameStr := readName()
	doctorSpecRead := readDoctor()
	dateTime := readDate()

	val, _ := cardIndex[nameStr]
	cardIndex[nameStr] = append(val, pationVisit{doctorSpec: doctorSpecRead, visitDate: dateTime})

	fmt.Printf("Запись: (%s - %s - %s) успешно сохранена в картотеке\n", nameStr, doctorSpecRead, dateTime.Format("2006-01-02"))
}

func getHistory(cardIndex map[string][]pationVisit) {
	patientsName := readName()

	for name, visit := range cardIndex {
		if name == patientsName {
			for _, pationVisit := range visit {
				fmt.Printf("%s %s\n", pationVisit.doctorSpec, pationVisit.visitDate.Format("2006-01-02"))
			}
		}
	}
}

type PatientNotFoundError struct{}

func (e PatientNotFoundError) Error() string {
	return "patient not found"
}

func getLastVisit(cardIndex map[string][]pationVisit) error {
	patientsName := readName()

	if _, exists := cardIndex[patientsName]; !exists {
		return PatientNotFoundError{}
	}

	doctorSpec := readDoctor()

	var lastDate time.Time

	var specCount int

	for name, visit := range cardIndex {
		if name == patientsName {
			for _, pationVisit := range visit {
				if pationVisit.doctorSpec == doctorSpec && lastDate.Before(pationVisit.visitDate) {
					lastDate = pationVisit.visitDate
					specCount++
				}
			}

		}
	}

	if specCount == 0 {
		fmt.Println("Пациент .. не посещал врача ..")
	} else {
		fmt.Println(lastDate.Format("2006-01-02"))
	}

	return nil
}

func main() {
	cardIndex := map[string][]pationVisit{} // имя пациента - pationVisit{специализация врача, дата посещения}

	for {
		fmt.Println("Введи команду: Save | GetHistory | GetLastVisit | q")
		inputStr, errStr := readInput()
		if errStr != nil {
			fmt.Println("Invalid input:", errStr)
			return
		}

		switch inputStr {
		case "Save":
			save(cardIndex)
		case "GetHistory":
			getHistory(cardIndex)
		case "GetLastVisit":
			err := getLastVisit(cardIndex)
			if errors.Is(err, PatientNotFoundError{}) {
				fmt.Println(err)
			}
		case "q":
			return

		default:
			fmt.Println("Неизвестная команда:", inputStr)
			fmt.Println("Ожидаю команды: Save | GetHistory | GetLastVisit | q")
		}
	}
}
