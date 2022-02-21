package models

import (
	"L0/utilites/errs"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"sync"
)

var (
	mtx              sync.Mutex
	postgresInstance *Postgres = nil
)

type Postgres struct {
	db *sql.DB
}

// initPostgresInstance инициализирует подключение к БД,
// а также новый экземпляр структуры Postgres
func initPostgresInstance() {
	mtx.Lock()
	defer mtx.Unlock()
	if postgresInstance == nil {

		// объявляем указатель на новый экземпляр структуры Postgres
		postgres := new(Postgres)

		// формируем строку подключения и открываем БД для подключений
		//connStr := "user=samurai host=localhost dbname=l0 password=0000 sslmode=disable"
		connStr := "user=samurai host=host.docker.internal dbname=l0 password=0000 sslmode=disable"
		base, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}

		// проверяем подключение к БД
		err = base.Ping()
		if err != nil {
			log.Fatal(err)
		}

		// добавляем указатель на БД в структуту Postgres
		// и передаем указатель на структуру в глобальную переменную
		postgres.db = base
		postgresInstance = postgres
	}
}

// fillStruct заполняет поля структуры данными
func (p *Postgres) fillStruct(strct interface{}, queryString string, args ...interface{}) error {
	// получаем срез указателей на поля структуры
	pointers, err := GetFieldPointers(strct)
	if err != nil {
		return err
	}
	// обрабатываем ответ на запрос и заполняем поля структуры полученными данными
	row := p.db.QueryRow(queryString, args...)
	if err := row.Scan(pointers...); err != nil {
		err = errs.New("[Postgres.fillStruct] " + err.Error())
		return err
	}
	return nil
}

// sendStruct отправляет в БД данные из полей структуры.
// Если структура зависимая, то добавляет Id связанной структуры.
func (p *Postgres) sendStruct(strct interface{}, queryString string, isMainStruct bool, mainId string) error {
	// получаем срез, содержащий данные всех полей структуры
	values, err := GetFieldValues(strct)
	if err != nil {
		return err
	}

	// если структура главная, выполняем запрос на добавление данных из среза в БД
	if isMainStruct {
		// дополняем полученную строку запроса плейсхолдерами параметров ($1, $2 и т.д)
		var phs []string
		for i := 1; i <= len(values); i++ {
			phs = append(phs, "$"+strconv.Itoa(i))
		}
		phStr := "(" + strings.Join(phs, ", ") + ")"
		queryString = queryString + phStr

		if _, err := p.db.Exec(queryString, values...); err != nil {
			err = errs.New("[Postgres.sendStruct] " + err.Error())
			return err
		}
		return nil
	}

	// если структура зависимая, то добавляем в конец среза Id связанной структуры,
	// а затем выполняем запрос на добавление данных из среза в БД
	values = append(values, mainId)

	// дополняем полученную строку запроса плейсхолдерами параметров ($1, $2 и т.д)
	var phs []string
	for i := 1; i <= len(values); i++ {
		phs = append(phs, "$"+strconv.Itoa(i))
	}
	phStr := "(" + strings.Join(phs, ", ") + ")"
	queryString = queryString + phStr

	if _, err := p.db.Exec(queryString, values...); err != nil {
		err = errs.New("[Postgres.sendStruct] " + err.Error())
		return err
	}
	return nil
}

// GetPostgres возвращает указатель на структуру БД Postgres,
// и если единственный экземпляр структуры не инициализирован,
// то инициализирует его
func GetPostgres() *Postgres {
	if postgresInstance == nil {
		initPostgresInstance()
	}
	return postgresInstance
}

// GetData возвращает сохрененные в БД данные
func (p *Postgres) GetData() (map[string]Order, error) {

	orders := make(map[string]Order)

	// получение из БД всех записей структуры Order
	ordersRows, err := p.db.Query("SELECT * FROM get_orders()")
	if err != nil {
		err = errs.New("[Postgres.GetData] " + err.Error())
		return nil, err
	}

	// обработка полученных из БД записей
	for ordersRows.Next() {

		// получение из БД данных структуры Order
		var order Order
		// получаем срез указателей на поля структуры
		orderPointers, err := GetFieldPointers(&order)
		if err != nil {
			return nil, err
		}
		// заполняем поля структуры полученными данными
		if err := ordersRows.Scan(orderPointers...); err != nil {
			err = errs.New("[Postgres.GetData] " + err.Error())
			return nil, err
		}
		// получаем id главной структуры Order
		orderId := order.OrderUid

		// получение из БД данных подструктуры Delivery
		var delivery Delivery
		if err := p.fillStruct(&delivery, "SELECT * FROM get_orders_delivery($1)", orderId); err != nil {
			return nil, err
		}

		// получение из БД данных подструктуры Payment
		var payment Payment
		if err := p.fillStruct(&payment, "SELECT * FROM get_orders_payment($1)", orderId); err != nil {
			return nil, err
		}

		// получение из БД данных подструктуры Items
		var items []Item
		itemsRows, err := p.db.Query("SELECT * FROM get_orders_items($1)", orderId)
		if err != nil {
			return nil, err
		}
		// обработка всех полученных строк и запись данных в Items
		for itemsRows.Next() {
			var item Item
			itemPointers, err := GetFieldPointers(&item)
			if err != nil {
				return nil, err
			}
			if err := itemsRows.Scan(itemPointers...); err != nil {
				return nil, err
			}
			items = append(items, item)
		}

		// добавление в структуру Order всех собранных подструктур
		order.Delivery = delivery
		order.Payment = payment
		order.Items = items

		id := order.OrderUid
		orders[id] = order
	}

	return orders, nil
}

// Insert сохраняет в БД данные
func (p *Postgres) Insert(order Order) error {

	// проверка есть ли в БД данные с таким же id
	orderId := order.OrderUid
	var check bool
	err := p.db.QueryRow("SELECT check_order($1)", orderId).Scan(&check)
	if check {
		err := errs.BadRequest.New("[Postgres.Insert] data already exist")
		err = errs.AddErrorContext(err, "Данные с таким id уже существуют")
		return err
	}
	// обработка ошибок
	if err != nil {
		err = errs.New("[Postgres.Insert] " + err.Error())
		return err
	}

	// сохранение в БД структуры Order
	if err := p.sendStruct(&order, "CALL insert_order", true, ""); err != nil {
		return err
	}

	// сохранение в БД подструктуры Delivery
	if err := p.sendStruct(&order.Delivery, "CALL insert_orders_delivery", false, order.OrderUid); err != nil {
		return err
	}

	// сохранение в БД подструктуры Payment
	if err := p.sendStruct(&order.Payment, "CALL insert_orders_payment", false, order.OrderUid); err != nil {
		return err
	}

	// сохранение в БД всех элементов подструктуры Items
	for _, item := range order.Items {
		if err := p.sendStruct(&item, "CALL insert_orders_item", false, order.OrderUid); err != nil {
			return err
		}
	}

	return nil
}
