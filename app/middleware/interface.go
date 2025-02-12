package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Middleware interface {
	Handler() gin.HandlerFunc
}

type RPCMiddleware interface {
	Handler(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
}
