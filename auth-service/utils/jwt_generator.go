package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	// Extract client ID from context as workaround
	fmt.Printf("GenerateBasic: client = %#v, userID = %s\n", data.Client, data.UserID)

	aud := ""
	if data.Client != nil {
		aud = data.Client.GetID()
	} else if cid := ctx.Value(ClientIDKey); cid != nil {
		aud = cid.(string)
	}

	fmt.Printf("🪪  Token generator input: clientID='%v', userID='%v'\n", aud, data.UserID)

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
		return "", "", err
	}

	if isGenRefresh {
		refresh = uuid.New().String()
	}

	return access, refresh, nil
}
