//go:build manual

package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
)

// тип структурного объекта с полем id из CSV файла
type RestaurantsMANUAL struct {
	ID uint64 `json:"id"`
	RestaurantsBASE
}

type InputType = []string // тип данных для загрузки из CSV в канал-буфер

type Restaurants = RestaurantsMANUAL

// Чтение строк из CSV и заливка в канал-буфер
func CSVLinesToChannel(file *os.File, ch chan InputType) {

	reader := initCSVReader(file) // инициализируем и настраиваем ридер для правильной разбивки строк на сегменты информации для конвертации в объекты Restaurants

	// Пропускаем заголовок
	if _, err := reader.Read(); err != nil {
		log.Fatalf("Ошибка чтения заголовка: %v", err)
	}

	go func() {
		defer close(ch)
		for {
			line, err := reader.Read()
			if err == io.EOF { // читаем файл до конца и выходим
				break
			}
			if err != nil {
				log.Printf("Ошибка чтения строки: %v", err)
				continue
			}

			if len(line) < 6 { // полей должно быть не менее шести: id, name, address, phone, longitude, latitude
				log.Printf("\n⚠️  Пропуск строки с недостаточным количеством полей: \n%+v\n ❗ полей должно быть не менее шести: id, name, address, phone, longitude, latitude", line)

				continue

			}
			ch <- line
		}
	}()
}

// Конвертация строки в объект Restaurants при ручном парсинге
func ConvertLineToRestaurantsOBJ(line InputType) (Restaurants, error) {

	// Парсим ID
	id, err := strconv.ParseUint(line[0], 10, 64)
	if err != nil {
		log.Printf("Ошибка парсинга ID: %v", err)
	}

	// Парсим координаты
	lat, err := strconv.ParseFloat(line[5], 64)
	if err != nil {
		log.Printf("Ошибка парсинга Latitude: %v", err)
	}

	lon, err := strconv.ParseFloat(line[4], 64)
	if err != nil {
		log.Printf("Ошибка парсинга Longitude: %v", err)
	}

	if err != nil {
		return Restaurants{}, err
	}

	return Restaurants{
		ID: id,
		RestaurantsBASE: RestaurantsBASE{
			Name:    line[1],
			Address: line[2],
			Phone:   line[3],
			Location: Location{
				Latitude:  lat,
				Longitude: lon,
			}},
	}, nil

}

func GetValue(r InputType) []byte {

	objRestaurant, _ := ConvertLineToRestaurantsOBJ(r)
	val, _ := json.Marshal(objRestaurant)

	return val
}
