package errs

import "errors"

var (
	// Основные ошибки бизнес-логики
	ErrNotFound      = errors.New("subscription not found")
	ErrAlreadyExists = errors.New("subscription already exists")
	ErrInvalidData   = errors.New("invalid data provided")

	// Ошибки валидации
	ErrInvalidDateFormat = errors.New("invalid date format, expected MM-YYYY")
	ErrInvalidUUID       = errors.New("invalid UUID format")
	ErrInvalidDateRange  = errors.New("invalid date range: start date must be before or equal to end date")
	ErrInvalidPrice      = errors.New("price must be a positive integer")
	ErrInvalidPagination = errors.New("invalid pagination parameters")

	// Ошибки доступа к данным
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery      = errors.New("database query error")

	// Ошибки обработки запросов
	ErrInvalidJSON          = errors.New("invalid JSON format")
	ErrMissingRequiredField = errors.New("missing required field")
)
