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

const errCategoryIDConv = "Category ID conversion error: %v"

func parseCategoryID(c echo.Context) (uint, error) {
	idConv, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Errorf(errCategoryIDConv, err)
		return 0, echo.NewHTTPError(400, "Invalid category ID format. ID must be an integer")
	}
	return uint(idConv), nil
}

func GetCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		idUint, err := parseCategoryID(c)
		if err != nil {
			return err
		}

		var category models.Category
		result := db.Scopes(models.ByCategoryID(idUint)).First(&category)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return echo.NewHTTPError(404, fmt.Sprintf("Category with ID %d not found", idUint))
			}
			c.Logger().Errorf("Database error fetching category %d: %v", idUint, result.Error)
			return echo.NewHTTPError(500, "Internal database error while retrieving category")
		}

		return c.JSON(http.StatusOK, category)
	}
}

func GetCategories(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var categories []models.Category

		if err := db.Find(&categories).Error; err != nil {
			c.Logger().Errorf("Database error fetching categories: %v", err)
			return echo.NewHTTPError(500, "Internal database error while fetching categories list")
		}

		return c.JSON(http.StatusOK, categories)
	}
}

func CreateCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.QueryParam("name")

		if name == "" {
			return echo.NewHTTPError(400, "Category name is required and cannot be empty")
		}

		newCat := &models.Category{CategoryName: name}
		result := db.Create(newCat)

		if result.Error != nil {
			c.Logger().Errorf("Database error creating category: %v", result.Error)
			return echo.NewHTTPError(500, "Failed to create new category in database")
		}

		return c.JSON(http.StatusOK, newCat)
	}
}

func UpdateCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		idUint, err := parseCategoryID(c)
		if err != nil {
			return echo.NewHTTPError(400, "Invalid category ID format for update")
		}

		name := c.QueryParam("name")

		result := db.Scopes(models.ByCategoryID(idUint)).Updates(models.Category{
			CategoryName: name,
		})

		if result.Error != nil {
			c.Logger().Errorf("Database error updating category %d: %v", idUint, result.Error)
			return echo.NewHTTPError(500, "Failed to update the specified category")
		}

		if result.RowsAffected == 0 {
			return echo.NewHTTPError(404, fmt.Sprintf("Cannot update: Category with ID %d does not exist", idUint))
		}

		return c.String(200, fmt.Sprintf("Successfully updated category %d", idUint))
	}
}

func DeleteCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		idUint, err := parseCategoryID(c)
		if err != nil {
			return echo.NewHTTPError(400, "Invalid category ID format for deletion")
		}

		result := db.Scopes(models.ByCategoryID(idUint)).Delete(models.Category{})

		if result.Error != nil {
			c.Logger().Errorf("Database error deleting category %d: %v", idUint, result.Error)
			return echo.NewHTTPError(500, "Failed to delete the specified category")
		}

		if result.RowsAffected == 0 {
			return echo.NewHTTPError(404, fmt.Sprintf("Cannot delete: Category with ID %d does not exist", idUint))
		}

		return c.String(200, fmt.Sprintf("Successfully deleted category %d", idUint))
	}
}