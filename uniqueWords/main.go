package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const RETURN string = "return"
const CONTINUE string = "continue"

func main() {

	for true {

		fmt.Print("\nВведите текст или имя файла(для завершения ввода нажмите дважды Ctrl+d): ")
		resUsIn, errUsIn := userInput() // ввод от пользователя
		if errUsIn != nil {
			fmt.Println("Не удалось прочитать ввод. Ошибка:", errUsIn, "\n", "завершение программы.")
			return
		}

		var errReadFile error
		var nameFile string
		if textOrFile(resUsIn) { // определяем текст это или файл и передаём в обработку
			nameFile = resUsIn
			resUsIn, errReadFile = readFile(resUsIn)
		}

		if errReadFile != nil {

			respons := handleError(errReadFile)

			if respons == RETURN {
				return
			}
			if respons == CONTINUE {
				continue
			}

		}

		resultClearData := clearData(resUsIn)

		if strings.TrimSpace(resultClearData) == "" {

			var emptyLine error = errors.New("Пустая строка")

			respons := handleError(emptyLine)

			if respons == RETURN {
				return
			}
			if respons == CONTINUE {
				continue
			}
		}

		resultSortWords := sortingOfWords(resultClearData)

		if nameFile == "" {
			nameFile = "nameless.txt"
		}

		errWriteFile := writeFile(nameFile, resUsIn, resultSortWords)

		if errWriteFile != nil {

			respons := handleError(errWriteFile)

			if respons == RETURN {
				return
			}
			if respons == CONTINUE {
				continue
			}

		}

		menuDisplays := true
		for true {
			if menuDisplays {
				fmt.Println("\nРезультат обработки сохраннён в файле", nameFile, "Вывести на экран - п/p(принт/print).")
			}
			fmt.Println("Следующая сортировка с/n(следующая/next).")
			fmt.Println("Выход - Ctrl+d: ")

			resUsIn, errUsIn = userInput()
			if errUsIn != nil {
				fmt.Println("Не удалось прочитать ввод. Ошибка:", errUsIn, "\n", "завершение программы.")
				return
			}

			if strings.ToLower(resUsIn) == "п" || strings.ToLower(resUsIn) == "p" {
				menuDisplays = false
				for str := range resultSortWords {
					fmt.Println(str)
				}
			} else if strings.ToLower(resUsIn) == "с" || strings.ToLower(resUsIn) == "n" {
				fmt.Println()
				break
			} else {
				fmt.Println()
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

	extension := []string{".txt", ".md", ".csv", ".log"}

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
		return "", err
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
		if !unicode.IsLetter(r) {
			return ' '
		}
		return unicode.ToLower(r)
	}, s)
	return str
}

func handleError(e error) string {
	cycle := true
	for true {
		if cycle {
			fmt.Println("\nПроизошла ошибка ", e)
		}
		fmt.Print("Повторить ввод/repeat enter (п/r + (Ctrl+d)) выход/exit (Ctrl+d): ")
		resUsIn, errUsIn := userInput()

		if errUsIn != nil {
			fmt.Println("\nНе удалось прочитать ввод. Ошибка:", errUsIn, "\n", "завершение программы.")
			return RETURN
		}

		if strings.ToLower(resUsIn) == "п" || strings.ToLower(resUsIn) == "r" {
			return CONTINUE
		}

		if resUsIn != "" {
			cycle = false
			fmt.Println("\nОшибка ввода")
			continue
		}

		if resUsIn == "" {
			fmt.Println("\nДо новых встреч!")
			return RETURN
		}
	}

	return ""
}
