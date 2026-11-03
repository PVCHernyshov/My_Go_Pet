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
	fmt.Print("\033[H\033[2J\n")
	var errRead error
	errRead = readFile(&tasks)

	if errRead != nil {
		if os.IsNotExist(errRead) {
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
			fmt.Println("Сохраняем и выходим.\nДо новых встреч!.\n")
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

	_, err := os.Stat("tasks.json")
	if err != nil {
		return err
	}

	file, err := os.Open("tasks.json")
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
	fmt.Println("Что бы добавить, введите: доб ДД.ММ.ГГГГ ЧЧ-ММ \"задача\"")
	fmt.Println("Что бы удалить, введите: уд 3")
	fmt.Println("Что бы изменить статус, введите: ст 3 выполнено")
	fmt.Println("Что бы просмотреть список, введите: список")
	fmt.Println("Что бы сохранить и выйти введите: вых")
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

	case "доб":
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

	case "уд":
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

	case "ст":
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

	case "вых":
		return exit
	default:
		newMessage()
		fmt.Println("\t\tКоманда не распознана.\n")
	}
	return nil
}

func addTask(s []string, t *[]Task) error {

	strDate := strings.Join(s[1:3], " ")
	task := strings.Join(s[3:], " ")
	date, err := time.Parse("02.01.2006 15-04", strDate)
	if err != nil {
		return err
	}
	id := 0
	findfreeID(&id, *t)
	newTask := Task{ID: id, Date: date, Task: task, Status: false}
	*t = append(*t, newTask)
	newMessage()
	fmt.Printf("\t\tДобавлена задача.\n\n\"%s\" на %s \n\n", task, date.Format("02.01.2006 15-04"))

	sortTasks(*t)

	err = writeFile(t)
	if err != nil {
		return err
	}

	return nil
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
		fmt.Printf("ID %d: %s %s - %s\n", task.ID, task.Date.Format("02.01.2006 15-04"), task.Task, stat)
	}
}

func findfreeID(id *int, t []Task) {
	for i := 1; i <= len(t)+1; i++ {
		flag := true
		for _, task := range t {
			if i == task.ID {
				flag = false
				break
			}
		}
		if flag {
			*id = i
			return
		}
	}
}

func changeOfStatus(n int, s string, t *[]Task) error {
	for j := range *t {
		if n == (*t)[j].ID {
			switch s {
			case "выполнено":
				(*t)[j].Status = true
				newMessage()
				fmt.Printf("\tСтатус задачи %s изменён на - Выполнено.\n\n", (*t)[j].Task)
			case "не выполнено":
				(*t)[j].Status = false
				newMessage()
				fmt.Printf("\tСтатус задачи %s изменён на - Не выполнено.\n\n", (*t)[j].Task)
			}
		}
	}
	err := writeFile(t)
	if err != nil {
		return err
	}
	return nil
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

	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	os.WriteFile("tasks.json", b, 0644)
	return nil
}
