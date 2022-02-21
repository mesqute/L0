package models

import (
	"log"
)

var (
	database Database
	cache    *Cache
)

// инициализирует кеш и БД
func init() {
	// инициализация БД и кеша
	database = GetPostgres()
	cache = GetCacheInstance()

	// восстанавливаем датасет из БД
	dataset, err := database.GetData()
	if err != nil {
		log.Fatal(err)
	}
	// передаем полученный датасет в кеш
	cache.SetCacheData(dataset)
}

// InsertData добавляет указанные данные в память сервиса (в БД и в кеш)
func InsertData(data Order) error {

	// добавляем данные в БД
	if err := database.Insert(data); err != nil {
		return err
	}

	// добавляем данные в кеш
	if err := cache.Insert(data); err != nil {
		return err
	}
	return nil
}

// GetData возвращает указанные данные из памяти сервиса (из кеша)
func GetData(id string) (Order, error) {
	// берем данные из кеша
	data, err := cache.Get(id)
	if err != nil {
		return Order{}, err
	}
	return data, nil
}
