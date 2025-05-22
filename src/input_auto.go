//go:build !manual

package main

import (
	// "bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"Go_Day03-1/src/types"
	"github.com/gocarina/gocsv"
)

// функция читает данные из CSV файла и записывает в канал. Принимает на вход файл и канал типа RestaurantsCSV
func CSVLinesToChannel(file *os.File, c chan types.InputType) {

	// Создаём буфер для файла, хотя это не обязательно
	// bufReader := bufio.NewReader(file)

	// Устанавливаем кастомный CSVReader, который будет читать после заголовка
	gocsv.SetCSVReader(func(_ io.Reader) gocsv.CSVReader {
		return initCSVReader(file) // вызов инициализатора для CSV-reader-а
	})

	// Стартуем горутину с потоком данных в канал
	go func() {
		// defer close(c) // в данном случае канал закрывать не надо, так как он автоматически закрывается gocsv
		if err := gocsv.UnmarshalToChan(file, c); err != nil {
			log.Fatalf("Ошибка при анмаршалинге CSV: %v", err)
		}
	}()
}

// функция принимает на вход объект RestaurantsCSV, а возвращает его в виде json типа []byte
func GetValue(r types.InputType) []byte {
	val, _ := json.Marshal(r.ToRestaurants()) // сначала объект конвертируется из RestaurantsCSV в RestaurantsBASE,
	// а потом он декодируется в json
	return val
}
