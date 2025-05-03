package api_gateway_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"time"
)

type AdministrativeDivision struct {
	ID        string     `json:"Id"`
	Name      string     `json:"Name"`
	Districts []District `json:"Districts,omitempty"`
}

type District struct {
	ID    string `json:"Id"`
	Name  string `json:"Name"`
	Wards []Ward `json:"Wards,omitempty"`
}

type Ward struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

type administrativeDivisionRepository struct {
	cache  pkg.ICache
	tracer pkg.Tracer
}

func NewAdministrativeDivisionRepository(cache pkg.ICache, tracer pkg.Tracer) IAdministrativeDivisionRepository {
	return &administrativeDivisionRepository{
		cache:  cache,
		tracer: tracer,
	}
}

func (r *administrativeDivisionRepository) GetProvinces(ctx context.Context) ([]AdministrativeDivision, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetProvinces"))
	defer span.End()

	// Check if data exists in cache
	provincesData, err := r.cache.Get(ctx, "provinces")
	if err == nil && provincesData != "" {
		// Data exists in cache, unmarshal it
		var provinces []AdministrativeDivision
		if err := json.Unmarshal([]byte(provincesData), &provinces); err != nil {
			span.RecordError(err)
			return nil, err
		}
		return provinces, nil
	}

	// Data not in cache, need to return empty result
	// We assume data is loaded into Redis during service initialization
	return []AdministrativeDivision{}, nil
}

func (r *administrativeDivisionRepository) GetDistrictsByProvinceID(ctx context.Context, provinceID string) ([]District, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetDistrictsByProvinceID"))
	defer span.End()

	// Check if data exists in cache
	key := fmt.Sprintf("districts:%s", provinceID)
	districtsData, err := r.cache.Get(ctx, key)
	if err == nil && districtsData != "" {
		// Data exists in cache, unmarshal it
		var districts []District
		if err := json.Unmarshal([]byte(districtsData), &districts); err != nil {
			span.RecordError(err)
			return nil, err
		}
		return districts, nil
	}

	// If not in specific key, try to get from provinces and extract
	provincesData, err := r.cache.Get(ctx, "provinces")
	if err == nil && provincesData != "" {
		var provinces []AdministrativeDivision
		if err := json.Unmarshal([]byte(provincesData), &provinces); err != nil {
			span.RecordError(err)
			return nil, err
		}

		// Find the province
		for _, province := range provinces {
			if province.ID == provinceID {
				// Cache districts for future use
				districtsJSON, _ := json.Marshal(province.Districts)
				r.cache.Set(ctx, key, string(districtsJSON), 24*time.Hour)
				return province.Districts, nil
			}
		}
	}

	return []District{}, nil
}

func (r *administrativeDivisionRepository) GetWardsByDistrictID(ctx context.Context, provinceID, districtID string) ([]Ward, error) {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetWardsByDistrictID"))
	defer span.End()

	// Check if data exists in cache
	key := fmt.Sprintf("wards:%s:%s", provinceID, districtID)
	wardsData, err := r.cache.Get(ctx, key)
	if err == nil && wardsData != "" {
		// Data exists in cache, unmarshal it
		var wards []Ward
		if err := json.Unmarshal([]byte(wardsData), &wards); err != nil {
			span.RecordError(err)
			return nil, err
		}
		return wards, nil
	}

	// If not in specific key, get districts and extract
	districts, err := r.GetDistrictsByProvinceID(ctx, provinceID)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Find the district
	for _, district := range districts {
		if district.ID == districtID {
			// Cache wards for future use
			wardsJSON, _ := json.Marshal(district.Wards)
			r.cache.Set(ctx, key, string(wardsJSON), 24*time.Hour)
			return district.Wards, nil
		}
	}

	return []Ward{}, nil
}
