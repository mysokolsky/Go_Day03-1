package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	// сначала нужно перейти в корневую папку Go_Day03-1
	// потом поскольку там нет файла go.mod + go.sum, создать их командой go mod init $(basename $(PWD)) или запустить цель make gomodinit
	// потом установить библиотеку go get github.com/elastic/go-elasticsearch/v8 или запустить цель make gogetelastic
	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	// Подключение к Elasticsearch
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Ошибка создания клиента: %s", err)
	}

	// Определяем схему индекса
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"name": {"type": "text"},
				"location": {"type": "geo_point"},
				"rating": {"type": "float"},
				"cuisine": {"type": "keyword"}
			}
		}
	}`

	// Создание индекса "places"
	res, err := es.Indices.Create("places", es.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))))
	if err != nil {
		log.Fatalf("Ошибка создания индекса: %s", err)
	}
	defer res.Body.Close()

	fmt.Println("Ответ Elasticsearch:", res)
}
