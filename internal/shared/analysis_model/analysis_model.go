package analysis_model

// Các struct model (SoilAnalysisRequest, Requests, Gardens, Soils)
type SoilAnalysisRequest struct {
	Request     Request `json:"request"`
	Garden      Garden  `json:"garden"`
	Soils       []Soil  `json:"soils"`
	UserProcess string  `json:"user_process"`
}
type Request struct {
	ID string `json:"id"`
}
type Garden struct {
	ID                 string `json:"id"`
	CompanyCode        string `json:"company_code"`
	FarmCode           string `json:"farm_code"`
	PlotCode           string `json:"plot_code"`
	PlotArea           int    `json:"plot_area"`
	SoilTypeCode       string `json:"soil_type_code"`
	AgeTree            int    `json:"age_tree"`
	PlantingYear       int    `json:"planting_year"`
	RootDepth          int    `json:"root_depth"`
	GrowthStatus       bool   `json:"growth_status"`
	StandardCode       bool   `json:"standard_code"`
	GroundCover        bool   `json:"ground_cover"`
	LeafColorCode      string `json:"leaf_color_code"`
	VarietyCode        string `json:"variety_code"`
	StandardPercentage int    `json:"standard_percentage"`
}

type Soil struct {
	ID             string  `json:"id"`
	ParameterCode  string  `json:"parameter_code"`
	ParameterValue float64 `json:"parameter_value"`
}
