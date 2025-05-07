package api_gateway_dto

import "time"

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
	CategoryName         string  `json:"category_name"`
}

type GetProductDetailRequest struct {
	ProductID string `uri:"productID" binding:"required,uuid"`
}

type GetProductDetailResponse struct {
	ProductID            string                            `json:"product_id"`
	ProductName          string                            `json:"name"`
	ProductDescription   string                            `json:"description"`
	CategoryID           int64                             `json:"category_id"`
	CategoryName         string                            `json:"category_name"`
	ProductAverageRating float32                           `json:"average_rating"`
	ProductTotalReviews  int64                             `json:"total_reviews"`
	Supplier             GetSupplierProductResponse        `json:"supplier"`
	ProductTags          []string                          `json:"product_tags"`
	Attributes           []ProductAttribute                `json:"attributes"`
	ProductVariants      []GetProductDetailVariantResponse `json:"product_variants"`
}

type ProductAttribute struct {
	AttributeID int64                `json:"attribute_id"`
	Name        string               `json:"name"`
	Values      AttributeOptionValue `json:"values"`
}

type AttributeOptionValue struct {
	OptionID int64  `json:"option_id"`
	Value    string `json:"value"`
}

type GetProductDetailVariantResponse struct {
	ProductVariantID string                 `json:"product_variant_id"`
	SKU              string                 `json:"sku"`
	VariantName      string                 `json:"variant_name"`
	Price            float64                `json:"price"`
	DiscountPrice    float64                `json:"discount_price"`
	Quantity         int64                  `json:"quantity"`
	IsDefault        bool                   `json:"is_default"`
	ShippingClass    string                 `json:"shipping_class"`
	ThumbnailURL     string                 `json:"thumbnail_url"`
	AltTextThumbnail string                 `json:"alt_text_thumbnail"`
	Currency         string                 `json:"currency"`
	AttributeValues  []VariantAttributePair `json:"attribute_values"`
}

type VariantAttributePair struct {
	AttributeName  string `json:"attribute_name"`
	AttributeValue string `json:"attribute_value"`
}

type GetSupplierProductResponse struct {
	SupplierID   string `json:"supplier_id"`
	CompanyName  string `json:"company_name"`
	Thumbnail    string `json:"thumbnail"`
	ContactPhone string `json:"contact_phone"`
}

type GetProductReviewsRequest struct {
	ProductID string `uri:"productID" binding:"required,uuid"`
	Limit     int64  `form:"limit,default=6" binding:"omitempty,gte=1"`
	Page      int64  `form:"page,default=1" binding:"omitempty,gte=1"`
}

type GetProductReviewsResponse struct {
	UserID        int64     `json:"user_id"`
	UserName      string    `json:"user_name"`
	UserAvatarURL string    `json:"user_avatar_url"`
	Rating        float32   `json:"rating"`
	Comment       string    `json:"comment"`
	HelpfulVotes  int64     `json:"helpful_votes"`
	CreatedAt     time.Time `json:"created_at"`
}
