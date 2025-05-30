package repository

import (
	"context"
	"encoding/json"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/order-and-payment/grpc/proto/order_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type delivererRepository struct {
	tracer pkg.Tracer
	db     pkg.Database
}

func NewDelivererRepository(tracer pkg.Tracer, db pkg.Database) IDelivererRepository {
	return &delivererRepository{
		tracer: tracer,
		db:     db,
	}
}

func (r *delivererRepository) RegisterDeliverer(ctx context.Context, data *order_proto_gen.RegisterDelivererRequest) error {
	ctx, span := r.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.RepositoryLayer, "RegisterDeliverer"))
	defer span.End()

	return r.db.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pkg.Tx) error {
		// check exists in delivery persons
		queryCheckDeliverer := `select exists (select 1 from delivery_persons where user_id = $1)`

		var isExists bool

		if err := tx.QueryRow(ctx, queryCheckDeliverer, data.UserId).Scan(&isExists); err != nil {
			return status.Errorf(codes.Internal, "check delivery person existence: %v", err)
		}

		if isExists {
			return status.Errorf(codes.AlreadyExists, "delivery person already exists")
		}

		querySelect := `select application_status from delivery_person_applications where user_id = $1`

		rows, err := tx.Query(ctx, querySelect, data.UserId)

		if err != nil {
			return status.Errorf(codes.Internal, "select delivery person application: %v", err)
		}

		defer rows.Close()

		for rows.Next() {
			var oldStatus common.DeliveryPersonApplicationStatus

			if err = rows.Scan(&oldStatus); err != nil {
				return status.Errorf(codes.Internal, "scan delivery person application: %v", err)
			}

			if oldStatus == common.DeliveryPersonApplicationStatusPending {
				return status.Errorf(codes.AlreadyExists, "delivery person application is pending")
			}
		}

		// insert into delivery_person_applications
		pgBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		serviceArea := struct {
			Country  string `json:"country"`
			City     string `json:"city"`
			District string `json:"district"`
			Ward     string `json:"ward"`
		}{
			Country:  data.ServiceArea.Country,
			City:     data.ServiceArea.City,
			District: data.ServiceArea.District,
			Ward:     data.ServiceArea.Ward,
		}

		rawServiceArea, err := json.Marshal(serviceArea)

		if err != nil {
			return status.Errorf(codes.Internal, "marshal service area: %v", err)
		}

		sqlInsert, args, err := pgBuilder.Insert("delivery_person_applications").
			Columns("user_id", "id_card_number", "id_card_front_image", "id_card_back_image",
				"vehicle_type", "vehicle_license_plate", "service_area", "application_status").
			Values(data.UserId, data.IdCardNumber, data.IdCardFrontImage, data.IdCardBackImage,
				data.VehicleType, data.VehicleLicensePlate, rawServiceArea, common.DeliveryPersonApplicationStatusPending).
			ToSql()

		if err != nil {
			span.RecordError(err)
			return status.Errorf(codes.Internal, "build sql query: %v", err)
		}

		if err = tx.Exec(ctx, sqlInsert, args...); err != nil {
			span.RecordError(err)
			return status.Errorf(codes.Internal, "execute sql query: %v", err)
		}

		return nil
	})
}
