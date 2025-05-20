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

	key := "date"
	value := formatDate

	funcMap := template.FuncMap{key: value} // инициализируем объект спецмапу шаблонов для записи в неё функциий и сразу инициализируем в ней запись алиаса date который ссылается на функцию formatDate

	// Создаём новый шаблон
	tmpl := template.New("новый шаблон")

	// подключаем мапу к шаблону
	tmpl = tmpl.Funcs(funcMap)

	// // Можно было сократить все предыдущие строки до одной строки
	// tmpl := template.New("новый шаблон").Funcs(map[string]interface{}{"date": formatDate}) // сразу создаём шаблон и инициализируем в него алиас date, который соотвествует вызову функции формата даты

	// Отдельно парсим файл шаблона (при этом создаётся подшаблон к имеющемуся)
	var err error
	tmpl, err = tmpl.ParseFiles("5_time_template.html") // здесь создаётся новый подшаблон "5_time_template.html", который является наследником для "новый шаблон"
	if err != nil {
		fmt.Println("Ошибка при парсинге файла:", err)
		return
	}

	// Выполняем шаблон "5_time_template.html", передавая в качестве аргумента для алиаса date текущее время
	err = tmpl.ExecuteTemplate(os.Stdout, "5_time_template.html", time.Now())
	if err != nil {
		fmt.Println("Ошибка при выполнении шаблона:", err)
	}

	// распечатка имеющихся шаблонов
	println("\n\n >>>> Шаблоны: <<<<\n")
	for _, t := range tmpl.Templates() {
		fmt.Println(t.Name())
	}

	// // Можно сделать вообще всё одной строчкой вот так:
	// template.Must(template.New("новый шаблон").Funcs(map[string]interface{}{"date": formatDate}).ParseFiles("5_time_template.html")).ExecuteTemplate(os.Stdout, "5_time_template.html", time.Now())

}
