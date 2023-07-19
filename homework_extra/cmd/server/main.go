package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"homework_extra/internal/adapters/adrepo"
	"homework_extra/internal/adapters/user_repo"
	"homework_extra/internal/app"
	grpcserver "homework_extra/internal/ports/grpc"
	"homework_extra/internal/ports/httpgin"
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
	err := InitConfig()
	if err != nil {
		log.Printf("config error: %s\n", err.Error())
		return
	}
	httpPort := viper.GetInt("http.port")

	grpcPort := viper.GetInt("grpc.port")
	network := viper.GetString("grpc.network")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		viper.GetString("adrepo.username"),
		viper.GetString("adrepo.password"),
		viper.GetString("adrepo.host"),
		viper.GetString("adrepo.port"),
		viper.GetString("adrepo.dbname"),
		viper.GetString("adrepo.sslmode"))

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Printf("adrepo connection error: %s\n", err.Error())
		return
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		_ = conn.Close(ctx)
	}(conn, context.Background())

	a := app.NewApp(adrepo.New(conn), user_repo.New())

	// configuring graceful shutdown
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
