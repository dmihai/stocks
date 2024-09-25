package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/dmihai/stocks/pkg/auth"
	"github.com/dmihai/stocks/pkg/data"
)

type Server struct {
	addr string
	auth *auth.Auth
	data *data.Store
}

func NewServer(addr string, auth *auth.Auth, data *data.Store) *Server {
	return &Server{
		addr: addr,
		auth: auth,
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

	r.POST("/login", gin.BasicAuth(s.auth.Accounts), s.login)
	r.GET("/top-gainers", s.getTopGainers)

	return r
}

func (s *Server) getTopGainers(c *gin.Context) {
	topGainers := s.data.GetTopGainers(20)
	c.JSON(http.StatusOK, topGainers)
}

func (s *Server) login(c *gin.Context) {
	user := c.GetString(gin.AuthUserKey)

	token, err := s.auth.GenerateJWT(user)
	if err != nil {
		s.error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *Server) error(c *gin.Context, err error) {
	c.AbortWithError(http.StatusInternalServerError, err)
}
