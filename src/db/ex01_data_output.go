// Общая схема:
// Подключиться к Elasticsearch из Go.

// Сделать запрос к индексу (например, получить все документы).

// Собрать HTML-страницу прямо в коде Go.

// Вернуть её в ответ на HTTP-запрос.



// В контексте Elasticsearch, hits — это ключ в ответе на поисковый запрос,
// который содержит найденные документы. Он является частью JSON-структуры,
// которую возвращает Elasticsearch в ответ на запрос типа search.

// 🧩 Структура ответа search в Elasticsearch
// Вот как примерно выглядит стандартный ответ от Elasticsearch:

// json

// {
//   "took": 5,
//   "timed_out": false,
//   "_shards": { ... },
//   "hits": {
//     "total": {
//       "value": 13649,
//       "relation": "eq"
//     },
//     "max_score": 1.0,
//     "hits": [
//       {
//         "_index": "places",
//         "_id": "AX1234...",
//         "_score": 1.0,
//         "_source": {
//           "Name": "Café Central",
//           "Address": "Vienna, Austria",
//           "Phone": "+43 1 531..."

//           // другие поля
//         }
//       },
//       {
//         "_index": "places",
//         "_id": "AX5678...",
//         "_score": 1.0,
//         "_source": {
//           "Name": "La Boulangerie",
//           "Address": "Paris, France",
//           "Phone": "+33 1 45..."

//           // другие поля
//         }
//       }
//       // ... и так далее
//     ]
//   }
// }

// 🔍 Разбор hits
// hits.total — общее количество документов, удовлетворяющих запросу.

// hits.hits — массив найденных документов.

// Каждый элемент в hits.hits — это один найденный документ, который содержит:

// _index — индекс, где найден документ.

// _id — внутренний ID документа в Elasticsearch.

// _source — содержимое документа (то, что ты загружал).

// 📘 В Go при работе с hits
// Если ты используешь github.com/elastic/go-elasticsearch/v8,
// ты обычно получаешь res.Body, читаешь JSON,
// и добираешься до result["hits"]["hits"]

// Типично, это делается через json.NewDecoder(res.Body).Decode(&result)
// и дальше обращение через мапы или структурные типы.


// 📌 Что тебе нужно:
// Выполнить search-запрос в Elasticsearch.

// Прочитать ответ и извлечь:

// Общее количество результатов (hits.total.value).

// Массив документов (hits.hits → _source).

// Преобразовать каждый _source в types.Place.

// Вернуть []types.Place, общее количество и ошибку.



package db

import "Go_Day03-1/src/types"

type Store interface {
	// возвращает список записей, общее количество найденных записей и (или) ошибку
	GetPlaces(limit int, offset int) ([]types.Place, int, error)

	// принимает:
	// limit - какое количество записей выводить на странице
	// offset - после какого id выводить записи из базы в количестве limit

	// возвращает:
	// []types.Place - слайс записей в количестве limit
	// int - общее количество записей с маркером "places"
	// или ошибку


}



func (es *ElasticClient) GetPlaces(limit int, offset int) ([]types.Place, int, error) {

	Println("OK")

	return []Place, 0, nil

	// тут надо написать много кода как я вытаскиваю данные из эластика
	// это и есть основная функция, реализующая доступ к эластик через интерфейс Store

}

