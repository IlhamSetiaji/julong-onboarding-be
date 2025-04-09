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
	EmployeeTaskHandler           handler.IEmployeeTaskHandler
	EmployeeTaskAttachmentHandler handler.IEmployeeTaskAttachmentHandler
	EventHandler                  handler.IEventHandler
	AnswerTypeHandler             handler.IAnswerTypeHandler
	SurveyTemplateHandler         handler.ISurveyTemplateHandler
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
			// employee tasks
			employeeTaskRoute := apiRoute.Group("/employee-tasks")
			{
				employeeTaskRoute.GET("", c.EmployeeTaskHandler.FindAllPaginated)
				employeeTaskRoute.GET("/employee-kanban", c.EmployeeTaskHandler.FindAllByEmployeeIDAndKanbanPaginated)
				employeeTaskRoute.GET("/employee", c.EmployeeTaskHandler.FindAllByEmployeeID)
				employeeTaskRoute.GET("/employee-paginated", c.EmployeeTaskHandler.FindAllPaginatedByEmployeeID)
				employeeTaskRoute.GET("/count", c.EmployeeTaskHandler.CountByKanbanAndEmployeeID)
				employeeTaskRoute.GET("/employee-kanban/count", c.EmployeeTaskHandler.CountKanbanProgressByEmployeeID)
				employeeTaskRoute.GET("/:id", c.EmployeeTaskHandler.FindByID)
				employeeTaskRoute.POST("", c.EmployeeTaskHandler.CreateEmployeeTask)
				employeeTaskRoute.PUT("/update", c.EmployeeTaskHandler.UpdateEmployeeTask)
				employeeTaskRoute.DELETE("/:id", c.EmployeeTaskHandler.DeleteEmployeeTask)
			}
			// employee task attachments
			employeeTaskAttachmentRoute := apiRoute.Group("/employee-task-attachments")
			{
				employeeTaskAttachmentRoute.GET("/:id", c.EmployeeTaskAttachmentHandler.FindByID)
				employeeTaskAttachmentRoute.DELETE("/:id", c.EmployeeTaskAttachmentHandler.DeleteEmployeeTaskAttachment)
			}
			// events
			eventRoute := apiRoute.Group("/events")
			{
				eventRoute.GET("", c.EventHandler.FindAllPaginated)
				eventRoute.GET("/:id", c.EventHandler.FindByID)
				eventRoute.POST("", c.EventHandler.CreateEvent)
				eventRoute.PUT("/update", c.EventHandler.UpdateEvent)
				eventRoute.DELETE("/:id", c.EventHandler.DeleteEvent)
			}
			// answer types
			answerTypeRoute := apiRoute.Group("/answer-types")
			{
				answerTypeRoute.GET("", c.AnswerTypeHandler.FindAll)
			}
			// survey templates
			surveyTemplateRoute := apiRoute.Group("/survey-templates")
			{
				surveyTemplateRoute.GET("", c.SurveyTemplateHandler.FindAllSurveyTemplatesPaginated)
				surveyTemplateRoute.GET("/:id", c.SurveyTemplateHandler.FindSurveyTemplateByID)
				surveyTemplateRoute.POST("", c.SurveyTemplateHandler.CreateSurveyTemplate)
				surveyTemplateRoute.PUT("/update", c.SurveyTemplateHandler.UpdateSurveyTemplate)
				surveyTemplateRoute.DELETE("/:id", c.SurveyTemplateHandler.DeleteSurveyTemplate)
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
	employeeTaskHandler := handler.EmployeeTaskHandlerFactory(log, viper)
	employeeTaskAttachmentHandler := handler.EmployeeTaskAttachmentHandlerFactory(log, viper)
	eventHandler := handler.EventHandlerFactory(log, viper)
	answerTypeHandler := handler.AnswerTypeHandlerFactory(log, viper)
	surveyTemplateHandler := handler.SurveyTemplateHandlerFactory(log, viper)
	return &RouteConfig{
		App:                           app,
		Log:                           log,
		Viper:                         viper,
		AuthMiddleware:                authMiddleware,
		UniversityHandler:             universityHandler,
		CoverHandler:                  coverHandler,
		TemplateTaskHandler:           templateTaskHandler,
		TemplateTaskAttachmentHandler: templateTaskAttachmentHandler,
		EmployeeTaskHandler:           employeeTaskHandler,
		EmployeeTaskAttachmentHandler: employeeTaskAttachmentHandler,
		EventHandler:                  eventHandler,
		AnswerTypeHandler:             answerTypeHandler,
		SurveyTemplateHandler:         surveyTemplateHandler,
	}
}
