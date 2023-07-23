package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	"time"
)

func HTTPServerWithGracefulShutdown(ctx context.Context, host string, port int, a app.App, eg *errgroup.Group) {
	srv := httpgin.NewHTTPServer(host, port, a)
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

// InitConfig initializes configuration file
func InitConfig() error {
	viper.SetConfigFile("config.yml")
	return viper.ReadInConfig()
}

// AdRepoConfig initializes connection to ads database
func AdRepoConfig(ctx context.Context, dbURL string) *pgx.Conn {
	// connecting to a database in the loop with delay 1 sec for correct starting in docker container
	for {
		conn, err := pgx.Connect(ctx, dbURL)
		if err != nil { // database haven't initialized in docker container yet
			log.Printf("adrepo connection error: %s\n", err.Error())
			time.Sleep(time.Second)
		} else { // database already initialized
			return conn
		}
	}
}

// UserRepoConfig initializes connection to users collection and returns last added ID
func UserRepoConfig(ctx context.Context, dbURL string, dbName string, collectionName string) (*mongo.Client, *mongo.Collection, int64, error) {
	for {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
		if err != nil {
			log.Printf("user_repo connection error: %s\n", err.Error())
			time.Sleep(time.Second)
		} else {
			// тут костыль для того, чтобы поле ID коллекции users сохраняло
			// уникальность даже при перезапуске, вообще конечно postgresql
			// подошла бы лучше, mongodb тут чисто по приколу

			collection := client.Database(dbName).Collection(collectionName)
			currentIDResult := collection.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.D{{"id", -1}}))
			if currentIDResult.Err() != nil && currentIDResult.Err() != mongo.ErrNoDocuments {
				return nil, nil, 0, currentIDResult.Err()
			}

			var user user_repo.UsersField
			err = currentIDResult.Decode(&user)
			if err == nil {
				return client, collection, user.ID, nil
			} else if err == mongo.ErrNoDocuments {
				return client, collection, 0, nil
			} else {
				return nil, nil, 0, err
			}
		}
	}
}

func main() {
	ctx := context.Background()
	err := InitConfig()
	if err != nil {
		log.Fatal("config error:", err)
	}
	httpHost := viper.GetString("http.host")
	httpPort := viper.GetInt("http.port")

	grpcPort := viper.GetInt("grpc.port")
	network := viper.GetString("grpc.network")

	// configuring adrepo
	adRepoURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		viper.GetString("adrepo.username"),
		viper.GetString("adrepo.password"),
		viper.GetString("adrepo.host"),
		viper.GetString("adrepo.port"),
		viper.GetString("adrepo.dbname"),
		viper.GetString("adrepo.sslmode"))
	adrepoConn := AdRepoConfig(ctx, adRepoURL)
	defer func(ctx context.Context, conn *pgx.Conn) {
		_ = conn.Close(ctx)
	}(ctx, adrepoConn)

	// configuring user_repo
	userRepoURL := fmt.Sprintf("mongodb://%s:%s",
		viper.GetString("user_repo.host"),
		viper.GetString("user_repo.port"))
	urCli, urColl, urCurrID, err := UserRepoConfig(ctx, userRepoURL, viper.GetString("user_repo.dbname"), viper.GetString("user_repo.collection"))
	if err != nil {
		log.Fatal("user_repo connection error:", err)
	}
	defer func(ctx context.Context, client *mongo.Client) {
		if err = urCli.Disconnect(ctx); err != nil {
			log.Fatal("user_repo disconnect error:", err)
		}
	}(ctx, urCli)

	a := app.NewApp(adrepo.New(adrepoConn), user_repo.New(urColl, urCurrID))

	// configuring graceful shutdown
	sigQuit := make(chan os.Signal, 1)
	defer close(sigQuit)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(ctx)
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
		HTTPServerWithGracefulShutdown(ctx, httpHost, httpPort, a, eg)
	}()

	// starting grpc server
	wg.Add(1)
	go func() {
		defer wg.Done()
		GRPCServerWithGracefulShutdown(ctx, grpcPort, network, a, eg)
	}()

	wg.Wait()
}
