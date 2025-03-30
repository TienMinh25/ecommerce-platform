package api_gateway_service

import (
	"context"
	"crypto/rsa"
	"fmt"
	api_gateway_dto "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/dto"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/TienMinh25/ecommerce-platform/internal/utils/errorcode"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"time"
)

// JwtKeyManager manages private key and public key
type JwtKeyManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	env        *env.EnvManager
	tracer     pkg.Tracer
}

type JwtPayload struct {
	UserID   int
	Email    string
	FullName string
	Role     []api_gateway_dto.RoleLoginResponse
}

type UserClaims struct {
	UserID   int                                 `json:"user_id"`
	Email    string                              `json:"email"`
	FullName string                              `json:"full_name"`
	Role     []api_gateway_dto.RoleLoginResponse `json:"role"`
	jwt.RegisteredClaims
}

func NewJwtService(env *env.EnvManager, tracer pkg.Tracer) (IJwtService, error) {
	privateKeyData, err := os.ReadFile(env.PrivateKeyPath)

	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)

	if err != nil {
		return nil, errors.Wrap(err, "jwt.ParseRSAPrivateKeyFromPEM")
	}

	publicKeyData, err := os.ReadFile(env.PublicKeyPath)

	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)

	if err != nil {
		return nil, errors.Wrap(err, "jwt.ParseRSAPublicKeyFromPEM")
	}

	fmt.Println("✅ Load keys for jwt successfully!")

	return &JwtKeyManager{privateKey: privateKey, publicKey: publicKey, env: env, tracer: tracer}, nil
}

func (km *JwtKeyManager) newUserClaims(userID int, email, fullname string, role []api_gateway_dto.RoleLoginResponse, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()

	if err != nil {
		return nil, utils.TechnicalError{
			Code:    http.StatusInternalServerError,
			Message: common.MSG_INTERNAL_ERROR,
		}
	}

	return &UserClaims{
		UserID:   userID,
		Email:    email,
		FullName: fullname,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}, nil
}

func (km *JwtKeyManager) GenerateToken(ctx context.Context, payload JwtPayload) (string, string, error) {
	ctx, span := km.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "GenerateToken"))
	defer span.End()

	claimsAccessToken, errClaims := km.newUserClaims(payload.UserID, payload.Email, payload.FullName, payload.Role, time.Duration(km.env.ExpireAccessToken)*time.Minute)

	if errClaims != nil {
		return "", "", errClaims
	}

	tokenInfo := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsAccessToken)

	accessToken, err := tokenInfo.SignedString(km.privateKey)

	if err != nil {
		return "", "", utils.TechnicalError{Code: http.StatusInternalServerError, Message: common.MSG_INTERNAL_ERROR}
	}

	refreshToken, err := utils.GenerateRandomString(26)

	if err != nil {
		return "", "", utils.TechnicalError{Code: http.StatusInternalServerError, Message: common.MSG_INTERNAL_ERROR}
	}

	return accessToken, refreshToken, nil
}

func (km *JwtKeyManager) VerifyToken(ctx context.Context, accessToken string) (*UserClaims, error) {
	ctx, span := km.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.ServiceLayer, "VerifyToken"))
	defer span.End()

	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)

		if !ok {
			return nil, utils.BusinessError{Code: http.StatusUnauthorized, Message: "Invalid signing method", ErrorCode: errorcode.UNAUTHORIZED}
		}

		return km.publicKey, nil
	})

	if !token.Valid {
		return nil, utils.BusinessError{Code: http.StatusUnauthorized, Message: "Invalid token", ErrorCode: errorcode.UNAUTHORIZED}
	}

	if err != nil {
		return nil, utils.TechnicalError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	claims, _ := token.Claims.(*UserClaims)

	return claims, nil
}
