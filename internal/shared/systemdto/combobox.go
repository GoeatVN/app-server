package systemdto

type ComboType string

const (
	AllRole ComboType = "AllRole"
)

type ComboboxRequestItem struct {
	ComboType ComboType `json:"comboType"`
	Param     string    `json:"param"`
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	Param3    string    `json:"param3"`
}

type ComboboxRequest struct {
	Data []ComboboxRequestItem `json:"data"`
}

type ComboboxDto struct {
	ID     int    `json:"id"`
	Name   string `json:"text"`
	Value  string `json:"value"`
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
	Value3 string `json:"value3"`
}

type ComboboxResponseItem struct {
	Data      []ComboboxDto `json:"data"`
	ComboType ComboType     `json:"comboType"`
}

type ComboboxResponse struct {
	Data []ComboboxResponseItem `json:"data"`
}
