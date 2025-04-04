package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

// type Industry struct {
// 	CompanyName                string `csv:"Organization Name"`
// 	LinkedIn                   string `csv:"LinkedIn"`
// 	Website                    string `csv:"Website"`
// 	TotalFundingAmount         int    `csv:"Total Funding Amount"`
// 	TotalFundingAmountCurrency string `csv:"Total Funding Amount Currency"`
// 	HeadquartersLocation       string `csv:"Headquarters Location"`
// }

// 5.57ms -> 600 records (read)
func main() {
	now := time.Now()
	// readChannel := make(chan Industry, 1)
	readChannel := make(chan RestaurantsCSV, 25)

	// readFilePath := "process.csv"
	readFilePath := "../materials/data.csv"

	// Open the CSV readFile
	readFile, err := os.OpenFile(readFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer readFile.Close()

	count := 0
	readFromCSV(readFile, readChannel)

	// Print the records
	for r := range readChannel {
		fmt.Println("========================================")
		fmt.Println(r)
		fmt.Println("========================================")
		fmt.Println()

		count++
	}

	fmt.Println(time.Since(now), count)
}

func readFromCSV(file *os.File, c chan RestaurantsCSV) {

	// Создаём новый CSV-ридер с настроенными параметрами
	reader := csv.NewReader(file)
	reader.Comma = '\t'         // Указываем разделитель (табуляция)
	reader.LazyQuotes = true    // Разрешаем необычные кавычки
	reader.FieldsPerRecord = -1 // Не проверяем количество полей

	// Запускаем асинхронное чтение и отправку в канал
	go func() {
		// Пытаемся распарсить весь файл в канал
		err := gocsv.UnmarshalToChan(file, c)
		if err != nil {
			fmt.Println("Ошибка при чтении CSV:", err)
			close(c) // Закрываем канал, если произошла ошибка
		}
	}()
}
