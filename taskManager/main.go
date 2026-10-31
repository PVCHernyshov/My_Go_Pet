package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Task struct {
	ID     int       `json:"id"`
	Date   time.Time `json:"date"`
	Task   string    `json:"task"`
	Status bool      `json:"status"`
}

func main() {

	fmt.Println("\n\n", heading(time.Now()))

	var tasks []Task
	err := readFile(&tasks)

	if err != nil {
		if os.IsNotExist(err) {
			printError()
		} else {
			fmt.Println(err)
			return
		}
	} else {
		parsingUsIn("список", &tasks)
	}

	for true {
		menu()

		resUsIn, errUsIn := userInput()

		if errUsIn != nil {
			fmt.Println("Не удалось прочитать ввод. Ошибка:", errUsIn, "\n", "завершение программы.")
			return
		}

		err = parsingUsIn(resUsIn, &tasks)
		if err != nil {
			return
		}
	}

}

func heading(t time.Time) string {
	days := map[string]string{"Monday": "Понедельник", "Tuesday": "Вторник", "Wednesday": "Среда",
		"Thursday": "Четверг", "Friday": "Пятница", "Saturday": "Суббота", "Sunday": "Воскресенье"}
	months := map[string]string{"January": "Января", "February": "Февраля", "March": "Марта",
		"April": "Апреля", "May": "Мая", "June": "Июня", "July": "Июля", "August": "Августа",
		"September": "Сентября", "October": "Октября", "November": "Ноября", "December": "Декабря"}

	res1 := t

	res2 := res1.Format("Monday")
	res3 := res1.Format("January")

	str := strings.Builder{}
	str.WriteString("==== Сегодня: ")
	str.WriteString(days[res2])
	str.WriteString(". ")
	str.WriteString(res1.Format("02"))
	str.WriteString(" ")
	str.WriteString(months[res3])
	str.WriteString(" ")
	str.WriteString(res1.Format("2006"))
	str.WriteString(" ")
	str.WriteString("года. ")
	str.WriteString(res1.Format("15:04"))
	str.WriteString(" ====")

	return str.String()

}

func readFile(t *[]Task) error {

	_, err := os.Stat("tasks.json")
	if err != nil {
		return err
	}

	file, err := os.OpenFile("tasks.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}
	err = json.Unmarshal(data, t)
	if err != nil {
		return err
	}
	return nil
}

func menu() {
	fmt.Println("Что бы добавить, введите: доб ДД.ММ.ГГГГ ЧЧ:ММ \"задача\"")
	fmt.Println("Что бы удалить, введите: уд 3")
	fmt.Println("Что бы изменить статус, введите: ст 3 выполнено")
	fmt.Println("Что бы просмотреть список, введите: список")
	fmt.Println("Что бы сохранить и выйти: Ctrl+d")
	fmt.Print(":: ")
}

func userInput() (string, error) {

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	str := scanner.Text()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return str, nil
}

func parsingUsIn(s string, t *[]Task) error {
	slice := strings.Split(strings.ToLower(s), " ")

	switch slice[0] {
	case "доб":
		err := addTask(slice, t)
		if err != nil {
			return err
		}

	case "уд":
		fmt.Println("Удаление")
	case "ст":
		fmt.Println("Замена статуса")
	case "список":
		if len(*t) == 0 {
			printError()
			return nil
		}
		fmt.Println("Печать списка")
	default:
		fmt.Print("\033[H\033[2J\n")
		fmt.Printf("\t\tКоманда не распознана.\n\n")
	}
	return nil
}

func addTask(s []string, t *[]Task) error {

	strDate := strings.Join(s[1:3], " ")
	task := strings.Join(s[3:], " ")
	date, err := time.Parse("02.01.2006 15:04", strDate)
	if err != nil {
		return err
	}
	id := len(*t) + 1
	newTask := Task{ID: id, Date: date, Task: task, Status: false}
	*t = append(*t, newTask)
	return nil

}

func printError() {
	fmt.Print("\033[H\033[2J\n")
	fmt.Println("\t\tСписок задач пуст.\n")
}
