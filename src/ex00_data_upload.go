package main

import (
	// "bufio"
	"Go_Day03-1/src/types"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
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
)

func main() {

	es := initElasticsearch() // инициализируем клиент для работы с базой Elastic

	CreateElasticIndex(es) // Создаём маркер для объединения объектов json в Elastic-е

	bi := initBulkIndexer(es) // инициализируем Bulk API для заливки json в Elastic

	now := time.Now() // засекаем время для вычисления длительности выполнения программы

	CSVFilePath := "../materials/data.csv"

	// Открываем CSV файл для чтения
	CSVFile, err := os.OpenFile(CSVFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Ошибка при открытии файла CSV: %s", err)
	}
	defer CSVFile.Close()

	var count uint64 = 0 // счётчик записей(строк) в CSV

	readChannel := make(chan types.InputType, 25) // Открываем канал на 25 записей. InputType - один из двух типов данных, прописанный в input_auto.go и input_manual.go,
	// который подставляется при условной компиляции go run -tags=manual . или go run .

	CSVLinesToChannel(CSVFile, readChannel) // вызов одной из функций парсинга, в зависимости от условной компиляции

	// Воркеры читают из канала
	var wg sync.WaitGroup
	workerCount := 5 // количество воркеров, которые будут одновременно параллельно заливать данные в Elastic
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			FromChannelToElastic(readChannel, bi, &count) // запуск функции заливки в Elastic
		}()
	}

	wg.Wait()

	// ⚠️ Закрываем BulkIndexer только после заершения работы всех воркеров!
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Ошибка при закрытии BulkIndexer: %s", err)
	}

	PrintResult(bi, es, &now, &count) // Выводим данные из Elastic в консоль

}

// ОСНОВНЫЕ РАБОЧИЕ ФУНКЦИИ

// Чтение из канала, конвертация в объекты Restaurants, маршалинг в json и заливка в Elastic
func FromChannelToElastic(
	ch chan types.InputType,
	bi esutil.BulkIndexer,
	count *uint64) {

	for r := range ch {
		// fmt.Println(r)

		val := GetValue(r)

		fmt.Println(string(val))

		bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.FormatUint(atomic.LoadUint64(count), 10), //strconv.FormatUint(objRestaurant.ID, 10),  // _id объекта (бывшей строчки из csv)
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

		atomic.AddUint64(count, 1) // безопасная инкрементация _id
	}
}

// СЛУЖЕБНЫЕ ФУНКЦИИ для подключения Elastic и других важных модулей

// Создание index для маркировки заливаемых в Elastic объектов json
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
		log.Fatalf("Ошибка при проверке наличия индекса 'places': %s", err)
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
		log.Fatalf("Ошибка чтения ответа от сервера после создания индекса: %s", err)
	}

	// Проверяем код ответа
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		log.Fatalf("Ошибка при создании индекса: %s\nОтвет: %s", res.Status(), string(body))
	}

	fmt.Println("✅ Индекс 'places' успешно создан!")
	fmt.Println("Ответ от сервера Elastic:", string(body))

}

// Считывание пароля для подключения Elastic
func getElasticPassword() string {
	filename := "../elastic_pass_MAC.txt"
	if runtime.GOOS == "linux" {
		filename = "../elastic_pass_WSL.txt"
	}
	data, _ := os.ReadFile(filename)
	return string(data)
}

// Инициализация клиента Elastic
func initElasticsearch() *elasticsearch.Client {

	// Создаем клиент для общения с сервером Elastic. С использованием HTTPS и аутентификации
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200", // Адрес сервера и порт для общения с Elastic
		},
		Username: "elastic",            // Имя пользователя
		Password: getElasticPassword(), // Пароль

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
		log.Fatalf("🛑 Ошибка при подключении к Elasticsearch:  \n%+v \n⚠️  Сначала запустите сервер Elastic при помощи команды 'make' в корневой папке проекта", err)
	}
	defer res.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Ошибка чтения ответа: %v", err)
	}

	// Выводим ответ
	log.Println("Ответ от сервера:", string(body))

	return es
}

// Инициализация и настройка Bulk API индексатора
func initBulkIndexer(es *elasticsearch.Client) esutil.BulkIndexer {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        es,
		Index:         "places",        // 🔁 Имя маркера для пометки заливаемых данных
		NumWorkers:    5,               // Кол-во воркеров которые одновременно будут заливать в Elastic
		FlushBytes:    5e+6,            // ограничиваем буфер данных для заливки в Elastic = 5 мегабайтам
		FlushInterval: 1 * time.Second, // или по времени не более 5 секунд
		OnFlushStart: func(ctx context.Context) context.Context {
			log.Println("🔁 Начало флаша данных...")
			return ctx
		},
		OnFlushEnd: func(ctx context.Context) {
			log.Println("✅ Флаш завершён")
		},
		OnError: func(ctx context.Context, err error) {
			log.Printf("❌ Ошибка при флаше: %s", err)
		},
	})
	if err != nil {
		log.Fatalf("Ошибка при создании BulkIndexer: %s", err)
	}

	return bi
}

// Инициализация и настройка CSV-читателя
func initCSVReader(file *os.File) *csv.Reader {

	// bufReader := bufio.NewReader(file) // создаём буфер для файла

	reader := csv.NewReader(file) // инициализируем ридер для CSV файла
	reader.Comma = '\t'           // разделитель полей в файле CSV - табуляция
	reader.LazyQuotes = true      // разрешаем некорректные или незакрытые кавычки в CSV
	reader.FieldsPerRecord = -1   // определяем, что в одной записи может быть разное количество полей

	return reader
}

// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ для анализа залитого в Elastic

// Основная функция вывода результатов и статистики на печать
func PrintResult(bi esutil.BulkIndexer, es *elasticsearch.Client, now *time.Time, count *uint64) {

	stats := bi.Stats()
	fmt.Printf("✅ Успешно: %d, 🛑 Ошибки: %d\n",
		stats.NumFlushed, stats.NumFailed)

	// Печатаем результат
	fmt.Println("⌛ Время:", time.Since(*now), "Кол-во загруженных:", atomic.LoadUint64(count))

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

// Подсчёт в Elastic количества залитых json-ов
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

// Чтение нескольких json в консоль
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
				Source types.Restaurants `json:"_source"`
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

// // Печать всех залитых json в консоль
// func scrollAllDocuments(es *elasticsearch.Client, index string) {
// 	// Начинаем запрос с указанием scroll параметра
// 	pageSize := 100 // сколько документов на странице (можно настроить)
// 	query := `{"query": {"match_all": {}}}`

// 	// Выполняем начальный запрос
// 	res, err := es.Search(
// 		es.Search.WithIndex(index),
// 		es.Search.WithScroll(time.Minute),
// 		es.Search.WithSize(pageSize),
// 		es.Search.WithBody(strings.NewReader(query)),
// 	)
// 	if err != nil {
// 		log.Fatalf("Ошибка при начальном поиске: %s", err)
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		log.Fatalf("Ошибка ответа при начальном поиске: %s", res.String())
// 	}

// 	var result map[string]interface{}
// 	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
// 		log.Fatalf("Ошибка декодирования ответа: %s", err)
// 	}

// 	scrollID, ok := result["_scroll_id"].(string)
// 	if !ok {
// 		log.Fatalf("Не удалось получить _scroll_id")
// 	}

// 	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
// 	totalDocuments := int(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
// 	fmt.Printf("Всего документов: %d\n", totalDocuments)

// 	// Выводим первые полученные документы
// 	for _, hit := range hits {
// 		doc := hit.(map[string]interface{})["_source"]
// 		docJSON, _ := json.MarshalIndent(doc, "", "  ")
// 		fmt.Println(string(docJSON))
// 	}

// 	// Считаем уже выведенные документы
// 	count := len(hits)

// 	// Продолжаем извлекать следующие страницы, пока не будет получено 0 документов
// 	for {
// 		// Формируем запрос scroll
// 		res, err := es.Scroll(
// 			es.Scroll.WithScrollID(scrollID),
// 			es.Scroll.WithScroll(time.Minute),
// 		)
// 		if err != nil {
// 			log.Fatalf("Ошибка при scroll запросе: %s", err)
// 		}
// 		defer res.Body.Close()

// 		if res.IsError() {
// 			log.Fatalf("Ошибка ответа scroll запроса: %s", res.String())
// 		}

// 		var scrollResult map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&scrollResult); err != nil {
// 			log.Fatalf("Ошибка декодирования scroll ответа: %s", err)
// 		}

// 		// Обновляем scrollID
// 		scrollID, ok = scrollResult["_scroll_id"].(string)
// 		if !ok {
// 			log.Fatalf("Не удалось обновить _scroll_id")
// 		}

// 		// Получаем hits
// 		hitsPage := scrollResult["hits"].(map[string]interface{})["hits"].([]interface{})
// 		if len(hitsPage) == 0 {
// 			break // если страниц больше нет, выходим
// 		}

// 		// Выводим каждый документ
// 		for _, hit := range hitsPage {
// 			doc := hit.(map[string]interface{})["_source"]
// 			docJSON, _ := json.MarshalIndent(doc, "", "  ")
// 			fmt.Println(string(docJSON))
// 		}

// 		count += len(hitsPage)
// 	}

// 	// Очищаем scroll-контекст
// 	_, err = es.ClearScroll(
// 		es.ClearScroll.WithScrollID(scrollID),
// 	)
// 	if err != nil {
// 		log.Printf("Ошибка при очистке scroll: %s", err)
// 	}

// 	fmt.Printf("Выведено документов: %d\n", count)
// }
