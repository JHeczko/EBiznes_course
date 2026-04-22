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

func GetCategory(db *gorm.DB) echo.HandlerFunc{
    return func (c echo.Context) error{
        id := c.Param("id")
        id_conv, err_conv := strconv.Atoi(id)
        if err_conv != nil{
            c.Logger().Errorf("Category ID conversion error: %v", err_conv)
            return echo.NewHTTPError(400, "Invalid category ID format. ID must be an integer")
        }

        id_uint := uint(id_conv)
        var category models.Category

        result := db.Scopes(models.ByCategoryID(id_uint)).First(&category)

        if result.Error != nil {
            if errors.Is(result.Error, gorm.ErrRecordNotFound){
                return echo.NewHTTPError(404, fmt.Sprintf("Category with ID %d not found", id_uint))
            }
            c.Logger().Errorf("Database error fetching category %d: %v", id_uint, result.Error)
            return echo.NewHTTPError(500 ,"Internal database error while retrieving category")
        }

        return c.JSON(http.StatusOK, category)
    }
}

func GetCategories(db *gorm.DB) echo.HandlerFunc{
        return func (c echo.Context) error {
            var categories []models.Category

            if err := db.Find(&categories).Error; err != nil {
                c.Logger().Errorf("Database error fetching categories: %v", err)
                return echo.NewHTTPError(500, "Internal database error while fetching categories list")
            }
            
			return c.JSON(http.StatusOK, categories)
        }
}

func CreateCategory(db *gorm.DB) echo.HandlerFunc{
    return func (c echo.Context) error{
            name := c.QueryParam("name")

            if name == "" {
                return echo.NewHTTPError(400, "Category name is required and cannot be empty")
            }

            new_cat := &models.Category{CategoryName: name}
            result := db.Create(new_cat)

            if result.Error != nil{
                c.Logger().Errorf("Database error creating category: %v", result.Error)
                return echo.NewHTTPError(500, "Failed to create new category in database")
            }

            return c.JSON(http.StatusOK, new_cat)
        }
}

func UpdateCategory(db *gorm.DB)echo.HandlerFunc{
    return func (c echo.Context) error{
        id := c.Param("id")
        id_conv, err_conv := strconv.Atoi(id)
        if err_conv != nil{
            c.Logger().Errorf("Category ID conversion error: %v", err_conv)
            return echo.NewHTTPError(400, "Invalid category ID format for update")
        }

        id_uint := uint(id_conv)
        name := c.QueryParam("name")

        result := db.Scopes(models.ByCategoryID(id_uint)).Updates(models.Category{
            CategoryName: name,
        })

        if result.Error != nil{
            c.Logger().Errorf("Database error updating category %d: %v", id_uint, result.Error)
            return echo.NewHTTPError(500, "Failed to update the specified category")
        }

        if result.RowsAffected == 0{
            return echo.NewHTTPError(404, fmt.Sprintf("Cannot update: Category with ID %d does not exist", id_uint))
        }

        return c.String(200, fmt.Sprintf("Successfully updated category %d", id_uint))
    }
}

func DeleteCategory(db *gorm.DB)echo.HandlerFunc{
    return func (c echo.Context) error{
        id := c.Param("id")
        id_conv, err_conv := strconv.Atoi(id)
        if err_conv != nil{
            c.Logger().Errorf("Category ID conversion error: %v", err_conv)
            return echo.NewHTTPError(400, "Invalid category ID format for deletion")
        }
        id_uint := uint(id_conv)
        
        result := db.Scopes(models.ByCategoryID(id_uint)).Delete(models.Category{})

        if result.Error != nil{
            c.Logger().Errorf("Database error deleting category %d: %v", id_uint, result.Error)
            return echo.NewHTTPError(500, "Failed to delete the specified category")
        }

        if result.RowsAffected == 0{
            return echo.NewHTTPError(404, fmt.Sprintf("Cannot delete: Category with ID %d does not exist", id_uint))
        }

        return c.String(200, fmt.Sprintf("Successfully deleted category %d", id_uint))   
    }
}