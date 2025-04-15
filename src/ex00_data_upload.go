package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	// "github.com/gocarina/gocsv"
)

// 5.57ms -> 600 records (read)
func main() {

	es := initElasticsearch()

	CreateElasticIndex(es)

	scrollAllDocuments(es, "places")

	bi := initBulkIndexer(es)

	now := time.Now()

	readChannel := make(chan RestaurantsCSV, 25) // создали канал ёмкостью 25 объектов типа RestaurantsCSV

	CSVFilePath := "../materials/data.csv"

	// Open the CSV readFile
	CSVFile, err := os.OpenFile(CSVFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer CSVFile.Close()

	var count int64 = 0 // количество строк
	readFromCSV(CSVFile, readChannel)

	// readChannel := CSVToChannel(CSVFile)

	// Воркеры читают из канала
	var wg sync.WaitGroup
	workerCount := 5 // количество воркеров
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range readChannel {
				fmt.Println(r)
				val, _ := json.Marshal(r.ToRestaurants())
				// objRestaurant, _ := parseLineToRestaurants(r)
				// val, _ := json.Marshal(objRestaurant)

				// fmt.Println(string(val))

				bi.Add(
					context.Background(),
					esutil.BulkIndexerItem{
						Action:     "index",
						DocumentID: strconv.FormatInt(atomic.LoadInt64(&count), 10), //strconv.FormatInt(count, 10), // _id объекта (бывшей строчки из csv)
						Body:       bytes.NewReader(val),
						OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
							if err != nil {
								log.Printf("Ошибка записи в Elasticsearch: %v", err)
							}
							if resp.Error.Type != "" {
								log.Printf("Ошибка в ответе Elasticsearch: %v", resp.Error)
							}
						},
					},
				)

				atomic.AddInt64(&count, 1) // безопасная инкрементация _id
			}
		}()
	}

	wg.Wait()

	// ⚠️ Закрываем BulkIndexer только после всех воркеров!
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Ошибка при закрытии BulkIndexer: %s", err)
	}
	stats := bi.Stats()
	fmt.Printf("✅ Успешно: %d, 🛑 Ошибки: %d\n",
		stats.NumFlushed, stats.NumFailed)

	// Печатаем результат
	fmt.Println("⌛ Время:", time.Since(now), "Кол-во загруженных:", atomic.LoadInt64(&count))

	res, err := es.Indices.Refresh(es.Indices.Refresh.WithIndex("places"))
	if err != nil {
		log.Fatalf("Ошибка при обновлении индекса: %s", err)
	}
	defer res.Body.Close()

	// Считаем и читаем данные из эластика
	countDocuments(es, "places")
	readFromElastic(es)

	// scrollAllDocuments(es, "places")

}

func scrollAllDocuments(es *elasticsearch.Client, index string) {
	// Начинаем запрос с указанием scroll параметра
	pageSize := 100 // сколько документов на странице (можно настроить)
	query := `{"query": {"match_all": {}}}`

	// Выполняем начальный запрос
	res, err := es.Search(
		es.Search.WithIndex(index),
		es.Search.WithScroll(time.Minute),
		es.Search.WithSize(pageSize),
		es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		log.Fatalf("Ошибка при начальном поиске: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Ошибка ответа при начальном поиске: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Ошибка декодирования ответа: %s", err)
	}

	scrollID, ok := result["_scroll_id"].(string)
	if !ok {
		log.Fatalf("Не удалось получить _scroll_id")
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	totalDocuments := int(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	fmt.Printf("Всего документов: %d\n", totalDocuments)

	// Выводим первые полученные документы
	for _, hit := range hits {
		doc := hit.(map[string]interface{})["_source"]
		docJSON, _ := json.MarshalIndent(doc, "", "  ")
		fmt.Println(string(docJSON))
	}

	// Считаем уже выведенные документы
	count := len(hits)

	// Продолжаем извлекать следующие страницы, пока не будет получено 0 документов
	for {
		// Формируем запрос scroll
		res, err := es.Scroll(
			es.Scroll.WithScrollID(scrollID),
			es.Scroll.WithScroll(time.Minute),
		)
		if err != nil {
			log.Fatalf("Ошибка при scroll запросе: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Fatalf("Ошибка ответа scroll запроса: %s", res.String())
		}

		var scrollResult map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&scrollResult); err != nil {
			log.Fatalf("Ошибка декодирования scroll ответа: %s", err)
		}

		// Обновляем scrollID
		scrollID, ok = scrollResult["_scroll_id"].(string)
		if !ok {
			log.Fatalf("Не удалось обновить _scroll_id")
		}

		// Получаем hits
		hitsPage := scrollResult["hits"].(map[string]interface{})["hits"].([]interface{})
		if len(hitsPage) == 0 {
			break // если страниц больше нет, выходим
		}

		// Выводим каждый документ
		for _, hit := range hitsPage {
			doc := hit.(map[string]interface{})["_source"]
			docJSON, _ := json.MarshalIndent(doc, "", "  ")
			fmt.Println(string(docJSON))
		}

		count += len(hitsPage)
	}

	// Очищаем scroll-контекст
	_, err = es.ClearScroll(
		es.ClearScroll.WithScrollID(scrollID),
	)
	if err != nil {
		log.Printf("Ошибка при очистке scroll: %s", err)
	}

	fmt.Printf("Выведено документов: %d\n", count)
}

func CreateElasticIndex(es *elasticsearch.Client) {

	// Читаем только содержимое properties
	schemaFile, err := os.ReadFile("../schema.json")
	if err != nil {
		log.Fatalf("Ошибка при чтении schema.json: %s", err)
	}

	// Оборачиваем содержимое schema.json в структуру с "mappings"
	mapping := fmt.Sprintf(`{"mappings": %s}`, string(schemaFile))

	res, err := es.Indices.Exists([]string{"places"})
	if err != nil {
		log.Fatalf("Ошибка при проверке индекса: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		res, err = es.Indices.Delete([]string{"places"})
		if err != nil {
			log.Fatalf("Ошибка при удалении индекса: %s", err)
		}
	}
	// Создание индекса "places"
	res, err = es.Indices.Create("places", es.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))))
	if err != nil {
		log.Fatalf("Ошибка создания индекса: %s", err)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(res.Body)
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

func getElasticPassword() string {
	filename := "../elastic_pass_MAC.txt"
	if runtime.GOOS == "linux" {
		filename = "../elastic_pass_WSL.txt"
	}
	data, _ := os.ReadFile(filename)
	return string(data)
}

func initElasticsearch() *elasticsearch.Client {
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

	return es
}

func initBulkIndexer(es *elasticsearch.Client) esutil.BulkIndexer {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        es,
		Index:         "places", // 🔁 Имя индекса
		NumWorkers:    5,        // Кол-во воркеров
		FlushBytes:    5e+6,     // Примерно 5MB
		FlushInterval: 5 * time.Second,
		OnFlushStart: func(ctx context.Context) context.Context {
			log.Println("▶ Начало флаша данных...")
			return ctx
		},
	})
	if err != nil {
		log.Fatalf("Ошибка при создании BulkIndexer: %s", err)
	}

	return bi
}

// func readFromCSV(file *os.File, c chan RestaurantsCSV) {

// 	// Настройка CSV-ридера
// 	gocsv.SetCSVReader(func(r io.Reader) gocsv.CSVReader {
// 		reader := csv.NewReader(r)
// 		reader.Comma = '\t'
// 		reader.LazyQuotes = true
// 		reader.FieldsPerRecord = -1
// 		return reader
// 	})

// 	// Запускаем асинхронное чтение и отправку в канал
// 	go func() {

// 		// Пытаемся распарсить весь файл в канал
// 		err := gocsv.UnmarshalToChan(file, c)
// 		if err != nil {

// 			fmt.Fprintf(os.Stderr, "Ошибка при чтении CSV:", err)
// 			close(c) // закрываем канал
// 		}
// 	}()
// }

// func readFromCSV(file *os.File, c chan RestaurantsCSV) {

// 	// Оборачиваем файл в буферизованный ридер
// 	bufReader := bufio.NewReader(file)
// 	// Устанавливаем кастомный CSVReader, который будет читать после заголовка
// 	gocsv.SetCSVReader(func(_ io.Reader) gocsv.CSVReader {

// 		// Создаём обычный CSV-ридер
// 		reader := csv.NewReader(bufReader)
// 		reader.Comma = '\t'
// 		reader.LazyQuotes = true
// 		reader.FieldsPerRecord = -1

// 		// Пропускаем заголовок вручную
// 		if _, err := reader.Read(); err != nil {
// 			log.Fatalf("Ошибка при чтении заголовка CSV: %v", err)
// 		}
// 		// Возвращаем уже подготовленный reader с пропущенным заголовком

// 		return reader
// 	})

// 	// Стартуем горутину с потоком данных в канал
// 	go func() {
// 		// defer close(c)
// 		if err := gocsv.UnmarshalToChan(bufReader, c); err != nil {
// 			log.Fatalf("Ошибка при анмаршалинге CSV: %v", err)
// 		}
// 	}()
// }

// func parseLineToRestaurants(record []string) (Restaurants, error) {

// 	// Парсим ID
// 	id, err := strconv.ParseUint(record[0], 10, 64)
// 	if err != nil {
// 		log.Printf("Ошибка парсинга ID: %v", err)
// 	}

// 	// Парсим координаты
// 	lat, err := strconv.ParseFloat(record[5], 64)
// 	if err != nil {
// 		log.Printf("Ошибка парсинга Latitude: %v", err)
// 	}

// 	lon, err := strconv.ParseFloat(record[4], 64)
// 	if err != nil {
// 		log.Printf("Ошибка парсинга Longitude: %v", err)
// 	}

// 	if err != nil {
// 		return Restaurants{}, err
// 	}

// 	return Restaurants{
// 		ID:      id,
// 		Name:    record[1],
// 		Address: record[2],
// 		Phone:   record[3],
// 		Location: Location{
// 			Latitude:  lat,
// 			Longitude: lon,
// 		},
// 	}, nil
// }

// func initCSVReader(file *os.File) *csv.Reader {

// 	bufReader := bufio.NewReader(file)

// 	reader := csv.NewReader(bufReader)
// 	reader.Comma = '\t'
// 	reader.LazyQuotes = true
// 	reader.FieldsPerRecord = -1

// 	return reader
// }

// func CSVToChannel(file *os.File) chan []string {

// 	reader := initCSVReader(file)

// 	// Пропускаем заголовок
// 	if _, err := reader.Read(); err != nil {
// 		log.Fatalf("Ошибка чтения заголовка: %v", err)
// 	}
// 	ch := make(chan []string, 100) // создали канал ёмкостью 100 объектов типа []string
// 	go func() {
// 		defer close(ch)
// 		for {
// 			record, err := reader.Read()
// 			if err == io.EOF {
// 				break
// 			}
// 			if err != nil {
// 				log.Printf("Ошибка чтения строки: %v", err)
// 				continue
// 			}

// 			if len(record) < 6 {
// 				log.Printf("Пропуск строки с недостаточным количеством полей: %+v", record)
// 				continue
// 			}
// 			ch <- record
// 		}
// 	}()
// 	return ch
// }

func readFromCSV(file *os.File, c chan RestaurantsCSV) {
	bufReader := bufio.NewReader(file)

	reader := csv.NewReader(bufReader)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	// Пропускаем заголовок
	if _, err := reader.Read(); err != nil {
		log.Fatalf("Ошибка чтения заголовка: %v", err)
	}

	go func() {
		defer close(c)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Ошибка чтения строки: %v", err)
				continue
			}

			if len(record) < 6 {
				log.Printf("Пропуск строки с недостаточным количеством полей: %+v", record)
				continue
			}

			// Парсим ID
			id, err := strconv.ParseUint(record[0], 10, 64)
			if err != nil {
				log.Printf("Ошибка парсинга ID: %v", err)
				continue
			}

			// Парсим координаты
			lat, err := strconv.ParseFloat(record[5], 64)
			if err != nil {
				log.Printf("Ошибка парсинга Latitude: %v", err)
				continue
			}

			lon, err := strconv.ParseFloat(record[4], 64)
			if err != nil {
				log.Printf("Ошибка парсинга Longitude: %v", err)
				continue
			}

			//Собираем структуру
			c <- RestaurantsCSV{
				ID:        id,
				Name:      record[1],
				Address:   record[2],
				Phone:     record[3],
				Latitude:  lat,
				Longitude: lon,
			}

		}
	}()

}

func countDocuments(es *elasticsearch.Client, index string) {
	res, err := es.Count(
		es.Count.WithIndex(index),
	)
	if err != nil {
		log.Fatalf("Ошибка при подсчёте документов: %s", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println("📊 Количество документов в индексе:", string(body))
}

func readFromElastic(es *elasticsearch.Client) {
	query := `{
		"query": {
			"match_all": {}
		}
	}`

	res, err := es.Search(
		es.Search.WithIndex("places"),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Ошибка поиска: %s", err)
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source Restaurants `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Ошибка при декодировании: %s", err)
	}

	fmt.Println("Печатаем результаты по индексу places из базы Elastic:")

	// Выводим найденные рестораны
	for _, hit := range result.Hits.Hits {
		fmt.Printf("🍽️  %s, %s [%s], { %v, %v }\n", hit.Source.Name, hit.Source.Address, hit.Source.Phone, hit.Source.Location.Latitude, hit.Source.Location.Longitude)
	}
}
