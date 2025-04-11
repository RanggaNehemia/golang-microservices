package utils

import (
	"context"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const ClientIDKey contextKey = "client_id"

type CustomJWTAccessGenerate struct {
	SignedKey     []byte
	SigningMethod jwt.SigningMethod
}

func NewCustomJWTAccessGenerate(key []byte, method jwt.SigningMethod) *CustomJWTAccessGenerate {
	return &CustomJWTAccessGenerate{
		SignedKey:     key,
		SigningMethod: method,
	}
}

func (cg *CustomJWTAccessGenerate) Token(
	ctx context.Context,
	data *oauth2.GenerateBasic,
	isGenRefresh bool,
) (access, refresh string, err error) {
	aud := ""
	if data.Client != nil {
		aud = data.Client.GetID()
	} else if cid := ctx.Value(ClientIDKey); cid != nil {
		aud = cid.(string)
	}

	claims := jwt.MapClaims{
		"iss":   "https://auth.sampledomain.com",
		"sub":   data.UserID,
		"aud":   aud,
		"iat":   time.Now().Unix(),
		"exp":   data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()).Unix(),
		"jti":   uuid.New().String(),
		"scope": data.TokenInfo.GetScope(),
	}

	token := jwt.NewWithClaims(cg.SigningMethod, claims)
	access, err = token.SignedString(cg.SignedKey)
	if err != nil {
		Logger.Error("Error on generating JWT Token", zap.Error(err))
		return "", "", err
	}

	if isGenRefresh {
		refresh = uuid.New().String()
	}

	Logger.Info("JWT Token generated")
	return access, refresh, nil
}
