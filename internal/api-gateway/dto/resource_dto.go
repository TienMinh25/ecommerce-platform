package api_gateway_dto

type CreateResource struct {
	Name string `json:"name" binding:"required,min=3,max=50,alphanumunicode"`
}

type CreateResourceResponse struct{}

type UpdateResource struct {
	Name string `json:"name" binding:"required,min=3,max=50,alphanumunicode"`
	ID   int    `json:"id" binding:"required,gt=0"`
}

type UpdateResourceResponse struct{}
