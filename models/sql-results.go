package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type SQLResult struct {
	result   sql.Result
	lastID   int64
	affected int64
	errors   []error
}

// ExtractStats извлекает статистику выполнения операции
func (r *SQLResult) ExtractStats() {
	if r.result == nil {
		r.errors = append(r.errors, fmt.Errorf("sql.Result is nil"))
		return
	}

	var err error

	// Получаем LastInsertId с обработкой ошибок
	if r.lastID, err = r.result.LastInsertId(); err != nil {
		r.errors = append(r.errors, fmt.Errorf("last insert id error: %w", err))
	}

	// Получаем RowsAffected с обработкой ошибок
	if r.affected, err = r.result.RowsAffected(); err != nil {
		r.errors = append(r.errors, fmt.Errorf("rows affected error: %w", err))
	}
}

// ToJSON генерирует JSON-ответ со статистикой
func (r *SQLResult) ToJSON() ([]byte, error) {
	// Вспомогательная структура для безопасной сериализации
	type response struct {
		LastInsertID int64    `json:"last_insert_id"`
		RowsAffected int64    `json:"rows_affected"`
		Errors       []string `json:"errors,omitempty"`
	}

	res := response{
		LastInsertID: r.lastID,
		RowsAffected: r.affected,
		Errors:       make([]string, 0, len(r.errors)),
	}

	// Конвертируем ошибки в строки
	for _, err := range r.errors {
		res.Errors = append(res.Errors, err.Error())
	}

	return json.MarshalIndent(res, "", "  ")
}

// Дополнительные методы для безопасного доступа к данным
func (r *SQLResult) HasErrors() bool {
	return len(r.errors) > 0
}

func (r *SQLResult) LastInsertID() int64 {
	return r.lastID
}

func (r *SQLResult) RowsAffected() int64 {
	return r.affected
}
