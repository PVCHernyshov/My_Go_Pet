package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID     int       `json:"id"`
	Date   time.Time `json:"date"`
	Task   string    `json:"task"`
	Status bool      `json:"status"`
}

var exit error = errors.New("exit")
var errorPrint bool

func main() {

	var tasks []Task
	var errRead error
	errRead = readFile(&tasks)

	if errRead != nil {
		if os.IsNotExist(errRead) {
			printError()
		} else if errors.Is(errRead, io.EOF) {
			printError()
		} else {
			fmt.Println(errRead)
			return
		}
	} else {
		parsingUsIn("список", &tasks)
	}

	for {
		menu()

		resUsIn, errUsIn := userInput()

		if errUsIn != nil {
			fmt.Println("Не удалось прочитать ввод. Ошибка:", errUsIn, "\n", "завершение программы.")
			return
		}

		err := parsingUsIn(resUsIn, &tasks)
		if errors.Is(err, exit) {
			fmt.Println("До новых встреч.\n")
			return
		}
		if err != nil {
			continue
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

	file, err := os.Open("tasks.json")
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(t)
	if err != nil {
		return err
	}

	return nil
}

func menu() {
	fmt.Println("Что бы добавить, введите: добавить ДД.ММ.ГГГГ ЧЧ-ММ \"задача\"")
	fmt.Println("Что бы удалить, введите: удалить №")
	fmt.Println("Что бы изменить статус, введите: статус № \"новый статус\"")
	fmt.Println("Что бы просмотреть список, введите: список")
	fmt.Println("Что бы сохранить и выйти введите: выход")
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
	slice := strings.Fields(s)
	if len(slice) == 0 {
		newMessage()
		fmt.Println("\t\tПустой ввод.\n")
		return nil
	}

	switch strings.ToLower(slice[0]) {

	case "добавить":
		if len(slice) < 3 {
			invalidInput()
			return nil
		}
		_, err := time.Parse("02.01.2006", slice[1])
		if err != nil {
			newMessage()
			fmt.Println("\t\tВведён не верный формат даты.\n")
			return nil
		}
		_, err = time.Parse("15-04", slice[2])
		if err != nil {
			newMessage()
			fmt.Println("\t\tВведён не верный формат времени.\n")
			return nil
		}

		err = addTask(slice, t)
		if err != nil {
			return err
		}

	case "удалить":
		if len(slice) != 2 {
			invalidInput()
			return nil
		}

		num, err := validationAndcheckID(slice, *t)
		if err != nil {
			newMessage()
			fmt.Println(err)
			return nil
		}

		err = deleteTask(num, t)
		if err != nil {
			newMessage()
			fmt.Println(err)
			return nil
		}

	case "статус":
		if len(slice) < 3 {
			invalidInput()
			return nil
		}

		num, err := validationAndcheckID(slice, *t)
		if err != nil {
			newMessage()
			fmt.Println(err)
			return nil
		}

		status := strings.Join(slice[2:], " ")

		if status != "выполнено" && status != "не выполнено" {
			newMessage()
			fmt.Println("\tНе верно указан статус для задачи.\n")
			return nil
		}
		err = changeOfStatus(num, status, t)
		if err != nil {
			newMessage()
			fmt.Println(err)
			return nil
		}

	case "список":
		if len(*t) == 0 {
			printError()
			return nil
		}
		newMessage()
		fmt.Println("\t\tСписок задач.\n")
		printTasks(*t)
		fmt.Println()

	case "выход":
		return exit
	default:
		newMessage()
		fmt.Println("\t\tКоманда не распознана.\n")
	}
	return nil
}

func addTask(s []string, t *[]Task) error {

	task, err := createTask(s, *t)
	if err != nil {
		return err
	}

	*t = append(*t, task)
	newMessage()
	fmt.Printf("\t\tДобавлена задача.\n\n\"%s\" на %s \n\n", task.Task, task.Date.Format("02.01.2006 15-04"))

	sortTasks(*t)

	return writeFile(t)
}

func parsingDataAddTask(s []string) (string, time.Time, error) {
	strDate := strings.Join(s[1:3], " ")
	task := strings.Join(s[3:], " ")
	date, err := time.Parse("02.01.2006 15-04", strDate)
	if err != nil {
		return "", time.Time{}, err
	}

	return task, date, nil
}

func printError() {
	if errorPrint {
		fmt.Print("\033[H\033[2J\n")
	}
	errorPrint = true
	fmt.Println(heading(time.Now()))
	fmt.Println("\t\tСписок задач пуст.\n")
}

func newMessage() {
	fmt.Print("\033[H\033[2J\n")
	fmt.Println(heading(time.Now()))
}

func sortTasks(t []Task) {
	slices.SortFunc(t, func(a, b Task) int {
		return a.Date.Compare(b.Date)
	})
}

func printTasks(t []Task) {
	var stat string
	for _, task := range t {
		if !task.Status {
			stat = "Невыполнено"
		} else {
			stat = "Выполнено"
		}

		if (task.Date.Before(time.Now())) && !task.Status {
			stat = "Просрочено."
		}
		fmt.Printf("№ %d: %s %s - %s\n", task.ID, task.Date.Format("02.01.2006 15-04"), task.Task, stat)
	}
}

func findfreeID(t []Task) int {
	for id := 1; ; id++ {
		flag := true
		for _, task := range t {
			if id == task.ID {
				flag = false
				break
			}
		}
		if flag {
			return id
		}
	}
}

func changeOfStatus(n int, s string, t *[]Task) error {
	for e := range *t {
		if n == (*t)[e].ID {
			switch s {
			case "выполнено":
				(*t)[e].Status = true
				s = "Выполнено"
			case "не выполнено":
				(*t)[e].Status = false
				s = "Не выполнено"
			}
			newMessage()
			fmt.Printf("\tСтатус задачи %s изменён на - %s.\n\n", (*t)[e].Task, s)
			return writeFile(t)
		}
	}

	return fmt.Errorf("Ошибка изменения статуса. Статус не изменён.")
}

func checkID(id int, t []Task) error {
	for _, task := range t {
		if id == task.ID {
			return nil
		}
	}
	return fmt.Errorf("\t\tID %d не найден.\n\n", id)
}

func invalidInput() {
	newMessage()
	fmt.Println("\tВведённые данные не соответствуют формату.\n")
}

func validationAndcheckID(slice []string, t []Task) (int, error) {
	num, err := strconv.Atoi(slice[1])
	if err != nil {
		return 0, fmt.Errorf("\t\tОшибочно введён ID.\n")
	}

	err = checkID(num, t)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func deleteTask(n int, t *[]Task) error {
	*t = slices.DeleteFunc(*t, func(id Task) bool {
		return id.ID == n
	})
	newMessage()
	fmt.Printf("\t\tЗадача ID: %d удалена.\n\n", n)

	err := writeFile(t)
	if err != nil {
		return err
	}
	return nil
}

func writeFile(t *[]Task) error {

	file, err := os.Create("tasks.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(t)
}

func createTask(s []string, t []Task) (Task, error) {
	task, date, err := parsingDataAddTask(s)
	if err != nil {
		return Task{}, err
	}

	id := findfreeID(t)

	newTask := Task{ID: id, Date: date, Task: task, Status: false}
	return newTask, nil

}
