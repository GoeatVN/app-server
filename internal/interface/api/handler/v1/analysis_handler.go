package v1

import (
	"app-server/internal/shared/analysis_model"

	"app-server/internal/usecase/soil_analysis"
	"app-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type SoilAnalysisHandler struct {
	soilAnalysisService soil_analysis.SoilAnalysisServiceInterface
}

func NewSoilAnalysisHandler(soilAnalysisService soil_analysis.SoilAnalysisServiceInterface) *SoilAnalysisHandler {
	return &SoilAnalysisHandler{soilAnalysisService: soilAnalysisService}
}

// Login handles user login
func (h *SoilAnalysisHandler) SaveSoilAnalysis(c *gin.Context) {
	/*
		{"request":{"id":""},"garden":{"id":"","company_code":"CT001","farm_code":"FRO001","plot_code":"LO001","plot_area":1000,"soil_type_code":"2","age_tree":10,"planting_year":2008,"root_depth":30,"growth_status":true,"standard_code":true,"ground_cover":true,"leaf_color_code":"XD","variety_code":"CS001","standard_percentage":70},"soils":[{"id":"","parameter_code":"PH","parameter_value":1.2},{"id":"","parameter_code":"NH3","parameter_value":1},{"id":"","parameter_code":"NO3","parameter_value":0.7},{"id":"","parameter_code":"P2O5","parameter_value":100},{"id":"","parameter_code":"K2O","parameter_value":60},{"id":"","parameter_code":"MG","parameter_value":65},{"id":"","parameter_code":"CU","parameter_value":70},{"id":"","parameter_code":"ZN","parameter_value":1.2}],"user_process":"1"}
	*/

	var request analysis_model.SoilAnalysisRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err)
		return
	}

	soilAnalysisResponse, err := h.soilAnalysisService.SaveSoilAnalysis(request)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("response_data", soilAnalysisResponse)
}
