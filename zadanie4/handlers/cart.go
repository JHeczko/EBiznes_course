package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"zadanie4/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)


func GetItems(db *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error{
		userId_conv, err_conv := strconv.Atoi(c.Param("user_id"))
		if err_conv != nil{
			return echo.NewHTTPError(400, "User ID must be integer")
		}
		userID := uint(userId_conv)


		var basket []models.Basket
		
		result := db.Preload("Product.Category").Where("UserID = ?", userID).Find(&basket)
		
		if result.Error != nil{
			return echo.NewHTTPError(500, "Database error")
		}

		return c.JSON(200, basket)
	}
}

func CreateItem(db *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error{
		// ==== USER ID
		userId_conv, err_conv := strconv.Atoi(c.Param("user_id"))
		if err_conv != nil{
			return echo.NewHTTPError(400, "User ID must be integer")
		}
		userID := uint(userId_conv)

		// ==== PROD ID
		prodID_conv, err_conv := strconv.Atoi(c.QueryParam("prod_id"))
		if err_conv != nil{
			return echo.NewHTTPError(400, "Product ID must be integer")
		}
		prodID := uint(prodID_conv)
		
		// ==== QUANTITY
		quantity := uint(1) 

		qParam := c.QueryParam("quantity")
		if qParam != "" {
			q_conv, err := strconv.Atoi(qParam)
			if err != nil {
				return echo.NewHTTPError(400, "Quantity must be integer")
			}
			quantity = uint(q_conv)
		}

		var item models.Basket

		err := db.Transaction(func (tx *gorm.DB) error {
			err := tx.Where("UserID = ? AND ProductID = ?",userID, prodID).First(&item).Error
			// we got match
			if err == nil{
				item.Quantity += quantity
				return tx.Save(&item).Error
			// we doesnt have item inside DB
			}else{
				if errors.Is(err, gorm.ErrRecordNotFound){
					item = models.Basket{
						UserID: userID,
						ProductID: prodID,
						Quantity: quantity,
					}

					return tx.Create(&item).Error
				}else{
					return echo.NewHTTPError(500, "Database error")
				}
			}
		})
	
		if err != nil{
			return echo.NewHTTPError(500, "Internal database error")
		}

        return c.JSON(200, item)
	}
}

func UpdateItem(db *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error{
		userId_conv, err_conv := strconv.Atoi(c.Param("user_id"))
		if err_conv != nil{
			return echo.NewHTTPError(400, "User ID must be integer")
		}
		userID := uint(userId_conv)

		// ==== PROD ID
		prodID_conv, err_conv := strconv.Atoi(c.QueryParam("prod_id"))
		if err_conv != nil{
			return echo.NewHTTPError(400, "Product ID must be integer")
		}
		prodID := uint(prodID_conv)

		// ==== QUANTITY
		q_conv, err_conv := strconv.Atoi(c.QueryParam("quantity"))
		if err_conv != nil{
			return echo.NewHTTPError(400, "You must provide quantity (as integer)")
		}
		quantity := uint(q_conv)

		var item models.Basket

		err := db.Transaction(func(tx *gorm.DB) error {
			result := tx.Where("UserID = ? AND ProductID = ?", userID, prodID).First(&item)
			if result.Error != nil{
				return result.Error
			}

			if quantity == 0{
				return tx.Delete(&item).Error
			}

			item.Quantity = quantity
			return tx.Save(&item).Error
		})

	
		// error checkinh
		if err != nil{
			if errors.Is(err, gorm.ErrRecordNotFound){
				return echo.NewHTTPError(404, "Cart item not found")	
			}

			return echo.NewHTTPError(500, "Database error while updating the cart quantity")
		}

		return c.String(200, fmt.Sprintf("Updated cart cart for %d and product id %d for quantity %d", userID, prodID, quantity))
	}
}

func DeleteItem(db *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error{
		// ==== USER ID
		userId_conv, err_conv := strconv.Atoi(c.Param("user_id"))
		if err_conv != nil{
			return echo.NewHTTPError(400, "User ID must be integer")
		}
		userID := uint(userId_conv)

		// ==== PROD ID
		prodID_conv, err_conv := strconv.Atoi(c.QueryParam("prod_id"))
		if err_conv != nil{
			return echo.NewHTTPError(400, "Product ID must be integer")
		}
		prodID := uint(prodID_conv)

		// Usuwa wiersz, który pasuje do obu ID-ków
		result := db.Where("UserID = ? AND ProductID = ?", userID, prodID).Delete(&models.Basket{})

		if result.Error != nil {
			return echo.NewHTTPError(500, "Błąd podczas usuwania")
		}

		// Opcjonalnie sprawdzasz, czy w ogóle coś zostało usunięte
		if result.RowsAffected == 0 {
			return echo.NewHTTPError(404, "Nie znaleziono takiego produktu w koszyku")
		}

		return c.String(200, fmt.Sprintf("Deleted for user %d, item %d ", userID, prodID))
	}
}