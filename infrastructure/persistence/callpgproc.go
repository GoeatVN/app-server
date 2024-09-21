package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
)

// Generic function to execute stored procedures with pgx v5
// paramsStruct is a struct with fields tagged by `pgParam:"param_name"`
func ExecuteStoredProcedure(conn *pgx.Conn, procedureName string, paramsStruct interface{}, resultStruct interface{}) error {
	// Extract parameters from struct using reflection and pgParam tag
	args, err := extractParams(paramsStruct)
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
	cursorName := "ref"
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

// Helper function to map rows from cursor to result struct
func mapRowsToStruct(rows pgx.Rows, resultStruct interface{}) error {
	// Get the reflection value of the result slice
	resultSlice := reflect.ValueOf(resultStruct).Elem()

	// Get the element type of the slice (i.e., the type of the struct)
	elemType := resultSlice.Type().Elem()

	// Get the column names from the rows
	columnNames := rows.FieldDescriptions()
	columnMap := make(map[string]int)
	for i, col := range columnNames {
		columnMap[string(col.Name)] = i
	}

	// Iterate over rows and map the results dynamically to the result struct
	for rows.Next() {
		// Create a new element of the result type
		resultElement := reflect.New(elemType).Elem()

		// Create a slice to hold scanned values (can use sql.NullXXX for nullable columns)
		values := make([]interface{}, elemType.NumField())

		// Map fields based on pgColumn tags or field names
		for i := 0; i < elemType.NumField(); i++ {
			fieldType := elemType.Field(i)
			field := resultElement.Field(i)

			// Get pgColumn tag or fallback to field name
			columnName := fieldType.Tag.Get("pgColumn")
			if columnName == "" {
				columnName = strings.ToLower(fieldType.Name) // Default to lowercase field name if no tag
			}

			// Find the column index in the result set
			if _, ok := columnMap[columnName]; ok {
				// Create a pointer for each field to scan the value into
				values[i] = field.Addr().Interface()
			} else {
				// If the column doesn't exist, set the value to null (sql.NullXXX)
				values[i] = new(sql.NullString) // Adjust based on field type
			}
		}

		// Scan row data into the struct fields
		err := rows.Scan(values...)
		if err != nil {
			return fmt.Errorf("failed to scan row data: %v", err)
		}

		// Append the result element to the result slice
		resultSlice.Set(reflect.Append(resultSlice, resultElement))
	}

	return nil
}

// Helper function to extract parameters from a struct using `pgParam` tags
func extractParams(paramsStruct interface{}) ([]interface{}, error) {
	val := reflect.ValueOf(paramsStruct)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("paramsStruct must be a struct")
	}

	var args []interface{}
	typ := reflect.TypeOf(paramsStruct)

	// Iterate over struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the pgParam tag value
		pgParam := fieldType.Tag.Get("pgParam")
		if pgParam == "" {
			return nil, fmt.Errorf("field %s does not have pgParam tag", fieldType.Name)
		}

		// Append the field's value to args
		args = append(args, field.Interface())
	}

	return args, nil
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
