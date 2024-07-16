package http

import (
	"context"
	"net"

	"github.com/gin-gonic/gin"

	"github.com/silvan-talos/tlp/example/user"
	"github.com/silvan-talos/tlp/log"
)

type Server struct {
	router *gin.Engine
}

func NewServer(us *user.Service) *Server {
	if us == nil {
		panic("user service is nil")
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	userRoutes := r.Group("/users")
	{
		h := userHandler{us: us}
		h.addRoutes(userRoutes)
	}

	return &Server{
		router: r,
	}
}

func (s *Server) Serve(lis net.Listener) error {
	log.Info(context.Background(), "Starting http server", "address", lis.Addr().String())
	return s.router.RunListener(lis)
}
