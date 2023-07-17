package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

// ServerLoggerInterceptor is my own interceptor for logging on server
func ServerLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		reply, err := handler(ctx, req)
		if err != nil {
			return nil, ErrorToGRPCError(err)
		}
		log.Printf("gRPC\tMethod: %s\tRequest: %v\tResponse: %v\n", info.FullMethod, req, reply)
		return reply, nil
	}
}

// ServerPanicInterceptor is my own interceptor for recovering panics
func ServerPanicInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered: %v", r)
			}
		}()

		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

/*
// ClientLoggerInterceptor is my own interceptor for logging on client
func ClientLoggerInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		log.Printf("gRPC\tMethod: %s\tRequest: %+v\tResponse: %+v\tError:%s", method, req, reply, err)
		if err != nil {
			return ErrorToGRPCError(err)
		} else {
			return nil
		}
	}
}*/
