// Последовательность работы с Elasticsearch:
// 1️⃣ Запускаешь Elasticsearch
// Проверяешь, что сервер работает

// 2️⃣ Создаёшь индекс с mappings и settings
// Определяешь схему данных (типы полей) перед загрузкой.
// Отправляешь запрос на создание индекса
// Теперь индекс places создан и готов к приёму данных.

// 3️⃣ Загружаешь CSV в Go и отправляешь в Elasticsearch
// Читаешь CSV в Go.
// Преобразуешь строки в JSON.
// Используешь Bulk API для массовой загрузки данных.

// 🚀 Почему сначала создаётся индекс?
// Если не создать индекс заранее, Elasticsearch сам определит типы данных, и они могут быть неправильными.
// Например, если он примет число 123 за long, а потом встретит 123.5, то будет ошибка.
// Создав индекс с mappings, ты гарантируешь правильные типы данных.
// 🔥 Вывод:
// Сначала создаём индекс (PUT /places).
// Потом загружаем CSV в этот индекс с помощью Bulk API.

package main

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type RestaurantsBASE struct {
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}
