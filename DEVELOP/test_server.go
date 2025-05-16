package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Обработка корневого пути
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Привет, мир!")
	})

	// Обработка пути /hello
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "незнакомец"
		}
		fmt.Fprintf(w, "Привет, %s!", name)
	})

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
