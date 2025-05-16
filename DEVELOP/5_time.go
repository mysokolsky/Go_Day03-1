// Задача 5: Форматирование времени
// Используй функции для форматирования времени в шаблоне.

// Передай текущее время в шаблон.

// Добавь функцию для форматирования даты через Funcs().

package main

import (
	"fmt"
	"html/template"
	"os"
	"time"
)

// Отдельно объявляем функцию форматирования даты
func formatDate(t time.Time) string {
	return t.Format("02-01-2006 15:04")
}


func main() {
	// Создаём новый шаблон
	tmpl := template.New("новый шаблон")

	funcMap := template.FuncMap{} // инициализируем пустую мапу
	funcMap["date"] = formatDate

	tmpl = tmpl.Funcs(funcMap)

	// Отдельно парсим файл шаблона (возвращает новый шаблон)
	var err error
	tmpl, err = tmpl.ParseFiles("5_time_template.html") // здесь создаётся новый подшаблон "5_time_template.html", который является ребёнком для первоначального "новый шаблон"
	if err != nil {
		fmt.Println("Ошибка при парсинге файла:", err)
		return
	}

	// Выполняем шаблон с текущим временем, указывая имя шаблона из файла
	err = tmpl.ExecuteTemplate(os.Stdout, "5_time_template.html", time.Now())
	if err != nil {
		fmt.Println("Ошибка при выполнении шаблона:", err)
	}

	// распечатка имеющихся шаблонов
	println("\n\n >>>> Шаблоны: <<<<\n")
	for _, t := range tmpl.Templates() {
		fmt.Println(t.Name())
	}
}
