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
	// Prepare variables to receive results
	var result *string
	var errorCode string
	var errorMessage string

	// Call stored procedure
	err = s.db.Raw("CALL save_soil_analysis_data(?, ?, ?, ?)",
		string(jsonData), &result, &errorCode, &errorMessage).Error
	if err != nil {
		log.Fatal(err)
	}

	// Check results

	if errorCode != "" {
		fmt.Printf("Operation failed. Error code: %s, Error message: %s\n", errorCode, errorMessage)

		return false, fmt.Errorf("operation failed. Error code: %s, Error message: %s", errorCode, errorMessage)
	}
	return true, nil
}
