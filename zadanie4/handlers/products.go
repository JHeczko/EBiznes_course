package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"zadanie4/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)



func GetProducts(db *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error {
		var products []models.Product
		db.Preload("Category").Find(&products)
		return c.JSON(http.StatusOK, products)
	}
}

func GetProduct(db *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error{
		id := c.Param("id")
		var products []models.Product

		result := db.Preload("Category").Where("ProductID = ?", id).Find(&products)

		if result.Error != nil {
			c.Logger().Error(result.Error) 
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Błąd podczas pobierania danych z bazy",
			})
		}


		return c.JSON(http.StatusOK, products)
	}
}

func CreateProduct(db *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error{
		name := c.QueryParam("name")
		priceStr := c.QueryParam("price")
		categoryStr := c.QueryParam("cat")

		price, _ := strconv.ParseFloat(priceStr,32)

		catInt, _ := strconv.Atoi(categoryStr)
		catID := uint(catInt)

		p := new(models.Product)
		p.ProductName=name
		p.Price=float32(price)
		p.CategoryID=&catID

		result := db.Preload("Category").Create(p)

		db.Preload("Category").Where("ProductID = ?", p.ProductID).First(p)

		if result.Error != nil{
			return echo.NewHTTPError(500, "Something went wrong. Not added the instance to database")
		}

		return c.JSON(http.StatusOK, p)
	}
}

func UpdateProduct(db *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error{
		idParam, convErr := strconv.Atoi(c.Param("id"))
		if convErr != nil{
			return c.String(500, "ID must be integer!")
		}
		id := uint(idParam)

		newName := c.QueryParam("name")

		price := c.QueryParam("price")
		priceFloat, convErr := strconv.ParseFloat(price,32)
		if convErr != nil{
			return c.String(500, "Price must be float!")
		}
		newPrice := float32(priceFloat)

		catID, convErr := strconv.Atoi(c.QueryParam("cat"))
		if convErr != nil{
			return c.String(500, "CatID must be integer!")
		}
		newCatID := uint(catID)

		result := db.Model(&models.Product{}).Where("ProductID = ?", id).Updates(map[string]interface{}{
            "ProductName": newName,
            "Price":       newPrice,
            "CategoryID":  newCatID,
        })

		if result.Error != nil{
			return c.String(500, "DataBase update error")
		}

		if result.RowsAffected == 0{
			return c.String(404, "No product finded with specific id")
		}

		return c.String(200,fmt.Sprintf("Updated product with id %d", id))
	}
}

func DeleteProduct(db *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error{
		id := c.Param("id")

		prod_id, err_conv := strconv.Atoi(id)
		if err_conv != nil{
			return c.String(http.StatusBadRequest, "ID must be integer like")
		}

		prod_uid := uint(prod_id)


		result := db.Delete(&models.Product{}, prod_uid)

		if result.Error != nil{
			return c.String(500, fmt.Sprintf("Delete error, couldnt delete item id %d", prod_uid))
		}

		if result.RowsAffected == 0 {
            return c.String(http.StatusNotFound, fmt.Sprintf("Did not found product of id %d! Not deleting anything", prod_uid))
        }

		return c.String(200, fmt.Sprintf("Deleted %d", prod_uid))	
	}
}