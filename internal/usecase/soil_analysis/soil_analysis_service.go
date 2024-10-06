package soil_analysis

import (
	"app-server/internal/shared/analysis_model"
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type SoilAnalysisServiceInterface interface {
	SaveSoilAnalysis(soilAnalysisRequest analysis_model.SoilAnalysisRequest) (bool, error)
}

type SoilAnalysisService struct {
	db *gorm.DB
}

func NewSoilAnalysisServiceInterface(db *gorm.DB) SoilAnalysisServiceInterface {
	return &SoilAnalysisService{db: db}
}

func (s *SoilAnalysisService) SaveSoilAnalysis(soilAnalysisRequest analysis_model.SoilAnalysisRequest) (bool, error) {

	// Convert struct to JSON
	jsonData, err := json.Marshal(soilAnalysisRequest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

	// Cấu trúc để lưu kết quả của các INOUT parameters
	type ProcedureResult struct {
		p_result        string
		p_error_code    string
		p_error_message string
	}
	var result ProcedureResult

	// Chuẩn bị câu lệnh SQL để gọi stored procedure
	query := `CALL demovcs.save_soil_analysis_data(?, ?, ?, ?)`

	// Thực thi câu lệnh SQL với tham số JSON và nhận kết quả từ INOUT parameters
	err = s.db.Raw(query, string(jsonData), result.p_result, result.p_error_code, result.p_error_message).Scan(&result).Error
	if err != nil {
		return false, err
	}

	fmt.Printf("Result: %s\n,", result.p_result)
	// Check results

	if result.p_error_code != "" {
		fmt.Printf("Operation failed. Error code: %s, Error message: %s\n", result.p_error_code, result.p_error_message)

		return false, fmt.Errorf("operation failed. Error code: %s, Error message: %s", result.p_error_code, result.p_error_message)
	}
	return true, nil
}
