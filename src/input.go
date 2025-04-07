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

import (
// "bytes"
// "crypto/tls"
// "fmt"
// "io"
// "log"
// "net/http"
// "os"
// "runtime"
// сначала нужно перейти в корневую папку Go_Day03-1
// потом поскольку там нет файла go.mod + go.sum, создать их командой go mod init $(basename $(PWD)) или запустить цель make gomodinit
// потом установить библиотеку go get github.com/elastic/go-elasticsearch/v8 или запустить цель make gogetelastic
// "github.com/elastic/go-elasticsearch/v8"
// сначала нужно закачать библиотеку командой go get github.com/gocarina/gocsv
// "encoding/csv"
// "github.com/gocarina/gocsv"
)

// ID
// Name
// Address
// Phone
// Longitude
// Latitude

type Location struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type Restaurants struct {
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

type RestaurantsCSV struct {
	Name      string  `csv:"Name"`
	Address   string  `csv:"Address"`
	Phone     string  `csv:"Phone"`
	Longitude float64 `csv:"Longitude"`
	Latitude  float64 `csv:"Latitude"`
}

func (r RestaurantsCSV) ToRestaurants() Restaurants {
	return Restaurants{
		Name:    r.Name,
		Address: r.Address,
		Phone:   r.Phone,
		Location: Location{
			Longitude: r.Longitude,
			Latitude:  r.Latitude,
		},
	}
}
