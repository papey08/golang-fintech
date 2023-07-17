package httpgin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"homework9/internal/app"
	"log"
	"net/http"
	"time"
)

// NewHTTPServer creates http.Server with routes and middlewares
func NewHTTPServer(port int, app app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("api/v1")
	AppRouter(api, app)
	return &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
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
