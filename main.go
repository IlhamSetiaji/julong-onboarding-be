package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/IlhamSetiaji/julong-onboarding-be/docs"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/config"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/rabbitmq"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func generateSwaggerDocs(mode string) {
	var swagPath = "swag"
	if mode == "production" {
		swagPath = "/go/bin/swag"

		// Check if swag exists
		if _, err := os.Stat(swagPath); os.IsNotExist(err) {
			fmt.Println("Swag CLI not found, skipping Swagger generation.")
			return
		}
	}
	cmd := exec.Command(swagPath, "init", "--parseDependency", "--parseInternal")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Generating Swagger documentation...")
	if err := cmd.Run(); err != nil {
		fmt.Println("Error generating Swagger docs:", err)
		os.Exit(1) // Exit if Swagger generation fails
	}
	fmt.Println("Swagger documentation generated successfully!")
}

// @title           Julong Recruitment API Docs
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api

// @securityDefinitions.apikey  BearerAuth
// @in header
// @name Authorization
// @description Bearer token for authentication
// @tokenUrl http://localhost:3000/api/login

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)
	generateSwaggerDocs(viper.GetString("app.env"))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		rabbitmq.InitConsumer(viper, log)
	}()

	go func() {
		defer wg.Done()
		rabbitmq.InitProducer(viper, log)
	}()

	app := gin.Default()
	app.Static("/storage", "./storage")
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("App-Name", viper.GetString("app.name"))
	})

	store := cookie.NewStore([]byte(viper.GetString("web.cookie.secret")))
	app.Use(sessions.Sessions(viper.GetString("web.session.name"), store))

	// setup CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(viper.GetString("frontend.urls"), ","), // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// setup custom csrf middleware
	app.Use(func(c *gin.Context) {
		if !shouldExcludeFromCSRF(c.Request.URL.Path) {
			csrf.Middleware(csrf.Options{
				Secret: viper.GetString("web.csrf_secret"),
				ErrorFunc: func(c *gin.Context) {
					c.String(http.StatusForbidden, "CSRF token mismatch")
					c.Abort()
				},
			})(c)
		}
		c.Next()
	})

	// setup routes
	routeConfig := route.NewRouteConfig(app, viper, log)
	routeConfig.SetupRoutes()

	// run server
	webPort := strconv.Itoa(viper.GetInt("web.port"))
	err := app.Run(":" + webPort)
	if err != nil {
		log.Panicf("Failed to start server: %v", err)
	}

	wg.Wait()
}

func shouldExcludeFromCSRF(path string) bool {
	return len(path) >= 4 && path[:4] == "/api"
}
