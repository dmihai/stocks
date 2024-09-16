package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/dmihai/stocks/pkg/data"
)

type Server struct {
	addr string
	data *data.Store
}

func NewServer(addr string, data *data.Store) *Server {
	return &Server{
		addr: addr,
		data: data,
	}
}

func (s *Server) Start() {
	r := s.setupRouter()

	r.Run(s.addr)
}

func (s *Server) setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ping", ping)
	r.GET("/top-gainers", s.getTopGainers)

	return r
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (s *Server) getTopGainers(c *gin.Context) {
	topGainers := s.data.GetTopGainers(20)
	c.JSON(http.StatusOK, topGainers)
}
