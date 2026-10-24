package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func main() {
	flag := true

	for flag {

		fmt.Print("Введите текст или имя файла(для завершения ввода нажмите Ctrl+d): ")
		resUsIn, errUsIn := userInput() // принимаем данные от пользователя

		if errUsIn != nil {
			fmt.Println(errUsIn)
			return
		}

		var errReadFile error
		var nameFile string
		if textOrFile(resUsIn) { // определяем текст это или файл и передаём в обработку
			nameFile = resUsIn
			resUsIn, errReadFile = readFile(resUsIn)
		}

		if errReadFile != nil {
			fmt.Println(errReadFile)
			return
		}

		resultClearData := clearData(resUsIn)

		if resultClearData == "" {
			fmt.Println("Получена пустая строка")
			return
		}

		resultSortWords := sortingOfWords(resultClearData)

		if nameFile == "" {
			nameFile = "nameless.txt"
		}
		errWriteFile := writeFile(nameFile, resUsIn, resultSortWords)
		if errWriteFile != nil {
			fmt.Println(errWriteFile)
			return
		}

		menu := true
		for menu {
			fmt.Println("\nУникальные слова отделены.")
			fmt.Println("Результат сохраннён в файле", nameFile, "Вывести на экран - п/p(принт/print).")
			fmt.Println("Следующая сортировка с/n(следующая/next).")
			fmt.Println("Выход - любой дрйгой символ.")
			fmt.Print(": ")

			resUsIn, errUsIn = userInput() // принимаем данные от пользователя

			if errUsIn != nil {
				fmt.Println(errUsIn)
				return
			}

			if strings.ToLower(resUsIn) == "п" || strings.ToLower(resUsIn) == "p" {
				for str := range resultSortWords {
					fmt.Println(str)
				}
			} else if strings.ToLower(resUsIn) == "с" || strings.ToLower(resUsIn) == "n" {
				fmt.Println()
				menu = false
				continue
			} else {
				flag, menu = false, false
				return
			}
		}

	}
}

func sortingOfWords(s string) map[string]struct{} {

	list := map[string]struct{}{}
	sliceStr := strings.Fields(s)
	for _, key := range sliceStr {
		list[key] = struct{}{}
	}
	return list
}

func userInput() (string, error) {
	return readFrom(os.Stdin)
}

func textOrFile(s string) bool {

	extension := []string{".txt", ".md", ".cvc", ".log"}

	for _, fileExtension := range extension {
		if strings.HasSuffix(s, fileExtension) {
			return true
		}
	}
	return false
}

func readFile(f string) (string, error) {

	file, err := os.Open(f)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	defer file.Close()

	return readFrom(file)
}

func readFrom(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	buffer := strings.Builder{}
	for scanner.Scan() {
		str := scanner.Text()
		if str == "<--->" {
			break
		}
		if len(str) != 0 {
			buffer.WriteString(str)
			buffer.WriteString(" ")
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	inputData := strings.TrimSpace(buffer.String())

	return inputData, nil
}

func writeFile(f, ui string, m map[string]struct{}) error {
	file, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, ui)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(file, "<--->")
	if err != nil {
		return err
	}

	for vol := range m {
		_, err := fmt.Fprintln(file, vol)
		if err != nil {
			return err
		}

	}

	return nil
}

func clearData(s string) string {
	str := strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return ' '
		}
		return unicode.ToLower(r)
	}, s)
	return str
}
