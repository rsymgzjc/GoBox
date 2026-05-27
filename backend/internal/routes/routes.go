package routes

import (
	"gobox/backend/internal/config"
	"gobox/backend/internal/controllers"
	"gobox/backend/internal/middleware"
	"gobox/backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, logger *zap.Logger, db *gorm.DB) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	if cfg.App.Env == "development" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(middleware.RequestLogger(logger))

	authService := services.NewAuthService(db, cfg, logger)
	toolService := services.NewToolService(db)

	authController := controllers.NewAuthController(authService)
	toolController := controllers.NewToolController(toolService)

	router.GET("/health", controllers.Health)

	api := router.Group("/api/v1")
	{
		api.POST("/auth/register/send-code", authController.SendRegisterCode)
		api.POST("/auth/register", authController.Register)
		api.POST("/auth/login", authController.Login)
		api.GET("/tools", toolController.List)
		api.POST("/tools/:slug/run", toolController.Run)
		api.GET("/stats/summary", toolController.Summary)
	}

	protected := api.Group("")
	protected.Use(middleware.Auth(cfg))
	{
		protected.GET("/me", authController.Profile)
		protected.PUT("/me/preferences", authController.SavePreferences)
	}

	admin := protected.Group("/admin")
	admin.Use(middleware.AdminOnly())
	{
		admin.GET("/summary", toolController.Summary)
	}

	return router, nil
}
