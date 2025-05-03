package api_gateway_dto

// ProvinceResponse represents a province/city
type ProvinceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DistrictResponse represents a district
type DistrictResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// WardResponse represents a ward
type WardResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
