// Используя шаблон, сгенерируй таблицу умножения для чисел от 1 до 10:

// Создай вложенные слайсы для таблицы.
// Передай их в шаблон.

package main

import (
	"html/template"
	"net/http"
)

// обработчик http-запроса
func handler(w http.ResponseWriter, r *http.Request) {
	var arr = [10][10]int{}
	for i := 1; i <= 10; i++ {
		for j := 1; j <= 10; j++ {
			arr[i-1][j-1] = i * j
		}
	}

	tmpl := template.Must(template.ParseFiles("6_multiply_template.html"))
	tmpl.Execute(w, arr)
}

func main() {

	http.HandleFunc("/", handler)     // вызов функции-обработчика при обращении к корневой директории сервера localhost "/"
	http.ListenAndServe(":8080", nil) // запускаем сервер для приёма запросов на порту 8080.
	// При этом программа перейдёт в режим ожидания, а нам надо будет любым способом, хоть с браузера, обратиться к порту 8080
	// После запуска программы открываем браузер и вводим адрес http://localhost:8080/

}
