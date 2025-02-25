package models

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

func fillNilValues(entry Entry, col string, columnType string) {
	var emptyValues = map[string]interface{}{
		"INT":              0,
		"INTEGER":          0,
		"SMALLINT":         0,
		"BIGINT":           0,
		"SERIAL":           0,
		"BIGSERIAL":        0,
		"FLOAT":            0.0,
		"DOUBLE PRECISION": 0.0,
		"REAL":             0.0,
		"NUMERIC":          0.0,
		"DECIMAL":          0.0,
		"BOOLEAN":          false,
		"VARCHAR":          "",
		"TEXT":             "",
		"CHAR":             "",
		"DATE":             time.Time{},
		"TIME":             time.Time{},
		"TIMESTAMP":        time.Time{},
		"TIMESTAMPTZ":      time.Time{},
		"INET":             "",
		"CIDR":             "",
	}

	if _, ok := emptyValues[columnType]; !ok {
		log.Printf("unknown column type for column %s: %s", col, columnType)
		return
	}

	entry[col] = emptyValues[columnType]
}

func convertDatabaseValue(columnType string, val interface{}, col string) (interface{}, error) {
	err := errors.New(fmt.Sprintf("unexpected type for column %s with value's type %T", col, val))

	switch columnType {
	case "INT", "INTEGER", "SMALLINT", "BIGINT", "SERIAL", "BIGSERIAL":
		if v, ok := val.(int64); ok {
			return v, nil
		}
		log.Printf("unexpected type for column %s: expected int64, got %T", col, val)
		return 0, nil // Значение по умолчанию для целых чисел

	case "FLOAT", "DOUBLE PRECISION", "REAL", "NUMERIC", "DECIMAL":
		if v, ok := val.(float64); ok {
			return v, nil
		}
		log.Printf("unexpected type for column %s: expected float64, got %T", col, val)
		return 0.0, nil // Значение по умолчанию для чисел с плавающей запятой

	case "BOOLEAN":
		if v, ok := val.(bool); ok {
			return v, nil
		}
		log.Printf("unexpected type for column %s: expected bool, got %T", col, val)
		return false, nil // Значение по умолчанию для булевых значений

	case "VARCHAR", "TEXT", "CHAR":
		if v, ok := val.([]byte); ok {
			return string(v), nil
		}
		log.Printf("unexpected type for column %s: expected []byte, got %T", col, val)
		return "", nil // Значение по умолчанию для строк

	case "DATE", "TIME", "TIMESTAMP", "TIMESTAMPTZ":
		if v, ok := val.(time.Time); ok {
			return v.Format(time.RFC3339), nil
		}
		log.Printf("unexpected type for column %s: expected time.Time, got %T", col, val)
		return time.Time{}.Format(time.RFC3339), nil // Значение по умолчанию для времени

	case "INET", "CIDR":
		if v, ok := val.([]byte); ok {
			return net.IP(v).String(), nil
		}
		log.Printf("unexpected type for column %s: expected []byte, got %T", col, val)
		return "", nil // Значение по умолчанию для IP-адресов

	default:
		if v, ok := val.([]byte); ok {
			return string(v), nil
		}
		log.Printf("unexpected type for column %s: expected []byte, got %T", col, val)
		return "", nil // Значение по умолчанию для неизвестных типов
	}
}
