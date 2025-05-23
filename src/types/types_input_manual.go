//go:build manual

package types

// тип структурного объекта с полем id из CSV файла
type RestaurantsMANUAL struct {
	ID uint64 `json:"id"`
	RestaurantsBASE
}

type InputType = []string // тип данных для загрузки из CSV в канал-буфер

type Restaurants = RestaurantsMANUAL
