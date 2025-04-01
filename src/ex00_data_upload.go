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
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	// сначала нужно перейти в корневую папку Go_Day03-1
	// потом поскольку там нет файла go.mod + go.sum, создать их командой go mod init $(basename $(PWD)) или запустить цель make gomodinit
	// потом установить библиотеку go get github.com/elastic/go-elasticsearch/v8 или запустить цель make gogetelastic
	"github.com/elastic/go-elasticsearch/v8"
)

// ID
// Name
// Address
// Phone
// Longitude
// Latitude

type Location struct {
	Longitude float64 `csv:"longitude" json:"lon"`
	Latitude  float64 `csv:"latitude" json:"lat"`
}

type Restaurants struct {
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

func getElasticPassword() string {
	filename := "../elastic_pass_MAC.txt"
	if runtime.GOOS == "linux" {
		filename = "../elastic_pass_WSL.txt"
	}
	data, _ := os.ReadFile(filename)
	return string(data)
}

func main() {
	// Создаем клиент Elasticsearch с использованием HTTPS и аутентификации
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200", // Ваш сервер Elasticsearch
		},
		Username: "elastic",            // Имя пользователя
		Password: getElasticPassword(), // Ваш пароль

		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Отключаем проверку сертификата (для разработки)
			},
		},
	})
	if err != nil {
		log.Fatalf("Ошибка при создании клиента: %v", err)
	}

	// Проверяем соединение с сервером
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Ошибка при подключении к Elasticsearch: %v", err)
	}
	defer res.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Ошибка чтения ответа: %v", err)
	}

	// Выводим ответ
	fmt.Println("Ответ от сервера:", string(body))

	// Определяем схему индекса
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 1
		},
		"mappings": {
			"properties": {
				"name": {"type": "text"},
				"address": {"type": "text"},
				"phone": {"type": "text"},
				"location": {"type": "geo_point"}
			}
		}
	}`

	// Создание индекса "places"
	res, err = es.Indices.Create("places", es.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))))
	if err != nil {
		log.Fatalf("Ошибка создания индекса: %s", err)
	}
	defer res.Body.Close()

	// Читаем тело ответа
	body, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Ошибка чтения ответа: %s", err)
	}

	// Проверяем код ответа
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		log.Fatalf("Ошибка при создании индекса: %s\nОтвет: %s", res.Status(), string(body))
	}

	fmt.Println("✅ Индекс 'places' успешно создан!")
	fmt.Println("Ответ от Elasticsearch:", string(body))
}
