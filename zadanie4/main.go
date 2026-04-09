package main

import (
	// "net/http"
	// "github.com/labstack/echo/v4"
	"zadanie4/handlers"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func main() {
	// databse init
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	if err != nil{
		panic("Blad otwarcia bazy danych")
	}

	// echo init
	e := echo.New()
	
	// products endpoints
	productsEndpoints := e.Group("/products")
	{
		productsEndpoints.GET("/:id",handlers.GetProduct(db))
		productsEndpoints.GET("", handlers.GetProducts(db))
		productsEndpoints.POST("", handlers.CreateProduct(db))
		productsEndpoints.PATCH("/:id", handlers.UpdateProduct(db))
		productsEndpoints.DELETE("/:id", handlers.DeleteProduct(db))

	}

	e.Logger.Fatal(e.Start(":13000"))

}