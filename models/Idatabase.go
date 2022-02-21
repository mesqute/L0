package models

// Database содержит методы, при помощи которых элементы сервиса взаимодействуют с БД
type Database interface {
	// Insert сохраняет в БД данные
	Insert(order Order) error
	// GetData возвращает сохрененные в БД данные, соответствующие заданному id
	GetData() (map[string]Order, error)
}
