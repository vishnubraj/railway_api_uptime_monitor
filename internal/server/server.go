package server

import (
	"context"
	"net/http"

	"railway-api-uptime-monitor/internal/config"
	"railway-api-uptime-monitor/internal/database"
	"railway-api-uptime-monitor/internal/handlers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config *config.Config
	server *http.Server
}

func New(db *database.Database, cfg *config.Config) *Server {
	// Set Gin mode
	if cfg.Port != "8080" { // Assume production if not default port
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Initialize handlers
	h := handlers.New(db)

	// Middleware
	router.Use(corsMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Serve static files
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*")

	// Routes
	setupRoutes(router, h)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	return &Server{
		router: router,
		config: cfg,
		server: server,
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func setupRoutes(router *gin.Engine, h *handlers.Handler) {
	// Dashboard
	router.GET("/", h.Dashboard)

	// API routes
	api := router.Group("/api")
	{
		api.GET("/health", h.HealthCheck)
		api.GET("/status", h.GetAllStatus)
		api.GET("/status/:name", h.GetAPIStatus)
		api.GET("/logs/:name", h.GetAPILogs)
		api.GET("/alerts", h.GetAlerts)
		api.GET("/stats", h.GetStats)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
