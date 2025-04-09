package api_gateway_handler

import (
	"context"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	api_gateway_service "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/service"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userManagementHandler struct {
	tracer  pkg.Tracer
	service api_gateway_service.IUserService
}

func NewUserManagementHandler(tracer pkg.Tracer, service api_gateway_service.IUserService) IUserManagementHandler {
	return &userManagementHandler{
		tracer:  tracer,
		service: service,
	}
}

// GetUserManagement implements IUserManagementHandler.
// GetUserManagement godoc
//
//	@Summary		Get list users
//	@Tags			users
//	@Description	Get list users
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			limit				query		int		true	"Limit number of records returned"
//	@Param			page				query		int		true	"page"
//
//	@Param			searchBy			query		string	false	"Provide for search by email, phone or fullname"	Enums(email, phone, fullname)
//	@Param			searchValue			query		string	false	"Value for search"
//	@Param			sortBy				query		string	false	"Sort by some attributes"	Enums(fullname, email, birthdate, updated_at, phone)
//	@Param			sortOrder			query		string	false	"Sort order asc or desc"	Enums(asc, desc)
//	@Param			emailVerify			query		boolean	false	"Filter by email verify"	Enums(true, false)
//	@Param			phoneVerify			query		boolean	false	"Filter by phone verify"	Enums(true, false)
//	@Param			status				query		string	false	"Filter by status of user"	Enums(active, inactive)
//	@Param			updatedAtStartFrom	query		string	false	"Start time for last update"
//	@Param			updatedAtEndFrom	query		string	false	"End time for last update"
//	@Param			roleID				query		integer	false	"Filter by role id"	Enums(1, 2, 3, 4)
//	@Success		200					{object}	api_gateway_dto.GetUserByAdminResponseDocs
//	@Failure		400					{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401					{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500					{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users [get]
func (h *userManagementHandler) GetUserManagement(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "GetUserManagement"))
	defer span.End()

	var data api_gateway_dto.GetUserByAdminRequest
	if err := ctx.ShouldBindQuery(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	res, totalItems, totalPages, hasNext, hasPrevious, err := h.service.GetUserManagement(ct, &data)

	if err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.PaginatedResponse[[]api_gateway_dto.GetUserByAdminResponse](ctx, res, data.Page, data.Limit, totalPages, totalItems, hasNext, hasPrevious)
}

// CreateUser implements IUserManagementHandler.
// CreateUser godoc
//
//	@Summary		create new user by admin
//	@Tags			users
//	@Description	create new user by admin
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			request	body		api_gateway_dto.CreateUserByAdminRequest	true	"Request body"
//	@Success		200		{object}	api_gateway_dto.CreateUserByAdminResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users [post]
func (h *userManagementHandler) CreateUser(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "CreateUser"))
	defer span.End()

	var data api_gateway_dto.CreateUserByAdminRequest

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := h.service.CreateUserByAdmin(ct, &data); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.CreateUserByAdminResponse{})
}

// UpdateUser implements IUserManagementHandler.
// UpdateUser godoc
//
//	@Summary		update user by admin
//	@Tags			users
//	@Description	update user by admin
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			userID	path		int											true	"user id on path"
//	@Param			request	body		api_gateway_dto.UpdateUserByAdminRequest	true	"Request body"
//	@Success		200		{object}	api_gateway_dto.UpdateUserByAdminResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/:userID [patch]
func (h *userManagementHandler) UpdateUser(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	_, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "UpdateUser"))
	defer span.End()

	var data api_gateway_dto.UpdateUserByAdminRequest
	var uri api_gateway_dto.UpdateUserByAdminRequestURI

	if err := ctx.ShouldBindUri(&uri); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := h.service.UpdateUserByAdmin(ctx, &data, uri.UserID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.UpdateUserByAdminResponse{})
}

// DeleteUserByID implements IUserManagementHandler.
// DeleteUserByID godoc
//
//	@Summary		delete user by admin
//	@Tags			users
//	@Description	delete user by admin
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//	@Param			userID	path		int	true	"user id on path"
//	@Success		200		{object}	api_gateway_dto.DeleteUserByAdminResponseDocs
//	@Failure		400		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		401		{object}	api_gateway_dto.ResponseErrorDocs
//	@Failure		500		{object}	api_gateway_dto.ResponseErrorDocs
//	@Router			/users/:userID [delete]
func (h *userManagementHandler) DeleteUserByID(ctx *gin.Context) {
	cRaw, _ := ctx.Get("tracingContext")
	c := cRaw.(context.Context)
	ct, span := h.tracer.StartFromContext(c, tracing.GetSpanName(tracing.HandlerLayer, "DeleteUserByID"))
	defer span.End()

	var data api_gateway_dto.DeleteUserByAdminRequest

	if err := ctx.ShouldBindUri(&data); err != nil {
		span.RecordError(err)
		utils.HandleValidateData(ctx, err)
		return
	}

	if err := h.service.DeleteUserByAdmin(ct, data.UserID); err != nil {
		span.RecordError(err)
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, api_gateway_dto.DeleteUserByAdminResponse{})
}
