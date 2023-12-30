package httpgin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"homework_extra/internal/app"
	"log"
	"net/http"
	"strconv"
	"time"
)

// NewHTTPServer creates http.Server with routes and middlewares
func NewHTTPServer(host string, port int, app app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("api/v1")
	AppRouter(api, app)
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}
}

// GracefulShutdown shutdowns server on signals from sigQuit
func GracefulShutdown(ctx context.Context, srv *http.Server, eg *errgroup.Group) {
	eg.Go(func() error {
		log.Printf("HTTP server started\n")
		defer log.Printf("HTTP server closed\n")

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(shCtx); err != nil {
				log.Printf("can't close http server listening on %s: %s", srv.Addr, err.Error())
			}
			close(errCh)
		}()

		go func() {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("HTTP server unable to listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("HTTP server shutting down: %s\n", err.Error())
	}

	log.Println("HTTP server was successfully shutdown")
}

// getID returns int64 given in endpoint string by key or error if endpoint is invalid
func getID(c *gin.Context, key string) (int64, error) {
	return strToID(c.Param(key))
}

// strToID exists only for fuzz test
func strToID(strID string) (int64, error) {
	if id, err := strconv.Atoi(strID); err != nil {
		return -1, err
	} else {
		return int64(id), err
	}
}