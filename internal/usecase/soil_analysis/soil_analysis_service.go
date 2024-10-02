package soil_analysis

import (
	"app-server/internal/shared/analysis_model"
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type SoilAnalysisServiceInterface interface {
	AddNewAnalysis(soilAnalysisRequest analysis_model.SoilAnalysisRequest) (bool, error)
}

type SoilAnalysisService struct {
	db *gorm.DB
}

func NewSoilAnalysisServiceInterface(db *gorm.DB) SoilAnalysisServiceInterface {
	return &SoilAnalysisService{db: db}
}

func (s *SoilAnalysisService) AddNewAnalysis(soilAnalysisRequest analysis_model.SoilAnalysisRequest) (bool, error) {

	/*
		// Tạo dữ liệu mẫu
			soilAnalysisRequest := SoilAnalysisRequest{
				Requests: Requests{
					Code:         "REQ011",
					RequestDate:  "2023-05-15",
					Requester:    "Emma Watson",
					Approver:     "Tom Hardy",
					ApprovalDate: time.Now(),
					Status:       "Approved",
					CreatedBy:    "system",
					UpdatedBy:    "system",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
				Gardens: Gardens{
					CompanyCode:        "COMP011",
					FarmCode:           "FARM011",
					PlotCode:           "LO011",
					PlotArea:           155.75,
					SoilTypeCode:       "1",
					AgeTree:            4.30,
					GrowthStatus:       true,
					StandardCode:       "STD011",
					GroundCover:        true,
					LeafColorCode:      "XD",
					VarietyCode:        "VAR011",
					CreatedBy:          "system",
					UpdatedBy:          "system",
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
					StandardPercentage: 70,
				},
				Soils: []Soils{
					{
						ParameterCode:  "MN",
						ParameterValue: 0.31,
						CreatedBy:      "system",
						UpdatedBy:      "system",
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					},
					// Add other soil items here...
				},
			}
	*/

	// Convert struct to JSON
	jsonData, err := json.Marshal(soilAnalysisRequest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
	// Prepare variables to receive results
	var result *string
	var errorCode int
	var errorMessage string

	// Call stored procedure
	err = s.db.Raw("CALL demovcs.insert_soil_analysis(?, ?, ?, ?)",
		string(jsonData), &result, &errorCode, &errorMessage).Error
	if err != nil {
		log.Fatal(err)
	}

	// Check results

	if errorCode != 0 {
		fmt.Printf("Operation failed. Error code: %d, Error message: %s\n", errorCode, errorMessage)
		//c.Error("Operation failed. Error code: %d, Error message: %s\n", errorCode, errorMessage)
		return false, fmt.Errorf("operation failed. Error code: %d, Error message: %s", errorCode, errorMessage)
	}
	return true, nil
}
