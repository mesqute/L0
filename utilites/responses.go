package utilites

import (
	"L0/utilites/errs"
	"encoding/json"
	"net/http"
	"strings"
)

// ErrorMethodNotAllowed отправляет код 405 с описанием доступных методов
func ErrorMethodNotAllowed(w http.ResponseWriter, allowMethods ...string) {

	// объединение в строку всех методов переданных в параметрах функции
	allowMethodsString := strings.Join(allowMethods, ", ")

	// передача в заголовок ответа список доступных методов
	w.Header().Set("Allow", allowMethodsString)

	// отправка ответа с кодом и описанием ошибки
	http.Error(w, "Метод запрещен", http.StatusMethodNotAllowed)
}

// Respond формирует и отправляет ответ в формате JSON
func Respond(w http.ResponseWriter, status int, data interface{}) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		err = errs.New("[Respond] " + err.Error())
		return err
	}
	return nil
}
