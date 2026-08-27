package main

import (
	"reflect"
	"testing"
)

// getTopKWords — вспомогательная функция, которая из строки и K возвращает
// срез из K наиболее частых слов, используя твои функции.
func getTopKWords(input string, k int) []string {
	words := getWorldsString(input)          // получаем слайс слов
	freqMap := getWorldsMap(words)           // считаем частоты
	sortedWords := getSortWorldsMap(freqMap) // сортируем по частоте и лексикографически

	// ограничиваем количество слов
	limit := k
	if limit > len(sortedWords) {
		limit = len(sortedWords)
	}

	// извлекаем только сами слова
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = sortedWords[i].Word
	}
	return result
}

func TestTopKFrequentWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		k        int
		expected []string
	}{
		{
			name:     "обычный случай: K < количества уникальных слов",
			input:    "aa bb cc aa cc cc cc aa ab ac bb",
			k:        3,
			expected: []string{"cc", "aa", "bb"},
		},
		{
			name:     "пустая строка",
			input:    "",
			k:        5,
			expected: []string{},
		},
		{
			name:     "K больше количества уникальных слов",
			input:    "apple banana apple",
			k:        5,
			expected: []string{"apple", "banana"},
		},
		{
			name:     "одинаковая частота — лексикографический порядок",
			input:    "zebra apple apple zebra cat dog",
			k:        4,
			expected: []string{"apple", "zebra", "cat", "dog"},
		},
		{
			name:     "один и тот же сценарий с разными K",
			input:    "one two two three three three",
			k:        2,
			expected: []string{"three", "two"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTopKWords(tt.input, tt.k)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("getTopKWords(%q, %d) = %v, want %v",
					tt.input, tt.k, got, tt.expected)
			}
		})
	}
}

func TestTopKFrequentWordsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		k        int
		expected []string
	}{
		{
			name:     "K = 0",
			input:    "hello world hello",
			k:        0,
			expected: []string{},
		},
		{
			name:     "только одно слово",
			input:    "golang",
			k:        1,
			expected: []string{"golang"},
		},
		{
			name:     "все слова одинаковые",
			input:    "same same same",
			k:        2,
			expected: []string{"same"},
		},
		{
			name:     "слова с пробелами в начале и конце",
			input:    "   a b a   ",
			k:        2,
			expected: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTopKWords(tt.input, tt.k)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("getTopKWords(%q, %d) = %v, want %v",
					tt.input, tt.k, got, tt.expected)
			}
		})
	}
}
