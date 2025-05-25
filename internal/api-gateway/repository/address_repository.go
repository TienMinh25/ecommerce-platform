package api_gateway_repository

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"net/http"
)

type addressRepository struct {
	db     pkg.Database
	tracer pkg.Tracer
}

func NewAddressRepository(db pkg.Database, tracer pkg.Tracer) IAddressRepository {
	return &addressRepository{
		db:     db,
		tracer: tracer,
	}
}

func (a *addressRepository) GetCurrentAddressByUserID(ctx context.Context, limit, page, userID int) ([]api_gateway_models.Address, int, error) {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "GetCurrentAddressByUserID"))
	defer span.End()

	var totalItems int

	countQuery := "SELECT COUNT(*) FROM addresses WHERE user_id = $1"

	if err := a.db.QueryRow(ctx, countQuery, userID).Scan(&totalItems); err != nil {
		span.RecordError(err)

		return nil, 0, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	query := `SELECT a.id, a.recipient_name, a.phone, a.street, a.district,
			a.province, a.ward, a.postal_code, a.country, a.is_default, a.longtitude, a.latitude,
			a.address_type_id, at.address_type
			FROM addresses a
			INNER JOIN address_types at ON at.id = a.address_type_id
			WHERE a.user_id = $1
			ORDER BY a.is_default DESC, a.created_at ASC
			LIMIT $2 OFFSET $3`

	rows, err := a.db.Query(ctx, query, userID, limit, limit*(page-1))

	if err != nil {
		span.RecordError(err)
		return nil, 0, utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	addresses := make([]api_gateway_models.Address, 0)

	for rows.Next() {
		var address api_gateway_models.Address

		if err = rows.Scan(&address.ID, &address.RecipientName, &address.Phone, &address.Street,
			&address.District, &address.Province, &address.Ward, &address.PostalCode, &address.Country, &address.IsDefault,
			&address.Longtitude, &address.Latitude, &address.AddressTypeID, &address.AddressType); err != nil {
			span.RecordError(err)
			return nil, 0, utils.TechnicalError{
				Message: common.MSG_INTERNAL_ERROR,
				Code:    http.StatusInternalServerError,
			}
		}

		addresses = append(addresses, address)
	}

	return addresses, totalItems, nil
}

func (a *addressRepository) SetDefaultAddressByID(ctx context.Context, addressID, userID int) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "SetDefaultAddressByID"))
	defer span.End()

	return a.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}, func(tx pkg.Tx) error {
		queryUpdateOldDefault := `UPDATE addresses SET is_default = $1 WHERE user_id = $2`

		if err := tx.Exec(ctx, queryUpdateOldDefault, false, userID); err != nil {
			return utils.TechnicalError{
				Message: common.MSG_INTERNAL_ERROR,
				Code:    http.StatusInternalServerError,
			}
		}

		queryUpdate := `UPDATE addresses SET is_default = $1 WHERE user_id = $2 and id = $3`

		if err := tx.Exec(ctx, queryUpdate, true, userID, addressID); err != nil {
			return utils.TechnicalError{
				Message: common.MSG_INTERNAL_ERROR,
				Code:    http.StatusInternalServerError,
			}
		}

		return nil
	})
}

func (a *addressRepository) CreateNewAddress(ctx context.Context, data *api_gateway_dto.CreateAddressRequest, userID int) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "CreateNewAddress"))
	defer span.End()

	return a.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		if data.IsDefault {
			queryUpdate := `UPDATE addresses SET is_default = $1 WHERE user_id = $2`

			if err := tx.Exec(ctx, queryUpdate, false, userID); err != nil {
				return utils.TechnicalError{
					Message: common.MSG_INTERNAL_ERROR,
					Code:    http.StatusInternalServerError,
				}
			}
		}

		queryInsert := `INSERT INTO addresses (user_id, recipient_name, phone, street, district, province,
                       	ward, postal_code, country, is_default, longtitude, latitude, address_type_id)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

		var args []interface{}
		args = append(args, userID)
		args = append(args, data.RecipientName)
		args = append(args, data.Phone)
		args = append(args, data.Street)
		args = append(args, data.District)
		args = append(args, data.Province)
		args = append(args, data.Ward)
		args = append(args, data.PostalCode)
		args = append(args, data.Country)
		args = append(args, data.IsDefault)
		args = append(args, data.Longitude)
		args = append(args, data.Latitude)
		args = append(args, data.AddressTypeID)

		if err := tx.Exec(ctx, queryInsert, args...); err != nil {
			return utils.TechnicalError{
				Message: common.MSG_INTERNAL_ERROR,
				Code:    http.StatusInternalServerError,
			}
		}

		return nil
	})
}

func (a *addressRepository) UpdateAddressByID(ctx context.Context, data *api_gateway_dto.UpdateAddressRequest, userID, addressID int) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "UpdateAddressByID"))
	defer span.End()

	return a.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		if data.IsDefault {
			queryUpdate := `UPDATE addresses SET is_default = $1 WHERE user_id = $2`

			if err := tx.Exec(ctx, queryUpdate, false, userID); err != nil {
				return utils.TechnicalError{
					Message: common.MSG_INTERNAL_ERROR,
					Code:    http.StatusInternalServerError,
				}
			}
		}

		queryUpdate := `UPDATE addresses 
						SET recipient_name = $1, phone = $2, street = $3, district = $4,
						    province = $5, ward = $6, postal_code = $7, country = $8, is_default = $9,
						    longtitude = $10, latitude = $11, address_type_id = $12
						WHERE user_id = $13 AND id = $14`

		var args []interface{}
		args = append(args, data.RecipientName)
		args = append(args, data.Phone)
		args = append(args, data.Street)
		args = append(args, data.District)
		args = append(args, data.Province)
		args = append(args, data.Ward)
		args = append(args, data.PostalCode)
		args = append(args, data.Country)
		args = append(args, data.IsDefault)
		args = append(args, data.Longitude)
		args = append(args, data.Latitude)
		args = append(args, data.AddressTypeID)
		args = append(args, userID)
		args = append(args, addressID)

		if err := tx.Exec(ctx, queryUpdate, args...); err != nil {
			return utils.TechnicalError{
				Message: common.MSG_INTERNAL_ERROR,
				Code:    http.StatusInternalServerError,
			}
		}

		return nil
	})
}

func (a *addressRepository) DeleteAddressByID(ctx context.Context, addressID int) error {
	ctx, span := a.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "DeleteAddressByID"))
	defer span.End()

	// select to check it is default address or not
	querySelect := `SELECT is_default from addresses WHERE id = $1`
	var isDefault bool

	if err := a.db.QueryRow(ctx, querySelect, addressID).Scan(&isDefault); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.BusinessError{
				Message:   "Address not found",
				Code:      http.StatusNotFound,
				ErrorCode: errorcode.NOT_FOUND,
			}
		}

		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	if isDefault {
		return utils.BusinessError{
			Message:   "You can not delete this address default",
			Code:      http.StatusBadRequest,
			ErrorCode: errorcode.BAD_REQUEST,
		}
	}

	queryDelete := `DELETE FROM addresses WHERE id = $1`
	if err := a.db.Exec(ctx, queryDelete, addressID); err != nil {
		return utils.TechnicalError{
			Message: common.MSG_INTERNAL_ERROR,
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}
