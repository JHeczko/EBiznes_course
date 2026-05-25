package main

import (
	"net/http"
	"zadanie4/handlers"
	"zadanie4/middleware_handlers"
	"github.com/labstack/echo/v4/middleware"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log"
	"os"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}


	// databse init
	// db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
	// 	Logger: logger.New(
	// 		log.New(os.Stdout, "\r\n", log.LstdFlags),
	// 		logger.Config{
	// 			LogLevel:                  logger.Error,
	// 			IgnoreRecordNotFoundError: true,
	// 			Colorful:                  true,
	// 		},
	// 	),
	// })
	db, err := gorm.Open(sqlite.Open("/home/appuser/app/database.db"), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel:                  logger.Error,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})


	if err != nil{
		panic("Blad otwarcia bazy danych")
	}

	// echo init
	e := echo.New()
	
	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{
        http.MethodGet,
        http.MethodPost,
        http.MethodPatch,
        http.MethodDelete,
    },
    AllowHeaders: []string{
        echo.HeaderOrigin,
        echo.HeaderContentType,
        echo.HeaderAccept,
        "ngrok-skip-browser-warning",
		echo.HeaderAuthorization,
    },
    AllowCredentials: false,
	}))

	// products endpoints
	productsEndpoints := e.Group("/products")
	{
		productsEndpoints.GET("/:id",handlers.GetProduct(db))
		productsEndpoints.GET("", handlers.GetProducts(db))
		productsEndpoints.POST("", handlers.CreateProduct(db))
		productsEndpoints.PATCH("/:id", handlers.UpdateProduct(db))
		productsEndpoints.DELETE("/:id", handlers.DeleteProduct(db))
	}

	categoryEndpoints := e.Group("/category")
	{
		categoryEndpoints.GET("", handlers.GetCategories(db))
		categoryEndpoints.GET("/:id", handlers.GetCategory(db))
		categoryEndpoints.POST("", handlers.CreateCategory(db))
		categoryEndpoints.PATCH("/:id", handlers.UpdateCategory(db))
		categoryEndpoints.DELETE("/:id", handlers.DeleteCategory(db))
	}

	userEndpoint := "/:user_id"

	cartEndpoints := e.Group("/cart")
	cartEndpoints.Use(middleware_handlers.AuthMiddleware)
	{
		cartEndpoints.GET(userEndpoint, handlers.GetItems(db))
		cartEndpoints.POST(userEndpoint, handlers.CreateItem(db))
		cartEndpoints.PATCH(userEndpoint, handlers.UpdateItem(db))
		cartEndpoints.DELETE(userEndpoint, handlers.DeleteItem(db))
	}

	paymentsEndpoints := e.Group("/payments")
	paymentsEndpoints.Use(middleware_handlers.AuthMiddleware)
	{
		paymentsEndpoints.GET(userEndpoint, handlers.GetPayments(db))
		paymentsEndpoints.POST(userEndpoint, handlers.AddPayment(db))
	}

	authGroup := e.Group("/auth")
	{
		authGroup.POST("/register", handlers.Register(db))
		authGroup.POST("/login", handlers.Login(db))
		authGroup.GET("/google/login", handlers.HandleGoogleLogin)
		authGroup.GET("/google/callback", handlers.HandleGoogleCallback(db))
		authGroup.GET("/github/login", handlers.HandleGithubLogin)
    	authGroup.GET("/github/callback", handlers.HandleGithubCallback(db))
	}


	e.Logger.Fatal(e.Start(":13000"))
}