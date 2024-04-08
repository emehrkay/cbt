package ctx

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var token string
	ctxToken, err := ContextGetToken(ctx)
	if err == nil {
		token = *ctxToken
	}

	ctx = metadata.NewOutgoingContext(ctx,
		metadata.New(
			map[string]string{
				"jwt": token,
			},
		),
	)

	return invoker(ctx, method, req, reply, cc, opts...)
}
