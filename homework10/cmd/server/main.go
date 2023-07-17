package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/user_repo"
	"homework10/internal/app"
	grpcserver "homework10/internal/ports/grpc"
	"homework10/internal/ports/httpgin"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func HTTPServerWithGracefulShutdown(ctx context.Context, port int, a app.App, eg *errgroup.Group) {
	srv := httpgin.NewHTTPServer(port, a)
	httpgin.GracefulShutdown(ctx, srv, eg)
}

func GRPCServerWithGracefulShutdown(ctx context.Context, port int, network string, a app.App, eg *errgroup.Group) {
	lis, err := net.Listen(network, fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("unable to create listener: %s\n", err.Error())
	}
	srv := grpcserver.NewGRPCServer(a)

	grpcserver.GracefulShutdown(ctx, srv, lis, eg)
}

// InitConfig initialises configuration file
func InitConfig() error {
	viper.SetConfigFile("config.yml")
	return viper.ReadInConfig()
}

func main() {
	a := app.NewApp(adrepo.New(), user_repo.New())

	err := InitConfig()
	if err != nil {
		log.Printf("config error: %s\n", err.Error())
	}
	httpPort := viper.GetInt("http.port")

	grpcPort := viper.GetInt("grpc.port")
	network := viper.GetString("grpc.network")

	sigQuit := make(chan os.Signal, 1)
	defer close(sigQuit)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	wg := new(sync.WaitGroup)

	// starting http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		HTTPServerWithGracefulShutdown(ctx, httpPort, a, eg)
	}()

	// starting grpc server
	wg.Add(1)
	go func() {
		defer wg.Done()
		GRPCServerWithGracefulShutdown(ctx, grpcPort, network, a, eg)
	}()

	wg.Wait()
}
