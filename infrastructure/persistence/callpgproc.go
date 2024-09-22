package persistence

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// Generic function to execute stored procedures with pgx v5
func ExecuteStoredProcedure(conn *pgx.Conn, procedureName string, paramsStruct interface{}, resultStruct interface{}) error {
	// Extract parameters from struct using reflection and pgParam tag
	args, cursorName, err := extractParams(paramsStruct)
	if err != nil {
		return fmt.Errorf("failed to extract parameters: %v", err)
	}

	// Open a transaction with context
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Call the stored procedure and capture the refcursor
	if cursorName == "" {
		cursorName = "ref_cur" // Default cursor name if not provided
	}

	err = tx.QueryRow(context.Background(), fmt.Sprintf("SELECT %s(%s)", procedureName, buildPlaceholders(len(args))), args...).Scan(&cursorName)
	if err != nil {
		return fmt.Errorf("failed to execute stored procedure: %v", err)
	}

	// Fetch all rows from the cursor
	rows, err := tx.Query(context.Background(), fmt.Sprintf("FETCH ALL IN %s", cursorName))
	if err != nil {
		return fmt.Errorf("failed to fetch rows from cursor: %v", err)
	}
	defer rows.Close()

	// Map the rows into resultStruct slice
	err = mapRowsToStruct(rows, resultStruct)
	if err != nil {
		return fmt.Errorf("failed to map rows: %v", err)
	}

	// Commit the transaction
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
func ExecuteStoredProcedureWithCursor[T any](conn *pgx.Conn, procedureName string, paramsStruct interface{}, batchSize int) ([]T, error) {
	args, cursorName, err := extractParams(paramsStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to extract parameters: %v", err)
	}

	if cursorName == "" {
		cursorName = "ref_cur" // Default cursor name if not provided
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Thêm tên cursor vào danh sách tham số
	args = append(args, cursorName)

	// Gọi function với cursor
	_, err = tx.Exec(context.Background(), fmt.Sprintf("SELECT %s(%s)", procedureName, buildPlaceholders(len(args))), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute function: %v", err)
	}

	var results []T

	for {
		// Fetch một batch các rows
		rows, err := tx.Query(context.Background(), fmt.Sprintf("FETCH %d FROM %s", batchSize, cursorName))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch from cursor: %v", err)
		}

		batchResult, err := mapRowsToSlice[T](rows)
		rows.Close()

		if err != nil {
			return nil, fmt.Errorf("failed to map rows: %v", err)
		}

		if len(batchResult) == 0 {
			// Không còn rows nào nữa
			break
		}

		// Append batch result vào slice kết quả cuối cùng
		results = append(results, batchResult...)
	}

	// Đóng cursor
	_, err = tx.Exec(context.Background(), fmt.Sprintf("CLOSE %s", cursorName))
	if err != nil {
		return nil, fmt.Errorf("failed to close cursor: %v", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return results, nil
}

func ExecuteStoredProcedureWithTable(conn *pgx.Conn, procedureName string, paramsStruct interface{}, resultStruct interface{}) error {
	args, cursorName, err := extractParams(paramsStruct)
	if err != nil {
		return fmt.Errorf("failed to extract parameters: %v", err)
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Declare a cursor
	cursorName = fmt.Sprintf("%s_cursor_%d", procedureName, time.Now().UnixNano())
	_, err = tx.Exec(context.Background(), fmt.Sprintf("DECLARE %s CURSOR FOR SELECT * FROM %s(%s)", cursorName, procedureName, buildPlaceholders(len(args))), args...)
	if err != nil {
		return fmt.Errorf("failed to declare cursor: %v", err)
	}

	resultSlice := reflect.ValueOf(resultStruct).Elem()

	for {
		// Fetch a batch of rows
		rows, err := tx.Query(context.Background(), fmt.Sprintf("FETCH %d FROM %s", 1000, cursorName))
		if err != nil {
			return fmt.Errorf("failed to fetch from cursor: %v", err)
		}

		batchResult := reflect.New(resultSlice.Type()).Elem()
		err = mapRowsToStruct(rows, batchResult.Addr().Interface())
		rows.Close()

		if err != nil {
			return fmt.Errorf("failed to map rows: %v", err)
		}

		if batchResult.Len() == 0 {
			// No more rows
			break
		}

		// Append batch result to the final result slice
		resultSlice.Set(reflect.AppendSlice(resultSlice, batchResult))
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func normalizeFieldName(name string) string {
	return strings.ReplaceAll(strings.ToUpper(name), "_", "")
}

func mapRowsToSlice[T any](rows pgx.Rows) ([]T, error) {
	var results []T

	fields := rows.FieldDescriptions()

	for rows.Next() {
		var result T
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		resultValue := reflect.ValueOf(&result).Elem()
		resultType := resultValue.Type()

		for i, field := range fields {
			dbFieldName := normalizeFieldName(field.Name)

			for j := 0; j < resultType.NumField(); j++ {
				structField := resultType.Field(j)
				structFieldValue := resultValue.Field(j)

				pgColumnTag := structField.Tag.Get("pgColumn")
				var fieldToCompare string

				if pgColumnTag != "" {
					fieldToCompare = normalizeFieldName(pgColumnTag)
				} else {
					fieldToCompare = normalizeFieldName(structField.Name)
				}

				if fieldToCompare == dbFieldName {
					if structFieldValue.IsValid() && structFieldValue.CanSet() {
						if values[i] != nil {
							// Xử lý các kiểu dữ liệu đặc biệt nếu cần
							switch structFieldValue.Kind() {
							case reflect.Struct:
								// Xử lý các kiểu như time.Time
								if timeValue, ok := values[i].(time.Time); ok {
									structFieldValue.Set(reflect.ValueOf(timeValue))
								}
							default:
								// Cho các kiểu dữ liệu cơ bản
								structFieldValue.Set(reflect.ValueOf(values[i]))
							}
						}
					}
					break
				}
			}
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Helper function to map rows from cursor to result struct
func mapRowsToStruct(rows pgx.Rows, resultStruct interface{}) error {
	resultSlice := reflect.ValueOf(resultStruct).Elem()
	elemType := resultSlice.Type().Elem()

	columnNames := rows.FieldDescriptions()
	columnMap := make(map[string]int)
	for i, col := range columnNames {
		columnMap[strings.ToLower(string(col.Name))] = i
	}

	for rows.Next() {
		resultElement := reflect.New(elemType).Elem()
		values := make([]interface{}, len(columnNames))
		fieldPointers := make([]interface{}, len(columnNames))

		for i := 0; i < elemType.NumField(); i++ {
			fieldType := elemType.Field(i)
			field := resultElement.Field(i)

			columnName := fieldType.Tag.Get("pgColumn")
			if columnName == "" {
				columnName = strings.ToLower(fieldType.Name)
			}

			if colIdx, ok := columnMap[columnName]; ok {
				fieldPointers[colIdx] = field.Addr().Interface()
				values[colIdx] = fieldPointers[colIdx]
			}
		}

		err := rows.Scan(values...)
		if err != nil {
			return fmt.Errorf("failed to scan row data: %v", err)
		}

		resultSlice.Set(reflect.Append(resultSlice, resultElement))
	}

	return nil
}

// Helper function to extract parameters from a struct using `pgParam` tags
func extractParams(paramsStruct interface{}) ([]interface{}, string, error) {
	val := reflect.ValueOf(paramsStruct)
	if val.Kind() != reflect.Struct {
		return nil, "", fmt.Errorf("paramsStruct must be a struct")
	}

	var args []interface{}
	var cursorName string
	typ := reflect.TypeOf(paramsStruct)

	// Iterate over struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Check for pgCur tag
		if pgCur := fieldType.Tag.Get("pgCur"); pgCur != "" {
			cursorName = pgCur
			continue // Skip adding this to args
		}

		// Get the pgParam tag value
		pgParam := fieldType.Tag.Get("pgParam")
		if pgParam == "" {
			return nil, "", fmt.Errorf("field %s does not have pgParam tag", fieldType.Name)
		}

		// Append the field's value to args
		args = append(args, field.Interface())
	}

	return args, cursorName, nil
}

// Helper function to create placeholders for SQL query
func buildPlaceholders(count int) string {
	placeholders := ""
	for i := 1; i <= count; i++ {
		placeholders += fmt.Sprintf("$%d", i)
		if i < count {
			placeholders += ", "
		}
	}
	return placeholders
}
