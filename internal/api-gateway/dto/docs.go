package api_gateway_dto

type ResponseSuccessDocs[T any] struct {
	Data     T        `json:"data,omitempty"`
	Metadata Metadata `json:"metadata"`
}

type ResponseSuccessPagingationDocs[T any] struct {
	Data     T                      `json:"data,omitempty"`
	Metadata MetadataWithPagination `json:"metadata"`
}

type ResponseErrorDocs struct {
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

type DeleteAddressTypeResponseDocs = ResponseSuccessDocs[DeleteAddressTypeByAdminResponse]
type UpdateAddressTypeResponseDocs = ResponseSuccessDocs[UpdateAddressTypeByAdminResponse]
type CreateAddressTypeResponseDocs = ResponseSuccessDocs[CreateAddressTypeByAdminResponse]
type ListAddressTypesResponseDocs = ResponseSuccessDocs[[]GetAddressTypeQueryResponse]
type GetAddressTypeByIdResponseDocs = ResponseSuccessDocs[GetAddressTypeByIdResponse]
