// Задача 4: Информация о пользователе
// Создай шаблон для вывода данных о пользователе: имя, возраст и список интересов.

// Создай структуру User с полями Name, Age и Interests.

// Заполни данные и передай их в шаблон.

package main

import (
	"html/template"
	"os"
)

type User struct {
	Name      string
	Age       uint
	Interests []string
}

func main() {

	user := User{
		Name:      "Arnold",
		Age:       12,
		Interests: []string{"футбол", "шахматы", "помощь старикам"},
	}

	tmpl := template.Must(template.ParseFiles("4_user_info_template.html"))
	tmpl.Execute(os.Stdout, user)

}
