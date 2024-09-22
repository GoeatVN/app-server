package persistence

import (
	"context"
	"fmt"
	"log"

	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func ExecuteStoredProcedureWithCursor[T any](procedureName string, paramsStruct interface{}) ([]T, error) {
	ctx := context.Background()
	// Sử dụng connection pool thay vì kết nối đơn
	pool, err := connectToDBPool(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	args, cursorName, err := extractParams(paramsStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to extract parameters: %w", err)
	}

	if cursorName == "" {
		cursorName = "ref_cur" // Default cursor name if not provided
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Thêm tên cursor vào danh sách tham số
	args = append(args, cursorName)

	// Gọi function với cursor
	_, err = tx.Exec(ctx, fmt.Sprintf("SELECT %s(%s)", procedureName, buildPlaceholders(len(args))), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute function: %w", err)
	}

	var results []T
	batchSize := 1000
	for {
		// Fetch một batch các rows
		rows, err := tx.Query(ctx, fmt.Sprintf("FETCH %d FROM %s", batchSize, cursorName))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch from cursor: %w", err)
		}

		batchResult, err := mapRowsToSlice[T](rows)
		//pgx.CollectRows(rows, pgx.RowToStructByName[T])
		if err != nil {
			return nil, fmt.Errorf("failed to collect rows: %w", err)
		}

		if len(batchResult) == 0 {
			// Không còn rows nào nữa
			break
		}

		// Append batch result vào slice kết quả cuối cùng
		results = append(results, batchResult...)
	}

	// Đóng cursor
	_, err = tx.Exec(ctx, fmt.Sprintf("CLOSE %s", cursorName))
	if err != nil {
		return nil, fmt.Errorf("failed to close cursor: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return results, nil
}

func connectToDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"))

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Cấu hình pool (tùy chỉnh theo nhu cầu)
	config.MaxConns = 10
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pool, nil
}
func ExecuteStoredProcedureWithTable(procedureName string, paramsStruct interface{}, resultStruct interface{}) error {
	// Kết nối db postgres
	conn, err := connectToDB()
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
func connectToDB() (conn *pgx.Conn, err error) {
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	connInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		host,
		port,
		user,
		dbname,
		password)
	// Kết nối đến database qua biến môi trường DATABASE_URL
	conn, err = pgx.Connect(context.Background(), connInfo)
	if err != nil {
		log.Fatalf("Không thể kết nối đến database: %v\n", err)
	}

	return conn, err
}

func normalizeFieldName(name string) string {
	return strings.ReplaceAll(strings.ToUpper(name), "_", "")
}

func mapRowsToSlice[T any](rows pgx.Rows) ([]T, error) {
	var results []T
	resultType := reflect.TypeOf((*T)(nil)).Elem()

	columnNames := rows.FieldDescriptions()
	columnMap := make(map[string]int)
	for i, col := range columnNames {
		columnMap[strings.ToLower(string(col.Name))] = i
	}

	for rows.Next() {
		resultElement := reflect.New(resultType).Elem()
		values := make([]interface{}, len(columnNames))
		fieldPointers := make([]interface{}, len(columnNames))

		for i := 0; i < resultType.NumField(); i++ {
			fieldType := resultType.Field(i)
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
			return nil, fmt.Errorf("failed to scan row data: %v", err)
		}

		results = append(results, resultElement.Interface().(T))
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
