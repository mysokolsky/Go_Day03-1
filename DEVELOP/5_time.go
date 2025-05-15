// Задача 5: Форматирование времени
// Используй функции для форматирования времени в шаблоне.

// Передай текущее время в шаблон.

// Добавь функцию для форматирования даты через Funcs().


package main

import (
	"os"
	"html/template"
	"time"
"fmt"
)


// Отдельно объявляем функцию форматирования даты
func formatDate(t time.Time) string {
    return t.Format("02-01-2006 15:04")
}

// Функция для создания шаблона с зарегистрированной функцией date
func createTemplate() *template.Template {
    return template.New("это имя может быть пустым").Funcs(template.FuncMap{
        "date": formatDate,
    })
}

func main() {
    // Создаём новый шаблон с функцией
    tmpl := createTemplate()

    // Отдельно парсим файл шаблона (возвращает новый шаблон)
    var err error
    tmpl, err = tmpl.ParseFiles("5_time_template.html")
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