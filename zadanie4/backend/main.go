package main

import (
	"net/http"
	"zadanie4/handlers"
	
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log"
	"os"
)


func main() {
	// databse init
	//db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open("/root/app/database.db"), &gorm.Config{
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

	cartEndpoints := e.Group("/cart")
	{
		cartEndpoints.GET("/:user_id", handlers.GetItems(db))
		cartEndpoints.POST("/:user_id", handlers.CreateItem(db))
		cartEndpoints.PATCH("/:user_id", handlers.UpdateItem(db))
		cartEndpoints.DELETE("/:user_id", handlers.DeleteItem(db))
	}

	paymentsEndpoints := e.Group("/payments")

	{
		paymentsEndpoints.GET("/:user_id", handlers.GetPayments(db))
		paymentsEndpoints.POST("/:user_id", handlers.AddPayment(db))
	}

	e.Logger.Fatal(e.Start(":13000"))
}