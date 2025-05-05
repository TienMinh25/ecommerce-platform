package api_gateway_dto

type GetCategoriesRequest struct {
	ParentID       *int64  `form:"parent_id" binding:"omitempty,gte=1"`
	ProductKeyword *string `form:"product_keyword"`
}

type GetCategoriesResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ImageURL     string `json:"image_url"`
	ParentID     *int   `json:"parent_id"`
	ProductCount *int64 `json:"product_count"`
	IsSelected   *bool  `json:"is_selected"`
}
