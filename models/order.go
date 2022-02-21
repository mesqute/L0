package models

import (
	"L0/utilites/errs"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type Order struct {
	OrderUid          string   `json:"order_uid" validate:"required"`
	TrackNumber       string   `json:"track_number" validate:"required"`
	Entry             string   `json:"entry" validate:"required"`
	Delivery          Delivery `json:"delivery" validate:"required"`
	Payment           Payment  `json:"payment" validate:"required"`
	Items             []Item   `json:"items" validate:"required"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id" validate:"required"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SmId              int      `json:"sm_id" validate:"numeric"`
	DateCreated       string   `json:"date_created" validate:"required"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount" validate:"numeric"`
	PaymentDt    int    `json:"payment_dt" validate:"numeric"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost" validate:"numeric"`
	GoodsTotal   int    `json:"goods_total" validate:"numeric"`
	CustomFee    int    `json:"custom_fee" validate:"numeric"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id" validate:"required,numeric"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price" validate:"numeric"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale" validate:"numeric"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price" validate:"numeric"`
	NmId        int    `json:"nm_id" validate:"numeric"`
	Brand       string `json:"brand"`
	Status      int    `json:"status" validate:"required,numeric"`
}

var Validator = validator.New()

// Validate проверяет поля структуры на корректность
func (o Order) Validate() error {
	var errorContext []string

	err := Validator.Struct(o)
	if err != nil {
		for _, fieldError := range err.(validator.ValidationErrors) {
			errString := "Поле: " + fieldError.Field() + ", Тег ошибки: " + fieldError.Tag()
			errorContext = append(errorContext, errString)
		}
	}

	if len(errorContext) != 0 {
		err := errs.BadRequest.New("[Order.Validate] validation failed")
		err = errs.AddErrorContext(err, strings.Join(errorContext, "\n"))
		return err
	}
	return nil
}

// GetFieldPointers возвращает срез содержащий указатели на поля структуры
func GetFieldPointers(strct interface{}) ([]interface{}, error) {

	// получаем из входных данных структуру reflect.Value
	s := reflect.ValueOf(strct).Elem()

	// проверяем является ли полученная переменная структурой
	if s.Kind().String() != "struct" {
		err := errs.New("[GetFieldPointers] input value is not struct")
		return nil, err
	}

	// считываем количество полей структуры Order
	num := s.NumField()
	// инициализируем срез для записи полей
	var fields []interface{}

	for i := 0; i < num; i++ {
		field := s.Field(i)
		// проверка можно ли получить значение поля
		// (не хранит указатели, функции, структуры и т.д.)
		if field.CanInterface() && field.CanAddr() &&
			field.Kind().String() != "struct" &&
			field.Kind().String() != "map" &&
			field.Kind().String() != "slice" {
			// добавляем в срез вывода значение поля структуры
			fields = append(fields, field.Addr().Interface())
		}
	}

	return fields, nil

}

// GetFieldValues возвращает срез содержащий значения полей структуры
func GetFieldValues(strct interface{}) ([]interface{}, error) {

	// получаем из структуры Order структуру reflect.Value
	s := reflect.ValueOf(strct).Elem()

	// проверяем является ли полученная переменная структурой
	if s.Kind().String() != "struct" {
		err := errs.New("[GetFieldPointers] input value is not struct")
		return nil, err
	}

	// считываем количество полей структуры Order
	num := s.NumField()
	// инициализируем срез для записи полей
	var fields []interface{}

	for i := 0; i < num; i++ {
		field := s.Field(i)
		// проверка можно ли получить значение поля
		// (не хранит указатели, функции, структуры и т.д.)
		if field.CanInterface() &&
			field.Kind().String() != "struct" &&
			field.Kind().String() != "map" &&
			field.Kind().String() != "slice" {
			// добавляем в срез вывода значение поля структуры
			fields = append(fields, field.Interface())
		}
	}

	return fields, nil
}
