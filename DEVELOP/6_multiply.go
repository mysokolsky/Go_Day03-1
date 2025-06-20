// Используя шаблон, сгенерируй таблицу умножения для чисел от 1 до 10:

// Создай вложенные слайсы для таблицы.
// Передай их в шаблон.

package main

import (
	"html/template"
	"net/http"
	"log"
	"os/exec"
	"runtime"
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


func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default: // Linux и WSL
		cmd = exec.Command("xdg-open", url)
	}

	err := cmd.Start()
	if err != nil {
		log.Printf("Не удалось открыть браузер: %v", err)
	}
}





func main() {

	http.HandleFunc("/", handler)     // вызов функции-обработчика при обращении к корневой директории сервера localhost "/"


	go openBrowser("http://localhost:8080") // Открываем браузер в отдельной горутине


	log.Printf("Запускаем сервер http://localhost:8080")
	
	http.ListenAndServe(":8080", nil) // запускаем сервер для приёма запросов на порту 8080.
	// При этом программа перейдёт в режим ожидания, а нам надо будет любым способом, хоть с браузера, обратиться к порту 8080
	// После запуска программы открываем браузер и вводим адрес http://localhost:8080/




}
