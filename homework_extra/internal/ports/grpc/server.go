package grpc

import (
	"context"
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework_extra/internal/app"
	"homework_extra/internal/ports/grpc/pb"
	"log"
	"net"
)

type Server struct {
	pb.AdServiceServer
	App app.App
}

// NewGRPCServer creates server with interceptors
func NewGRPCServer(a app.App) *grpc.Server {
	chain := grpcmiddleware.ChainUnaryServer(
		ServerLoggerInterceptor(),
		ServerPanicInterceptor(),
	)
	s := grpc.NewServer(grpc.UnaryInterceptor(chain))
	pb.RegisterAdServiceServer(s, &Server{
		AdServiceServer: nil,
		App:             a,
	})
	return s
}

// GracefulShutdown shutdowns server on signals from sigQuit
func GracefulShutdown(ctx context.Context, srv *grpc.Server, lis net.Listener, eg *errgroup.Group) {
	eg.Go(func() error {
		log.Printf("GRPC server started\n")
		defer log.Printf("GRPC server closed\n")

		errCh := make(chan error)

		defer func() {
			srv.GracefulStop()
			_ = lis.Close()
			close(errCh)
		}()

		go func() {
			if err := srv.Serve(lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("GRPC server unable to listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("GRPC server shutting down: %s\n", err.Error())
	}
	log.Println("GRPC server was successfully shutdown")
}
