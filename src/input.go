// Последовательность чтения CSV и загрузки данных в Elasticsearch:
// 1️⃣ Запускаешь сервер Elasticsearch
// Проверяешь, что он работает

// 2️⃣ Создаёшь индекс places с настройкой mappings для приёма корректных json объектов
// Определяешь схему данных (типы полей) перед загрузкой.
// Создаёшь индекс (это маркер, который служит для пометки однотипных данных)

// 3️⃣ Загружаешь CSV в объекты, преобразуешь в json и отправляешь в Elasticsearch
// Читаешь данные из CSV и конвертируешь в структурные объекты
// Преобразуешь структурные объекты в JSON.
// Используешь Bulk API для массовой загрузки объектов на сервер Elastic.

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
