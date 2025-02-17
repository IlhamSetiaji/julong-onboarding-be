package route

import (
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/handler"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouteConfig struct {
	App               *gin.Engine
	Log               *logrus.Logger
	Viper             *viper.Viper
	AuthMiddleware    gin.HandlerFunc
	UniversityHandler handler.IUniversityHandler
}

func (c *RouteConfig) SetupRoutes() {
	c.App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	c.App.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	c.SetupAPIRoutes()
}

func (c *RouteConfig) SetupAPIRoutes() {
	apiRoute := c.App.Group("/api")
	{
		apiRoute.Use(c.AuthMiddleware)
		{
			// universities
			universityRoute := apiRoute.Group("/universities")
			{
				universityRoute.GET("", c.UniversityHandler.FindAll)
			}
		}
	}
}

func NewRouteConfig(app *gin.Engine, viper *viper.Viper, log *logrus.Logger) *RouteConfig {
	authMiddleware := middleware.NewAuth(viper)
	universityHandler := handler.UniversityHandlerFactory(log, viper)
	return &RouteConfig{
		App:               app,
		Log:               log,
		Viper:             viper,
		AuthMiddleware:    authMiddleware,
		UniversityHandler: universityHandler,
	}
}
