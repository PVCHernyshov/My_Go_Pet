package main

import (
	"reflect"
	"testing"
)

func TestClearData(t *testing.T) {
	tests := []struct {
		Name         string
		inputData    string
		expectedData string
	}{ // не буквы меняются на пробелы = появляются лишние пробелвы
		{"обычная строка", "Простая строка. И знаки препинания!", "простая строка  и знаки препинания "},
		{"строка с цифрами", "строка12 456строка", "строка      строка"},
		{"разный регистр", "РазНыЙ РЕГИстр", "разный регистр"},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result := clearData(tt.inputData)

			if tt.expectedData != result {
				t.Errorf("ожидалось \"%s\" получено \"%s\"", tt.expectedData, result)
			}
		})
	}
}

func TestSortingOfWords(t *testing.T) {
	case1 := map[string]struct{}{"просто": {}, "пустая": {}, "строка": {}}
	case2 := map[string]struct{}{"строка": {}}
	case3 := map[string]struct{}{}

	tests := []struct {
		Name         string
		inputData    string
		expectedData map[string]struct{}
	}{
		{"простая строка с лишними пробелами", "просто    пустая     строка", case1},
		{"строка с дубликатами", "строка строка", case2},
		{"пустая строка", "", case3},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result := sortingOfWords(tt.inputData)

			if !reflect.DeepEqual(result, tt.expectedData) {
				t.Errorf("ожидалось \"%v\" получено \"%v\"", tt.expectedData, result)
			}
		})
	}
}
