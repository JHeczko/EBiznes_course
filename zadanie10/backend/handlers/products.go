package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"zadanie4/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const errProductIDMsg = "Product ID must be an integer"

func parseProductIDParam(c echo.Context) (uint, error) {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Errorf("ID conversion error: %v", err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, errProductIDMsg)
	}
	return uint(idParam), nil
}

func GetProducts(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var products []models.Product

		query := db.Model(&models.Product{}).Preload("Category")

		if c.QueryParam("cheap") == "true" {
			query = query.Scopes(models.FilterCheapProduct(200))
		}

		if err := query.Find(&products).Error; err != nil {
			c.Logger().Errorf("Database error fetching products: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve products list")
		}
		return c.JSON(http.StatusOK, products)
	}
}

func GetProduct(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := parseProductIDParam(c)
		if err != nil {
			return err
		}

		var product models.Product
		result := db.Preload("Category").Scopes(models.ByProductID(id)).First(&product)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Product with ID %d not found", id))
			}
			c.Logger().Errorf("Database error fetching product %d: %v", id, result.Error)
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

		if err := db.Create(p).Error; err != nil {
			c.Logger().Errorf("Database error creating product: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not save product to database")
		}

		if err := db.Preload("Category").First(p, p.ProductID).Error; err != nil {
			c.Logger().Warnf("Product %d created, but failed to preload category: %v", p.ProductID, err)
		}

		return c.JSON(http.StatusCreated, p)
	}
}

func UpdateProduct(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := parseProductIDParam(c)
		if err != nil {
			return err
		}

		newName := c.QueryParam("name")
		priceStr := c.QueryParam("price")
		catIDParam := c.QueryParam("cat")

		priceFloat, err := strconv.ParseFloat(priceStr, 32)
		if err != nil {
			c.Logger().Errorf("Price conversion error for %s: %v", priceStr, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Price must be a valid decimal number")
		}
		newPrice := float32(priceFloat)

		catInt, err := strconv.Atoi(catIDParam)
		if err != nil {
			c.Logger().Errorf("Category ID conversion error for %s: %v", catIDParam, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Category ID must be an integer")
		}
		newCatID := uint(catInt)

		result := db.Model(&models.Product{}).Scopes(models.ByProductID(id)).Updates(map[string]interface{}{
			"ProductName": newName,
			"Price":       newPrice,
			"CategoryID":  newCatID,
		})

		if result.Error != nil {
			c.Logger().Errorf("Database error updating product %d: %v", id, result.Error)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update product in database")
		}

		if result.RowsAffected == 0 {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No product found with ID %d", id))
		}

		return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Updated product with id %d", id)})
	}
}

func DeleteProduct(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		prodID, err := parseProductIDParam(c)
		if err != nil {
			return err
		}

		result := db.Delete(&models.Product{}, prodID)

		if result.Error != nil {
			c.Logger().Errorf("Database error deleting product %d: %v", prodID, result.Error)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete product from database")
		}

		if result.RowsAffected == 0 {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Product with ID %d not found", prodID))
		}

		return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Deleted product with id %d", prodID)})
	}
}