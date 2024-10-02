package analysis_model

import "time"

// Các struct model (SoilAnalysisRequest, Requests, Gardens, Soils)
type SoilAnalysisRequest struct {
	Requests Requests `json:"requests"`
	Gardens  Gardens  `json:"gardens"`
	Soils    []Soils  `json:"soils"`
}

type Requests struct {
	Code         string    `json:"code"`
	RequestDate  string    `json:"request_date"`
	Requester    string    `json:"requester"`
	Approver     string    `json:"approver"`
	ApprovalDate time.Time `json:"approval_date"`
	Status       string    `json:"status"`
	CreatedBy    string    `json:"created_by"`
	UpdatedBy    string    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Gardens struct {
	RequestID          string    `json:"request_id"`
	CompanyCode        string    `json:"company_code"`
	FarmCode           string    `json:"farm_code"`
	PlotCode           string    `json:"plot_code"`
	PlotArea           float64   `json:"plot_area"`
	SoilTypeCode       string    `json:"soil_type_code"`
	AgeTree            float64   `json:"age_tree"`
	GrowthStatus       bool      `json:"growth_status"`
	StandardCode       string    `json:"standard_code"`
	GroundCover        bool      `json:"ground_cover"`
	LeafColorCode      string    `json:"leaf_color_code"`
	VarietyCode        string    `json:"variety_code"`
	CreatedBy          string    `json:"created_by"`
	UpdatedBy          string    `json:"updated_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	StandardPercentage int       `json:"standard_percentage"`
}

type Soils struct {
	RequestID      string    `json:"request_id"`
	ParameterCode  string    `json:"parameter_code"`
	ParameterValue float64   `json:"parameter_value"`
	CreatedBy      string    `json:"created_by"`
	UpdatedBy      string    `json:"updated_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
