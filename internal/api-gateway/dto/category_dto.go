package api_gateway_dto

type GetCategoriesRequest struct {
	ParentID *int64 `form:"parent_id" binding:"omitempty,gte=1"`
}

type GetCategoriesResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	ParentID *int   `json:"parent_id"`
}
