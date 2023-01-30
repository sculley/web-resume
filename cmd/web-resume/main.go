package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/sculley/someadmin-go/config"
	"github.com/sculley/web-resume/internal/middleware"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

// Experience is the struct that holds the values for each web resume experience
// in the about page
type Experience struct {
	Company         string   `mapstructure:"company"`
	URL             string   `mapstructure:"url"`
	Role            string   `mapstructure:"role"`
	StartDate       string   `mapstructure:"start_date"`
	EndDate         string   `mapstructure:"end_date"`
	Accomplishments []string `mapstructure:"accomplishments"`
}

// WebResume is the struct that holds the all the values from the config file
// used to render the web resume
type WebResume struct {
	Title             string `mapstructure:"title"`
	FullName          string `mapstructure:"full_name"`
	GithubProfile     string `mapstructure:"github_profile"`
	HomePageSubtitle  string `mapstructure:"home_page_subtitle"`
	AboutPageSubtitle string `mapstructure:"about_page_subtitle"`
	Experience        []Experience
	Skills            []string `mapstructure:"skills"`
	Certifications    []string `mapstructure:"certifications"`
}

func main() {
	// Set log level
	setLogLevel()

	if getEnvVar("LOG_FORMAT", "text") == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	// Load config
	fileConfig := config.FileConfig{
		Path: getEnvVar("CONFIG_PATH", "./config"),
		Name: "config",
		Type: "yaml",
	}

	// Initialize config struct
	var config WebResume

	// Load config into struct
	fileConfig.Load(&config)

	// Initialize Gin
	gin.SetMode(getEnvVar("GIN_MODE", "release")) // Set gin mode to release by default, can be overridden by GIN_MODE env var (e.g. GIN_MODE=debug)
	r := gin.New()                                // empty engine
	if getEnvVar("LOG_FORMAT", "text") == "json" {
		r.Use(middleware.DefaultStructuredLogger()) // adds our new middleware
	} else {
		r.Use(ginlogrus.Logger(log.New()))
	}
	r.Use(gin.Recovery()) // adds gin's default recovery middleware

	// Disable trusting proxy
	r.SetTrustedProxies(nil)

	// Load templates
	r.HTMLRender = loadTemplates()

	// Serve static files, use STATIC_PATH env var if set, otherwise use ./static
	r.Static("/static", getEnvVar("STATIC_PATH", "./static"))

	// Setup health route, useful for load balancer health checks
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Setup index route
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"full_name":          config.FullName,
			"github_profile":     config.GithubProfile,
			"header_class":       "mb-auto", // The home page centers the main content
			"home_active":        "active",
			"home_page_subtitle": config.HomePageSubtitle,
			"title":              config.Title,
			"year":               time.Now().Year(),
		})
	})

	// Setup about route
	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.html", gin.H{
			"about_active":        "active",
			"about_page_subtitle": config.AboutPageSubtitle,
			"certifications":      config.Certifications,
			"experience":          config.Experience,
			"full_name":           config.FullName,
			"github_profile":      config.GithubProfile,
			"header_class":        "mb-5", // The about page does not center the main content
			"skills":              config.Skills,
			"title":               config.Title,
			"year":                time.Now().Year(),
		})
	})

	// Run the server
	log.Info("Starting server on port " + getEnvVar("PORT", "8080"))
	if err := r.Run(":" + getEnvVar("PORT", "8080")); err != nil {
		log.Error(err)
	}
}

// setLogLevel sets the log level based on the LOG_LEVEL env var
func setLogLevel() {
	switch getEnvVar("LOG_LEVEL", "info") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func getEnvVar(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		v = d
	}

	return v
}

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// Get all files in layouts/ directory
	layouts, err := filepath.Glob(getEnvVar("TEMPLATES_PATH", "./templates") + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Get all files in includes/ directory
	includes, err := filepath.Glob(getEnvVar("TEMPLATES_PATH", "./templates") + "/includes/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
