package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

// SQLRows is used to store the results of SELECT queries
type SQLRows struct {
	rows    *sql.Rows
	columns []string
	types   []string
	data    []map[string]any
	errors  []error
}

func (r *SQLRows) ParseMeta() error {
	if r.rows == nil {
		err := fmt.Errorf("parsing metadata: *sql.Rows is empty")
		return err
	}

	columns, err := r.rows.Columns()
	if err != nil {
		r.errors = append(r.errors, fmt.Errorf("failed to get columns: %w", err))
		return fmt.Errorf("failed to get coumns: %v", err)
	}

	r.columns = make([]string, len(columns))
	copy(r.columns, columns)

	types, err := r.rows.ColumnTypes()
	if err != nil {
		r.errors = append(r.errors, fmt.Errorf("failed to get columns types: %w", err))
		return fmt.Errorf("failed to get column types: %v", err)
	}
	r.types = make([]string, len(types))
	for i, typ := range types {
		r.types[i] = typ.DatabaseTypeName()
	}

	return nil
}

func (r *SQLRows) ScanBatch() {
	defer r.rows.Close()

	for r.rows.Next() {
		rowMap := make(map[string]any)
		typedRow := make(map[string]any)

		values := make([]any, len(r.columns))
		ptrs := make([]any, len(values))
		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := r.rows.Scan(ptrs...); err != nil {
			r.errors = append(r.errors, fmt.Errorf("scan failed: %w", err))
			continue
		}

		for i, col := range r.columns {
			// Получаем тип PostgreSQL для текущей колонки
			pgType := r.types[i]
			value := values[i]

			// Приводим значение к нужному типу
			converted, err := convertValue(pgType, value)
			if err != nil {
				r.errors = append(r.errors, fmt.Errorf("conversion failed for column %s: %w", col, err))
				continue
			}

			typedRow[col] = converted
			rowMap[col] = value // Оригинальное значение, если нужно
		}

		r.data = append(r.data, typedRow) // Используем приведенные значения
	}

	if err := r.rows.Err(); err != nil {
		r.errors = append(r.errors, fmt.Errorf("rows error: %w", err))
	}
}

// Функция преобразования значений
func convertValue(pgType string, value any) (any, error) {
	if value == nil {
		return nil, nil
	}

	switch pgType {
	case "INT", "INTEGER", "SERIAL", "SMALLINT":
		return convertInt(value)
	case "BIGINT", "BIGSERIAL":
		return convertBigInt(value)
	// case "FLOAT", "REAL":
	// 	return convertFloat32(value)
	case "DOUBLE PRECISION", "FLOAT", "REAL":
		return convertFloat64(value)
	case "NUMERIC", "DECIMAL", "MONEY":
		return convertNumeric(value)
	case "BOOLEAN":
		return convertBool(value)
	case "VARCHAR", "TEXT", "CHAR":
		return convertString(value)
	case "DATE", "TIME", "TIMESTAMP", "TIMESTAMPTZ":
		return convertTime(value)
	case "INET":
		return convertIP(value)
	case "CIDR":
		return convertIPNet(value)
	case "BYTEA":
		return convertBytea(value)
	default:
		return nil, fmt.Errorf("unsupported type: %s", pgType)
	}
}

// Реализации конкретных преобразований
func convertInt(v any) (int, error) {
	switch val := v.(type) {
	case int64:
		return int(val), nil
	case int32:
		return int(val), nil
	case int16:
		return int(val), nil
	case int:
		return val, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

func convertBigInt(v any) (int64, error) {
	switch val := v.(type) {
	case int64:
		return int64(val), nil
	case int32:
		return int64(val), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

func convertFloat64(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

func convertNumeric(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	case []uint8: // Для строковых представлений
		str := string(val)
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse numeric: %v", err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to numeric", v)
	}
}

func convertBool(v any) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case int64:
		return val != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

func convertString(v any) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case []uint8:
		return string(val), nil
	default:
		return "", fmt.Errorf("cannot convert %T to string", v)
	}
}

func convertTime(v any) (time.Time, error) {
	switch val := v.(type) {
	case time.Time:
		return val, nil
	case string:
		return time.Parse(time.RFC3339, val)
	default:
		return time.Time{}, fmt.Errorf("cannot convert %T to time.Time", v)
	}
}

func convertIP(v any) (net.IP, error) {
	switch val := v.(type) {
	case net.IP:
		return val, nil
	case string:
		ip := net.ParseIP(val)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address: %s", val)
		}
		return ip, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to net.IP", v)
	}
}

func convertIPNet(v any) (net.IPNet, error) {
	switch val := v.(type) {
	case net.IPNet:
		return val, nil
	case string:
		_, ipnet, err := net.ParseCIDR(val)
		if err != nil {
			return net.IPNet{}, err
		}
		return *ipnet, nil
	default:
		return net.IPNet{}, fmt.Errorf("cannot convert %T to net.IPNet", v)
	}
}

func convertBytea(v any) ([]byte, error) {
	switch val := v.(type) {
	case []byte:
		return val, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to []byte", v)
	}
}

func (r *SQLRows) ToJSON() ([]byte, error) {
	// Создаем алиас для кастомной сериализации
	type customData map[string]interface{}
	var serializableData []customData

	// Конвертируем данные с учетом специальных типов
	for _, row := range r.data {
		customRow := make(customData)
		for key, value := range row {
			switch v := value.(type) {
			case time.Time:
				customRow[key] = v.Format(time.RFC3339) // Сериализуем время в строку
			case net.IP:
				customRow[key] = v.String() // IP-адрес в строковое представление
			case net.IPNet:
				customRow[key] = v.String() // IP-сеть в строку CIDR
			case []byte:
				customRow[key] = string(v) // Бинарные данные в base64 или строку
			default:
				customRow[key] = value
			}
		}
		serializableData = append(serializableData, customRow)
	}

	return json.MarshalIndent(serializableData, "", "  ")
}
