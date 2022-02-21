package errs

import (
	"errors"
	"log"
	"net/http"
)

const (
	NoType = ErrorType(iota)
	BadRequest
	NotFound
)

type ErrorType uint

type customError struct {
	errorType     ErrorType
	originalError error
	contextInfo   string
}

// Error возвращат текст сообщения customError
func (error customError) Error() string {
	return error.originalError.Error()
}

// New создает новый customError с указанным ErrorType
func (t ErrorType) New(msg string) error {
	return customError{
		errorType:     t,
		originalError: errors.New(msg),
	}
}

// New создает новый customError с не указанным ErrorType
func New(msg string) error {
	return customError{
		errorType:     NoType,
		originalError: errors.New(msg),
	}

}

// AddErrorContext возвращает customError с добавленным контекстным сообщением
func AddErrorContext(err error, msg string) error {
	// проверяем является ли err экземпляром структуры customError
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			originalError: customErr.originalError,
			contextInfo:   msg,
		}
	}
	// если не является, возвращаем новый экземпляр
	return customError{
		errorType:     NoType,
		originalError: err,
		contextInfo:   msg,
	}
}

// GetType возвращает тип ошибки
func GetType(err error) ErrorType {
	// проверяем является ли err экземпляром структуры customError
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}
	// если не является, возвращаем дефольный тип NoType
	return NoType
}

// GetContext возвращает контекст ошибки
func GetContext(err error) string {
	// проверяем является ли err экземпляром структуры customError
	if customErr, ok := err.(customError); ok {
		return customErr.contextInfo
	}
	// если не является, возвращаем новый экземпляр
	return ""
}

// HandleError обрабатывает ошибки
func HandleError(w http.ResponseWriter, err error) {
	var status int

	// в зависимости от типа ошибки выбираем статус код ответа
	errorType := GetType(err)
	switch errorType {
	case BadRequest:
		status = http.StatusBadRequest
	case NotFound:
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	// чтобы не засорять логи, логируем только внутренние ошибки сервера,
	// которые имеют тип ошибки по умолчанию: NoType
	if errorType == NoType {
		log.Print(err)
	}

	// если есть контекст, то отправляем его в теле сообщения об ошибке
	errorContext := GetContext(err)
	if errorContext != "" {
		http.Error(w, errorContext, status)
		return
	}

	// если нет контекста, то отправляем готовый текст статус кода
	http.Error(w, http.StatusText(status), status)

}
