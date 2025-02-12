package middleware

type MiddlewareCollections struct {
	RPCMiddleware RPCMiddleware
}

func RegisterMiddleware() MiddlewareCollections {
	return MiddlewareCollections{
		RPCMiddleware: newJWTgRPCMiddleware(),
	}
}
