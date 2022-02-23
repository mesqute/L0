package models

import (
	"L0/utilites/errs"
	"sync"
)

var (
	mtxInstance sync.Mutex
	instance    *Cache = nil
)

// Cache представляет собой in-memory кэш
type Cache struct {
	mtx  sync.RWMutex
	data map[string]Order
}

// инициализирует единственный экземпляр структуры Cache
func initCacheInstance() {
	mtxInstance.Lock()
	defer mtxInstance.Unlock()
	if instance == nil {
		cache := new(Cache)
		cache.data = make(map[string]Order)
		instance = cache
	}
}

// GetCacheInstance возвращает указатель на структуру Cache
func GetCacheInstance() *Cache {
	if instance == nil {
		initCacheInstance()
	}
	return instance
}

// SetCacheData передает в Кэш новый датасет
func (c *Cache) SetCacheData(data map[string]Order) {
	c.data = data
}

// Get возвращает данные из кэша по id.
// Если не находит данные, то возвращает ошибку.
func (c *Cache) Get(id string) (Order, error) {
	//проверем, инициализирован ли кеш
	if c.data == nil {
		err := errs.New("[Cache.Get] кеш не инициализирован")
		return Order{}, err
	}

	// используем блокировку для чтения (не блокирует чтение для остальных, но блокирует запись)
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	// считываем данные и их наличие
	val, ok := c.data[id]
	// если данных нет, то возвращаем ошибку
	if !ok {
		err := errs.NotFound.New("[Cache.Get] data not found")
		err = errs.AddErrorContext(err, "Объект с таким id не найден")
		return Order{}, err
	}
	// если данные есть, то возвращаем их
	return val, nil
}

func (c *Cache) Insert(data Order) error {
	id := data.OrderUid

	//проверем, инициализирован ли кеш
	if c.data == nil {
		err := errs.New("[Cache.Get] кеш не инициализирован")
		return err
	}

	// проверка, существуют ли в кэше данные с таким же id.
	// если существуют, то возвращаем ошибку
	if _, ok := c.data[id]; ok {
		err := errs.BadRequest.New("[Cache.Insert]: data already exist")
		err = errs.AddErrorContext(err, "Данные с таким id уже существуют")
		return err
	}

	// используем полную блокировку для записи
	c.mtx.Lock()
	defer c.mtx.Unlock()

	// записываем данные в кэш
	c.data[id] = data

	return nil
}
