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
}
