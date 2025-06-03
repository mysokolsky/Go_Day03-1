// Общая схема:
// Подключиться к Elasticsearch из Go.

// Сделать запрос к индексу (например, получить все документы).

// Собрать HTML-страницу прямо в коде Go.

// Вернуть её в ответ на HTTP-запрос.

package db

import "Go_Day03-1/src/types"

type Store interface {
	// возвращает список записей, общее количество найденных записей и (или) ошибку
	GetPlaces(limit int, offset int) ([]types.Place, int, error)

	// принимает:
	// limit - какое количество записей выводить на странице
	// offset - после какого id выводить записи из базы в количестве limit

	// возвращает:
	// []types.Place - слайс записей в количестве limit
	// int - общее количество записей с маркером "places"
	// или ошибку
}

func GetPlaces(limit int, offset int) ([]types.Place, int, error) {

}