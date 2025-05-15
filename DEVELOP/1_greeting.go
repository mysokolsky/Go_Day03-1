// Задача 1: Приветствие

// Загрузить шаблон.

// Передать имя в шаблон через Execute().

package main

import (
    "html/template"
    "os"
)

func main() {
    tmpl := template.Must(template.ParseFiles("1_greeting_template.html"))
    tmpl.Execute(os.Stdout, "Парняга!")
}