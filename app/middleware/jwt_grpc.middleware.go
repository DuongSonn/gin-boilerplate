package middleware

import (
	"context"
	"oauth-server/package/errors"
	_jwt "oauth-server/package/jwt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type jwtRPCMiddleware struct{}

func newJWTgRPCMiddleware() RPCMiddleware {
	return &jwtRPCMiddleware{}
}

// build jwt middleware for grpc
func (m *jwtRPCMiddleware) Handler(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, errors.GetMessage(errors.ErrCodeUnauthorized))
	}

	authorization := md["authorization"]
	if len(authorization) < 1 {
		return nil, status.Error(codes.Unauthenticated, errors.GetMessage(errors.ErrCodeUnauthorized))
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	payload, err := _jwt.ParseRPCToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, errors.GetMessage(errors.ErrCodeUnauthorized))
	}
	ctx = context.WithValue(ctx, _jwt.RPC_CONTEXT_KEY, payload)

	return handler(ctx, req)
}
