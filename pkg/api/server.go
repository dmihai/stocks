package api

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/dmihai/stocks/pkg/auth"
	"github.com/dmihai/stocks/pkg/data"
	"github.com/dmihai/stocks/pkg/stocks"
)

const (
	authKey = "auth"
)

type Server struct {
	addr   string
	auth   *auth.Auth
	data   *data.Store
	stocks stocks.Client
}

func NewServer(addr string, auth *auth.Auth, data *data.Store, stocks stocks.Client) *Server {
	return &Server{
		addr:   addr,
		auth:   auth,
		data:   data,
		stocks: stocks,
	}
}

func (s *Server) Start() {
	r := s.setupRouter()

	r.Run(s.addr)
}

func (s *Server) setupRouter() *gin.Engine {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	r.Use(cors.New(corsConfig))

	r.POST("/login", gin.BasicAuth(s.auth.Accounts), s.login)
	r.POST("/exchange", s.extractBearer, s.exchange, s.login)
	api := r.Group("/api", s.extractBearer, s.validateAuth)
	api.GET("/top-gainers", s.getTopGainers)
	api.GET("/symbol-details/:symbol", s.getSymbolDetails)

	return r
}

func (s *Server) getTopGainers(c *gin.Context) {
	topGainers := s.data.GetTopGainers(20)
	c.JSON(http.StatusOK, topGainers)
}

func (s *Server) getSymbolDetails(c *gin.Context) {
	symbol := c.Param("symbol")

	details, err := s.stocks.GetSymbolDetails(symbol)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if details == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, details)
}

func (s *Server) extractBearer(c *gin.Context) {
	bearerToken := c.Request.Header.Get("Authorization")
	bearerTokenParts := strings.Split(bearerToken, " ")

	if len(bearerTokenParts) != 2 || strings.ToLower(bearerTokenParts[0]) != "bearer" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Set(authKey, bearerTokenParts[1])
}

func (s *Server) validateAuth(c *gin.Context) {
	token := c.GetString(authKey)

	username, err := s.auth.ParseJWT(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(gin.AuthUserKey, username)
}

func (s *Server) exchange(c *gin.Context) {
	refreshToken := c.GetString(authKey)

	username := s.auth.FindUserByRefreshToken(refreshToken)
	if username == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(gin.AuthUserKey, username)
}

func (s *Server) login(c *gin.Context) {
	user := c.GetString(gin.AuthUserKey)

	token, err := s.auth.GenerateJWT(user)
	if err != nil {
		s.error(c, err)
		return
	}

	refreshToken := s.auth.GenerateRefreshToken(user)

	c.JSON(http.StatusOK, gin.H{
		"token":        token,
		"refreshToken": refreshToken,
	})
}

func (s *Server) error(c *gin.Context, err error) {
	c.AbortWithError(http.StatusInternalServerError, err)
}
