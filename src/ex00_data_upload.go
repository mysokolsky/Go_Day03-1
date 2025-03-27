package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"

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

func main() {
	// Создаем клиент Elasticsearch с использованием HTTPS и аутентификации
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200", // Ваш сервер Elasticsearch
		},
		Username: "elastic",              // Имя пользователя
		Password: "aWl-v68rt1TSmzKuqPk-", // Ваш пароль
		// Insecure: true,                   // Включить, если не хотите проверять сертификат SSL (не рекомендуется для продакшн-окружения)
		// Примечание: можете добавить параметры для пропуска проверки сертификатов, если нужно
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
				"ID": {"type": "text"},
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
