package xgrpc

import (
	"context"

	"github.com/davidterranova/contacts/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func BasicAuthMiddleware(username string, password string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, auth.ErrUnauthorized
		}

		authMetadata := meta.Get("Authorization")
		if len(authMetadata) == 0 {
			return nil, auth.ErrUnauthorized
		}

		user, err := auth.BasicAuth(username, password)(authMetadata[0])
		if err != nil {
			return nil, auth.ErrUnauthorized
		}

		return handler(auth.ContextWithUser(ctx, user), req)
	}
}

func GrantAnyFn() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, auth.ErrUnauthorized
		}

		authMetadata := meta.Get("Authorization")
		if len(authMetadata) == 0 {
			return nil, auth.ErrUnauthorized
		}

		user, err := auth.GrantAnyAccess()(authMetadata[0])
		if err != nil {
			return nil, auth.ErrUnauthorized
		}

		return handler(auth.ContextWithUser(ctx, user), req)
	}
}
