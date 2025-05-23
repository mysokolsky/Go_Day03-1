//go:build !manual

package types

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
