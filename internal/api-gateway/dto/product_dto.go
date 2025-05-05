package api_gateway_dto

type GetProductsRequest struct {
	Limit       int64   `form:"limit,default=20" binding:"omitempty,gte=1"`
	Page        int64   `form:"page,default=1" binding:"omitempty,gte=1"`
	Keyword     *string `form:"keyword" binding:"omitempty"`
	CategoryIDs []int64 `form:"category_ids" binding:"omitempty"`
	MinRating   *int64  `form:"min_rating" binding:"omitempty,gte=1"`
}

type GetProductsResponse struct {
	ProductID            string  `json:"product_id"`
	ProductName          string  `json:"product_name"`
	ProductThumbnail     string  `json:"product_thumbnail"`
	ProductAverageRating float32 `json:"product_average_rating"`
	ProductTotalReviews  int64   `json:"product_total_reviews"`
	ProductCategoryID    int64   `json:"product_category_id"`
	ProductPrice         float64 `json:"product_price"`
	ProductDiscountPrice float64 `json:"product_discount_price"`
	ProductCurrency      string  `json:"product_currency"`
}
