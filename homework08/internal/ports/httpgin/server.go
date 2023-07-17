package httpgin

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, a app.App, logsWriter io.Writer) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{
		port: port,
		app:  gin.New(),
	}
	api := s.app.Group("api/v1")
	AppRouter(api, a, logsWriter)
	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
