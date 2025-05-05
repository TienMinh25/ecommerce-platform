package api_gateway_service

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/supplier-and-product/grpc/proto/partner_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"net/http"
)

type categoryService struct {
	tracer        pkg.Tracer
	partnerClient partner_proto_gen.PartnerServiceClient
}

func NewCategoryService(tracer pkg.Tracer, partnerClient partner_proto_gen.PartnerServiceClient) ICategoryService {
	return &categoryService{
		tracer:        tracer,
		partnerClient: partnerClient,
	}
}

func (c *categoryService) GetCategories(ctx context.Context, query api_gateway_dto.GetCategoriesRequest) ([]api_gateway_dto.GetCategoriesResponse, error) {
	ctx, span := c.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetCategories"))
	defer span.End()

	// call grpc to get result
	response, err := c.partnerClient.GetCategories(ctx, &partner_proto_gen.GetCategoriesRequest{
		ParentId:       query.ParentID,
		ProductKeyword: query.ProductKeyword,
	})

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	res := make([]api_gateway_dto.GetCategoriesResponse, len(response.Categories))

	for idx, categoryResponse := range response.Categories {
		res[idx] = api_gateway_dto.GetCategoriesResponse{
			ID:           int(categoryResponse.CategoryId),
			Name:         categoryResponse.Name,
			ImageURL:     categoryResponse.ImageUrl,
			ProductCount: categoryResponse.ProductCount,
			IsSelected:   categoryResponse.Selected,
		}
	}

	return res, nil
}
