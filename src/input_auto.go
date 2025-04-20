//go:build !manual

package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/gocarina/gocsv"
)

// вспомогательная структура для автоматического парсинга
type RestaurantsCSV struct {
	Name      string  `csv:"Name"`
	Address   string  `csv:"Address"`
	Phone     string  `csv:"Phone"`
	Latitude  float64 `csv:"Latitude"`
	Longitude float64 `csv:"Longitude"`
}

type InputType = RestaurantsCSV

type Restaurants = RestaurantsBASE

// вспомогательный метод для автоматического парсинга
func (r RestaurantsCSV) ToRestaurants() Restaurants {
	return Restaurants{
		Name:    r.Name,
		Address: r.Address,
		Phone:   r.Phone,
		Location: Location{
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
		},
	}
}

// функция читает данные из CSV файла и записывает в канал. Принимает на вход файл и канал типа RestaurantsCSV
func CSVLinesToChannel(file *os.File, c chan InputType) {

	// Устанавливаем кастомный CSVReader, который будет читать после заголовка
	gocsv.SetCSVReader(func(_ io.Reader) gocsv.CSVReader {
		return initCSVReader(file) // вызов инициализатора для CSV-reader-а
	})

	// Создаём буфер для файла, хотя это не обязательно
	bufReader := bufio.NewReader(file)

	// Стартуем горутину с потоком данных в канал
	go func() {
		// defer close(c) // в данном случае канал закрывать не надо, так как он автоматически закрывается gocsv
		if err := gocsv.UnmarshalToChan(bufReader, c); err != nil {
			log.Fatalf("Ошибка при анмаршалинге CSV: %v", err)
		}
	}()
}

// функция принимает на вход объект RestaurantsCSV, а возвращает его в виде json типа []byte
func GetValue(r InputType) []byte {
	val, _ := json.Marshal(r.ToRestaurants()) // сначала объект конвертируется из RestaurantsCSV в RestaurantsBASE,
	// а потом он декодируется в json
	return val
}
