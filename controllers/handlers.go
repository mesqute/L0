package controllers

import (
	"L0/models"
	u "L0/utilites"
	"L0/utilites/errs"
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
