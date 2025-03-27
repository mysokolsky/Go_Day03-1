package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("http://127.0.0.1:9200")
	if err != nil {
		log.Fatalf("Ошибка при подключении к Elasticsearch: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("Ответ от сервера:", resp.Status)
}
