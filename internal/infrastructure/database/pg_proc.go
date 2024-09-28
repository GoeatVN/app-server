package database

import (
	"app-server/internal/infrastructure/config"
	"context"
	"database/sql"
	"fmt"

	"log"

	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
Demo call procedure with cursor
// Example of a struct result (each procedure can have different results)
	type Result struct {
		CntName        string    `pgColumn:"cnt_name" json:"cntName"`
		CntInt         int       `pgColumn:"cnt_int" json:"cntInt"`
		CntTest        string    `pgColumn:"cnt_test" json:"cntTest"`
		CntDate        time.Time `pgColumn:"cnt_date" json:"cntDate"`
		CntTimeStamp   time.Time `pgColumn:"cnt_timestamp"`
		CntTime        time.Time `pgColumn:"cnt_time"`
		CtnNumeric     float64   `pgColumn:"cnt_numeric"`
		CtnBoolean     bool      `pgColumn:"cnt_boolean"`
		CtnTimestampTz time.Time `pgColumn:"cnt_timestamptz"`
		ctnTemp        float64   `pgColumn:"cnt_temp"`
		CntId          int       `pgColumn:"cnt_id" json:"cntId"`
	}
	// Struct containing parameters with `pgParam` tags

	type ContentParams struct {
		CntId   int       `pgParam:"prm_id"`
		CntName string    `pgParam:"prm_name"`
		CntDate time.Time `pgParam:"prm_date"`
		CurRef  string    `pgCur:"cur_ref"`
	}
	// Định nghĩa struct chứa tham số đầu vào cho stored procedure
	params := ContentParams{
		CntId:   1,
		CntName: "thong",
		CntDate: time.Now(),
	}

	// Khởi tạo slice để chứa kết quả trả về
	var results []Result
	results, err := persistence2.ExecuteStoredProcedureWithCursor[Result]("content_get_test", params)
	if err != nil {
		fmt.Println("Lỗi khi gọi stored procedure: %v\n", err)
	}

	// In kết quả
	for _, result := range results {
		fmt.Printf("CntId: %d, CntTimeStamp: %s\n", result.CntId, result.CntTimeStamp)

	}
*/

func ExecProcCursor[T any](procedureName string, paramsStruct interface{}) ([]T, error) {
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
	connInfo, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		connInfo.Database.Host,
		connInfo.Database.User,
		connInfo.Database.Password,
		connInfo.Database.Name,
		connInfo.Database.Port,
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Cấu hình pool (tùy chỉnh theo nhu cầu)
	config.MaxConns = int32(connInfo.Database.MaxConns)
	config.MinConns = int32(connInfo.Database.MinConns)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pool, nil
}

func normalizeFieldName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), "_", "")
}

func mapRowsToSlice[T any](rows pgx.Rows) ([]T, error) {
	var results []T
	resultType := reflect.TypeOf((*T)(nil)).Elem()

	// Lấy danh sách các tên cột từ kết quả truy vấn
	columnNames := rows.FieldDescriptions()
	columnMap := make(map[string]int)
	for i, col := range columnNames {
		// Sử dụng cột dưới dạng chữ thường để tránh lỗi không khớp tên
		columnMap[normalizeFieldName(string(col.Name))] = i
	}

	// Duyệt qua từng hàng dữ liệu
	for rows.Next() {
		// Tạo phần tử mới của kiểu kết quả (struct)
		resultElement := reflect.New(resultType).Elem()

		// Tạo slice chứa các giá trị tương ứng với số lượng cột
		values := make([]interface{}, len(columnNames))

		// Map các field của struct với các cột trong kết quả
		for i := 0; i < resultType.NumField(); i++ {
			fieldType := resultType.Field(i)
			field := resultElement.Field(i)

			// Lấy tên cột từ thẻ pgColumn, nếu không có thì dùng tên field trong struct
			columnName := normalizeFieldName(fieldType.Tag.Get("pgColumn"))
			if columnName == "" {
				columnName = normalizeFieldName(fieldType.Name)
			}

			// Tìm chỉ số của cột trong kết quả trả về
			if colIdx, ok := columnMap[columnName]; ok {
				// Sử dụng các kiểu sql.NullXXX để xử lý null
				switch field.Kind() {
				case reflect.Int, reflect.Int32, reflect.Int64:
					values[colIdx] = new(sql.NullInt64)
				case reflect.String:
					values[colIdx] = new(sql.NullString)
				case reflect.Float32, reflect.Float64:
					values[colIdx] = new(sql.NullFloat64)
				case reflect.Struct:
					// Kiểm tra nếu kiểu là time.Time thì dùng sql.NullTime
					if fieldType.Type == reflect.TypeOf(time.Time{}) {
						values[colIdx] = new(sql.NullTime)
					}
				default:
					// Nếu không phải kiểu null thì map trực tiếp
					values[colIdx] = field.Addr().Interface()
				}
			}
		}

		// Thực hiện quét dữ liệu từ hàng vào slice values (dựa theo chỉ số cột)
		err := rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row data: %v", err)
		}

		// Map giá trị từ slice values vào các field của struct, kiểm tra null
		for i := 0; i < resultType.NumField(); i++ {
			field := resultElement.Field(i)
			columnName := normalizeFieldName(resultType.Field(i).Tag.Get("pgColumn"))
			if columnName == "" {
				columnName = normalizeFieldName(resultType.Field(i).Name)
			}

			// Lấy chỉ số cột từ columnMap
			if colIdx, ok := columnMap[columnName]; ok {
				// Kiểm tra giá trị trả về và map dữ liệu nếu không null
				switch v := values[colIdx].(type) {
				case *sql.NullInt64:
					if v.Valid {
						field.SetInt(v.Int64)
					}
				case *sql.NullString:
					if v.Valid {
						field.SetString(v.String)
					}
				case *sql.NullFloat64:
					if v.Valid {
						field.SetFloat(v.Float64)
					}
				case *sql.NullTime:
					if v.Valid {
						// Chuyển đổi sql.NullTime thành time.Time
						field.Set(reflect.ValueOf(v.Time))
					}
				}
			}
		}

		// Thêm phần tử đã map vào kết quả
		results = append(results, resultElement.Interface().(T))
	}

	// Kiểm tra lỗi sau khi duyệt các hàng
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
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

// Generic function to execute stored procedures with pgx v5
func ExecuteStoredProcedure_Not_Used(conn *pgx.Conn, procedureName string, paramsStruct interface{}, resultStruct interface{}) error {
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
	err = mapRowsToStruct_Not_Used(rows, resultStruct)
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
func mapRowsToSlice_Not_Used[T any](rows pgx.Rows) ([]T, error) {
	var results []T
	resultType := reflect.TypeOf((*T)(nil)).Elem()

	columnNames := rows.FieldDescriptions()
	columnMap := make(map[string]int)
	for i, col := range columnNames {
		columnMap[normalizeFieldName(string(col.Name))] = i
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
				columnName = normalizeFieldName(fieldType.Name)
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
func mapRowsToStruct_Not_Used(rows pgx.Rows, resultStruct interface{}) error {
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
