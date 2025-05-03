package api_gateway_service

import (
	"context"
	"embed"
	"encoding/json"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_repository "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/repository"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/redis/go-redis/v9"
	"time"
)

type administrativeDivisionService struct {
	repo   api_gateway_repository.IAdministrativeDivisionRepository
	cache  pkg.ICache
	tracer pkg.Tracer
}

func NewAdministrativeDivisionService(
	repo api_gateway_repository.IAdministrativeDivisionRepository,
	cache pkg.ICache,
	tracer pkg.Tracer,
) IAdministrativeDivisionService {
	return &administrativeDivisionService{
		repo:   repo,
		cache:  cache,
		tracer: tracer,
	}
}

//go:embed hanh-chinh-viet-nam.json
var administrativeData embed.FS

func (s *administrativeDivisionService) LoadDataToCache(ctx context.Context) error {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "LoadDataToCache"))
	defer span.End()

	// Read embedded file
	data, err := administrativeData.ReadFile("hanh-chinh-viet-nam.json")
	if err != nil {
		span.RecordError(err)
		return err
	}

	// Parse the JSON data
	var provinces []api_gateway_repository.AdministrativeDivision
	if err := json.Unmarshal(data, &provinces); err != nil {
		span.RecordError(err)
		return err
	}

	// Store in Redis cache
	provincesJSON, err := json.Marshal(provinces)
	if err != nil {
		span.RecordError(err)
		return err
	}

	// Set provinces in cache (24 hour TTL)
	if err := s.cache.Set(ctx, "provinces", string(provincesJSON), redis.KeepTTL); err != nil {
		span.RecordError(err)
		return err
	}

	// Pre-cache districts and wards for faster access
	for _, province := range provinces {
		districtsJSON, _ := json.Marshal(province.Districts)
		s.cache.Set(ctx, "districts:"+province.ID, string(districtsJSON), 24*time.Hour)

		for _, district := range province.Districts {
			wardsJSON, _ := json.Marshal(district.Wards)
			s.cache.Set(ctx, "wards:"+province.ID+":"+district.ID, string(wardsJSON), 24*time.Hour)
		}
	}

	return nil
}

func (s *administrativeDivisionService) GetProvinces(ctx context.Context) ([]api_gateway_dto.ProvinceResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetProvinces"))
	defer span.End()

	provinces, err := s.repo.GetProvinces(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Map to response api_gateway_dto
	response := make([]api_gateway_dto.ProvinceResponse, 0, len(provinces))
	for _, province := range provinces {
		response = append(response, api_gateway_dto.ProvinceResponse{
			ID:   province.ID,
			Name: province.Name,
		})
	}

	return response, nil
}

func (s *administrativeDivisionService) GetDistricts(ctx context.Context, provinceID string) ([]api_gateway_dto.DistrictResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetDistricts"))
	defer span.End()

	districts, err := s.repo.GetDistrictsByProvinceID(ctx, provinceID)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Map to response api_gateway_dto
	response := make([]api_gateway_dto.DistrictResponse, 0, len(districts))
	for _, district := range districts {
		response = append(response, api_gateway_dto.DistrictResponse{
			ID:   district.ID,
			Name: district.Name,
		})
	}

	return response, nil
}

func (s *administrativeDivisionService) GetWards(ctx context.Context, provinceID, districtID string) ([]api_gateway_dto.WardResponse, error) {
	ctx, span := s.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GetWards"))
	defer span.End()

	wards, err := s.repo.GetWardsByDistrictID(ctx, provinceID, districtID)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Map to response api_gateway_dto
	response := make([]api_gateway_dto.WardResponse, 0, len(wards))
	for _, ward := range wards {
		response = append(response, api_gateway_dto.WardResponse{
			ID:   ward.ID,
			Name: ward.Name,
		})
	}

	return response, nil
}
