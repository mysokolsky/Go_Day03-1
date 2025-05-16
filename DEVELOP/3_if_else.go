// Задача 3: Условный вывод
// Добавь проверку статуса пользователя. Если он вошел, показывай приветствие, если нет — сообщение "Пожалуйста, войдите".

// Передай true или false в шаблон.

package main

import (
	"html/template"
	"os"
)

func main() {

	tmpl := template.Must(template.ParseFiles("3_if_else_template.html"))
	tmpl.Execute(os.Stdout, false)

}
