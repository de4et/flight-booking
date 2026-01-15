package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/de4et/flight-booking/internal/server/handlers"
	"github.com/de4et/flight-booking/internal/server/middleware"
	"github.com/de4et/flight-booking/internal/service"
)

func (s *Server) RegisterRoutes(searchService *service.MultipleSearchService) http.Handler {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.LogHandler())
	r.Use(middleware.MetricsHandler())
	r.Use(middleware.ErrorHandler())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	apiGroup := r.Group("/api/v1")
	apiGroup.GET("/search-result", handlers.NewSearchResultHandler(searchService).Handle)

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World!!?"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
