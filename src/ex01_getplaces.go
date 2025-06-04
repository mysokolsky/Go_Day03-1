package main

import (
	"Go_Day03-1/src/db"
	"Go_Day03-1/src/types"
	"log"
    "net/http"
)



// 1. Получать ?page=N через r.URL.Query().Get("page").

// 2. Конвертировать его в int.

// 3. Вычислять, с какой страницы начинать, и передавать в search_after.

// 4. Полученные записи из elastic передавать в Page и рендерить шаблон.


func handlePlaces(w http.ResponseWriter, r *http.Request) {
    // Парсим параметры из URL
    after := r.URL.Query().Get("after")
    limit := 10

    // Вызываем метод получения данных из интерфейса
    places, nextAfterID, err := db.GetPlaces(limit, after)
    if err != nil {
        http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
        return
    }

    // Строим данные для шаблона
    data := types.Page{
        Count:    len(places), // можно заменить на общее кол-во, если считаешь в GetPlaces
        Places:   places,
        NotFirst: after != "", // если after пустой — это первая страница
        NotLast:  nextAfterID != "", // если есть ещё страницы
        Prev:     "", // пока можно не реализовывать
        Next:     nextAfterID,
        Last:     "PLACEHOLDER_LAST", // можно позже реализовать
    }

    tmpl := template.Must(template.ParseFiles("page_template.html"))
    tmpl.Execute(w, data)
}

func main() {



	http.HandleFunc("/", handlePlaces)
	log.Println("Сервер запущен на http://localhost:8888")
	http.ListenAndServe(":8888", nil)


}