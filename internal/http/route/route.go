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
	App                           *gin.Engine
	Log                           *logrus.Logger
	Viper                         *viper.Viper
	AuthMiddleware                gin.HandlerFunc
	UniversityHandler             handler.IUniversityHandler
	CoverHandler                  handler.ICoverHandler
	TemplateTaskHandler           handler.ITemplateTaskHandler
	TemplateTaskAttachmentHandler handler.ITemplateTaskAttachmentHandler
}

func (c *RouteConfig) SetupRoutes() {
	c.App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to Julong Onboarding API",
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
			// covers
			coverRoute := apiRoute.Group("/covers")
			{
				coverRoute.GET("", c.CoverHandler.FindAllPaginated)
				coverRoute.GET("/:id", c.CoverHandler.FindByID)
				coverRoute.POST("", c.CoverHandler.CreateCover)
				coverRoute.POST("/upload", c.CoverHandler.UploadCover)
				coverRoute.PUT("/update", c.CoverHandler.UpdateCover)
				coverRoute.DELETE("/:id", c.CoverHandler.DeleteCover)
			}
			// template tasks
			templateTaskRoute := apiRoute.Group("/template-tasks")
			{
				templateTaskRoute.GET("", c.TemplateTaskHandler.FindAllPaginated)
				templateTaskRoute.GET("/:id", c.TemplateTaskHandler.FindByID)
				templateTaskRoute.POST("", c.TemplateTaskHandler.CreateTemplateTask)
				templateTaskRoute.PUT("/update", c.TemplateTaskHandler.UpdateTemplateTask)
				templateTaskRoute.DELETE("/:id", c.TemplateTaskHandler.DeleteTemplateTask)
			}
			// template task attachments
			templateTaskAttachmentRoute := apiRoute.Group("/template-task-attachments")
			{
				templateTaskAttachmentRoute.GET("/:id", c.TemplateTaskAttachmentHandler.FindByID)
				templateTaskAttachmentRoute.DELETE("/:id", c.TemplateTaskAttachmentHandler.DeleteTemplateTaskAttachment)
			}
		}
	}
}

func NewRouteConfig(app *gin.Engine, viper *viper.Viper, log *logrus.Logger) *RouteConfig {
	authMiddleware := middleware.NewAuth(viper)
	universityHandler := handler.UniversityHandlerFactory(log, viper)
	coverHandler := handler.CoverHandlerFactory(log, viper)
	templateTaskHandler := handler.TemplateTaskHandlerFactory(log, viper)
	templateTaskAttachmentHandler := handler.TemplateTaskAttachmentHandlerFactory(log, viper)
	return &RouteConfig{
		App:                           app,
		Log:                           log,
		Viper:                         viper,
		AuthMiddleware:                authMiddleware,
		UniversityHandler:             universityHandler,
		CoverHandler:                  coverHandler,
		TemplateTaskHandler:           templateTaskHandler,
		TemplateTaskAttachmentHandler: templateTaskAttachmentHandler,
	}
}
