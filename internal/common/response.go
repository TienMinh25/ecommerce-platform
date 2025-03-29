package common

type ResponseSuccess[T any] struct {
	Data     T        `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type ResponseSuccessPagingation[T any] struct {
	Data     T                      `json:"data"`
	Metadata MetadataWithPagination `json:"metadata"`
}

type ResponseError struct {
	Metadata Metadata    `json:"metadata"`
	Error    interface{} `json:"error,omitempty"`
}

type Metadata struct {
	Code int `json:"code"`
}

type MetadataWithPagination struct {
	Code       int         `json:"code"`
	Pagination *Pagination `json:"pagination"`
}

type Pagination struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}
