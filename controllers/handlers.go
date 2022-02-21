package controllers

import (
	"L0/models"
	u "L0/utilites"
	"L0/utilites/errs"
	"encoding/json"
	"io"
	"net/http"
	"text/template"
)

// home обрабатывает запросы на отображение интерфейса
func home(w http.ResponseWriter, r *http.Request) {
	// проверка на соответствие обрабатываемому методу Get
	if r.Method != http.MethodGet {
		u.ErrorMethodNotAllowed(w, http.MethodGet)
		return
	}
	// добавляем файл главной страницы в обработчик
	t, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		err = errs.New("[home] " + err.Error())
		errs.HandleError(w, err)
		return
	}
	// возвращаем обработанный файл главной страницы
	err = t.Execute(w, nil)
	if err != nil {
		err = errs.New("[home] " + err.Error())
		errs.HandleError(w, err)
		return
	}
}

// setData обрабатывает запросы на добавление записей
func setData(w http.ResponseWriter, r *http.Request) {
	// проверка на соответствие обрабатываемому методу Post
	if r.Method != http.MethodPost {
		u.ErrorMethodNotAllowed(w, http.MethodPost)
		return
	}
	// читаем данные из Body в байтовый срез
	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = errs.New("[ReadBody] " + err.Error())
		errs.HandleError(w, err)
		return
	}

	// инициализация экземпляра структуры Order
	// в которую будут записаны полученные в запросе данные
	var order models.Order

	// декодирование полученного json файла в ранее созданый экземпяр структуры Order
	if err := json.Unmarshal(b, &order); err != nil {
		err = errs.BadRequest.New("[json.Decode] " + err.Error())
		errs.HandleError(w, err)
		return
	}

	// проверка полученных данных на корректность
	if err := order.Validate(); err != nil {
		errs.HandleError(w, err)
		return
	}

	// добавление полученных данных в память сервиса
	if err := models.InsertData(order); err != nil {
		errs.HandleError(w, err)
		return
	}

	// отправка ответа с информацией об успешном добавлении данных
	w.WriteHeader(http.StatusCreated)
}

// getData обрабатывает запросы на выдачу записей по id
func getData(w http.ResponseWriter, r *http.Request) {
	// проверка на соответствие обрабатываемому методу Get
	if r.Method != http.MethodGet {
		u.ErrorMethodNotAllowed(w, http.MethodGet)
		return
	}
	// считывание параметра id из URL
	id := r.URL.Query().Get("id")
	if id == "" {
		err := errs.BadRequest.New("request dont have id")
		errs.HandleError(w, err)
		return
	}

	// получаем данные
	data, err := models.GetData(id)
	if err != nil {
		errs.HandleError(w, err)
		return
	}

	// отправляем данные
	err = u.Respond(w, http.StatusOK, data)
	if err != nil {
		errs.HandleError(w, err)
	}
}
