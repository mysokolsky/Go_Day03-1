// Создай слайс строк с задачами.

// Передай его в шаблон.

package main

// import "fmt"

import (
	"html/template"
	"os"
)

func main() {

	str := []string{"Первая", "Вторая", "Третья"}

	// fmt.Println(str)

	tmpl := template.Must(template.ParseFiles("2_task_list_template.html"))
	tmpl.Execute(os.Stdout, str)

}
