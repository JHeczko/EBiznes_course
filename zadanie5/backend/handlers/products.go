package handlers

import (
	"fmt"
	"errors"
	"net/http"
	"strconv"
	"zadanie4/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)



func GetProducts(db *gorm.DB) echo.HandlerFunc{
    return func (c echo.Context) error {
        var products []models.Product

		query := db.Model(&models.Product{}).Preload("Category")

		if c.QueryParam("cheap") == "true"{
			query = query.Scopes(models.FilterCheapProduct(200))
		}

        if err := query.Find(&products).Error; err != nil {
            c.Logger().Errorf("Database error fetching products: %v", err)
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve products list")
        }
        return c.JSON(http.StatusOK, products)
    }
}

func GetProduct(db *gorm.DB) echo.HandlerFunc{
    return func (c echo.Context) error{
        idParam, convErr := strconv.Atoi(c.Param("id"))
        if convErr != nil{
            c.Logger().Errorf("ID conversion error: %v", convErr)
            return echo.NewHTTPError(http.StatusBadRequest, "Product ID must be an integer")
        }
        id := uint(idParam)

        var product models.Product

        result := db.Preload("Category").Scopes(models.ByProductID(id)).First(&product)

        if result.Error != nil {
            if errors.Is(result.Error, gorm.ErrRecordNotFound) {
                return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Product with ID %d not found", id))
            }
            c.Logger().Errorf("Database error fetching product %s: %v", id, result.Error) 
            return echo.NewHTTPError(http.StatusInternalServerError, "Error while retrieving product details")
        }

        return c.JSON(http.StatusOK, product)
    }
}

func CreateProduct(db *gorm.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        name := c.QueryParam("name")
        priceStr := c.QueryParam("price")
        categoryStr := c.QueryParam("cat")

        if name == "" {
            return echo.NewHTTPError(http.StatusBadRequest, "Product name is required")
        }

        priceFloat, err := strconv.ParseFloat(priceStr, 32)
        if err != nil {
            c.Logger().Errorf("Price conversion error for %s: %v", priceStr, err)
            return echo.NewHTTPError(http.StatusBadRequest, "Invalid price format. Must be a decimal number")
        }
        price := float32(priceFloat)

        catInt, err := strconv.Atoi(categoryStr)
        if err != nil {
            c.Logger().Errorf("Category ID conversion error for %s: %v", categoryStr, err)
            return echo.NewHTTPError(http.StatusBadRequest, "Category ID must be a valid integer")
        }
        catID := uint(catInt)

        p := &models.Product{
            ProductName: name,
            Price:       price,
            CategoryID:  &catID,
        }

        result := db.Create(p)
        if result.Error != nil {
            c.Logger().Errorf("Database error creating product: %v", result.Error)
            return echo.NewHTTPError(http.StatusInternalServerError, "Could not save product to database")
        }

        if err := db.Preload("Category").First(p, p.ProductID).Error; err != nil {
            c.Logger().Warnf("Product %d created, but failed to preload category: %v", p.ProductID, err)
        }

        return c.JSON(http.StatusCreated, p)
    }
}

func UpdateProduct(db *gorm.DB) echo.HandlerFunc{
    return func (c echo.Context) error{
        idParam, convErr := strconv.Atoi(c.Param("id"))
        if convErr != nil{
            c.Logger().Errorf("ID conversion error: %v", convErr)
            return echo.NewHTTPError(http.StatusBadRequest, "Product ID must be an integer")
        }
        id := uint(idParam)

        newName := c.QueryParam("name")
        price := c.QueryParam("price")
        catIDParam := c.QueryParam("cat")

        priceFloat, convErr := strconv.ParseFloat(price, 32)
        if convErr != nil{
            c.Logger().Errorf("Price conversion error for %s: %v", price, convErr)
            return echo.NewHTTPError(http.StatusBadRequest, "Price must be a valid decimal number")
        }
        newPrice := float32(priceFloat)

        catID, convErr := strconv.Atoi(catIDParam)
        if convErr != nil{
            c.Logger().Errorf("Category ID conversion error for %s: %v", catIDParam, convErr)
            return echo.NewHTTPError(http.StatusBadRequest, "Category ID must be an integer")
        }
        newCatID := uint(catID)

        result := db.Model(&models.Product{}).Scopes(models.ByProductID(id)).Updates(map[string]interface{}{
            "ProductName": newName,
            "Price":       newPrice,
            "CategoryID":  newCatID,
        })

        if result.Error != nil{
            c.Logger().Errorf("Database error updating product %d: %v", id, result.Error)
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update product in database")
        }

        if result.RowsAffected == 0{
            return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No product found with ID %d", id))
        }

        return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Updated product with id %d", id)})
    }
}

func DeleteProduct(db *gorm.DB) echo.HandlerFunc{
    return func (c echo.Context) error{
        id := c.Param("id")

        prod_id, err_conv := strconv.Atoi(id)
        if err_conv != nil{
            c.Logger().Errorf("ID conversion error for deletion: %v", err_conv)
            return echo.NewHTTPError(http.StatusBadRequest, "Product ID must be an integer")
        }

        prod_uid := uint(prod_id)
        result := db.Delete(&models.Product{}, prod_uid)

        if result.Error != nil{
            c.Logger().Errorf("Database error deleting product %d: %v", prod_uid, result.Error)
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete product from database")
        }

        if result.RowsAffected == 0 {
            return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Product with ID %d not found", prod_uid))
        }

        return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Deleted product with id %d", prod_uid)})
    }
}