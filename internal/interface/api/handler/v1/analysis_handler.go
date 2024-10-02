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
func (h *SoilAnalysisHandler) AddNewSoilAnalysis(c *gin.Context) {

	var request analysis_model.SoilAnalysisRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, err)
		return
	}

	soilAnalysisResponse, err := h.soilAnalysisService.AddNewAnalysis(request)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("response_data", soilAnalysisResponse)
}
