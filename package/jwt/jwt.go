package _jwt

import (
	"oauth-server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RPCPayload struct {
	UserID uuid.UUID `json:"user_id"`
	Scope  string    `json:"scopes"`
	jwt.RegisteredClaims
}
type UserPayload struct {
	ID uuid.UUID `json:"id"`
}

type key string

const (
	RPC_CONTEXT_KEY key = "rpc_user"
)

func ParseRPCToken(tokenStr string) (*RPCPayload, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &RPCPayload{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RPCPayload); ok {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

func GenerateRPCToken(payload RPCPayload) (string, error) {
	conf := config.GetConfiguration().Jwt

	accessToken, err := GenerateToken(payload, conf.UserAccessTokenKey, 24*60*60)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

/*
Parameters:

- payload: The data payload to be included in the token.

- key: The secret key used for signing the token.

- expireTime: The expiration time of the token in seconds.

Returns:

string: The generated token.

error: An error if the token generation fails.
*/
func GenerateToken(payload interface{}, key string, expireTime int) (string, error) {
	conf := config.GetConfiguration().Jwt

	claims := struct {
		data interface{}
		jwt.RegisteredClaims
	}{
		data: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTime))),
			Issuer:    conf.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func VerifyToken(tokenString string, key string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return token, nil
}
