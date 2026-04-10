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
			return echo.NewHTTPError(400, "Conversion eror, id must be integer like")
		}

		id_uint := uint(id_conv)

		var category models.Category

		result := db.Where("CategoryID = ?", id_uint).First(&category)

		if result.Error != nil {
			c.Logger().Error(result.Error) 
			if errors.Is(result.Error, gorm.ErrRecordNotFound){
				return echo.NewHTTPError(404, "Category not found")
			}

			return echo.NewHTTPError(500 ,"Interal db error")
		}


		return c.JSON(http.StatusOK, category)
	}
}

func GetCategories(db *gorm.DB) echo.HandlerFunc{
		return func (c echo.Context) error {
			var categories []models.Category
			db.Find(&categories)
			return c.JSON(http.StatusOK, categories)
		}
}

func CreateCategory(db *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error{
			name := c.QueryParam("name")

			new_cat := new(models.Category)
			new_cat.CategoryName=name

			if name==""{
				return echo.NewHTTPError(400, "Provide name in order to create category. It cannot be empty string")
			}

			result := db.Create(new_cat)

			if result.Error != nil{
				return echo.NewHTTPError(500, "Something went wrong. Not added the instance to database")
			}

			return c.JSON(http.StatusOK, new_cat)
		}
}

func UpdateCategory(db *gorm.DB)echo.HandlerFunc{
	return func (c echo.Context) error{
		id := c.Param("id")
		id_conv, err_conv := strconv.Atoi(id)
		if err_conv != nil{
			return echo.NewHTTPError(400, "Conversion eror, id must be integer like")
		}

		id_uint := uint(id_conv)

		name := c.QueryParam("name")

		result := db.Where("CategoryID = ?", id_uint).Updates(models.Category{
			CategoryName: name,
		})

		if result.Error != nil{
			return echo.NewHTTPError(500, "Couldn't update the specified category")
		}

		if result.RowsAffected == 0{
			return echo.NewHTTPError(404, "Couldnt find the category")
		}

		return c.String(200, fmt.Sprintf("Updated category %d", id_uint))
	}
}

func DeleteCategory(db *gorm.DB)echo.HandlerFunc{
	return func (c echo.Context) error{
		id := c.Param("id")
		id_conv, err_conv := strconv.Atoi(id)
		if err_conv != nil{
			return echo.NewHTTPError(400, "Conversion eror, id must be integer like")
		}
		id_uint := uint(id_conv)
		
		result := db.Where("CategoryID = ?", id_uint).Delete(models.Category{})

		if result.Error != nil{
			return echo.NewHTTPError(500, "Couldn't delete the specified category")
		}

		if result.RowsAffected == 0{
			return echo.NewHTTPError(404, "Couldnt find the category")
		}

		return c.String(200, fmt.Sprintf("Deleted category %d", id_uint))	
	}
}