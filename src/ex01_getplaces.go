package main

// import (
// 	"Go_Day03-1/src/db"
// 	"Go_Day03-1/src/types"
// )

// func (r Restaurants) GetPlaces(limit int, offset int) ([]types.Place, int, error) {

// }


// 1. Получать ?page=N через r.URL.Query().Get("page").

// 2. Конвертировать его в int.

// 3. Вычислять, с какой страницы начинать, и передавать в search_after.

// 4. Полученные записи из elastic передавать в Page и рендерить шаблон.