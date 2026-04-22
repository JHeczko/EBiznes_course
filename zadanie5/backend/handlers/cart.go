package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"zadanie4/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)


func GetItems(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId_conv, err_conv := strconv.Atoi(c.Param("user_id"))
		if err_conv != nil {
			c.Logger().Errorf("UserID conversion failed: %v", err_conv)
			return echo.NewHTTPError(400, "Invalid User ID format. It must be an integer")
		}
		userID := uint(userId_conv)

		var basket []models.Basket

		query := db.Model(&models.Basket{}).Preload("Product.Category")

		if c.QueryParam("cheap") == "true"{
			query = query.Scopes(models.FilterBasketByPrice(300))
		}

		result := query.Scopes(models.ByUserID(userID)).Find(&basket)

		if result.Error != nil {
			c.Logger().Errorf("Database error while fetching basket for user %d: %v", userID, result.Error)
			return echo.NewHTTPError(500, "Failed to retrieve basket items from database")
		}

		return c.JSON(200, basket)
	}
}

func CreateItem(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId_conv, err_conv := strconv.Atoi(c.Param("user_id"))
		if err_conv != nil {
			c.Logger().Errorf("UserID conversion failed: %v", err_conv)
			return echo.NewHTTPError(400, "Invalid User ID format")
		}
		userID := uint(userId_conv)

		prodID_conv, err_conv := strconv.Atoi(c.QueryParam("prod_id"))
		if err_conv != nil {
			c.Logger().Errorf("ProductID conversion failed: %v", err_conv)
			return echo.NewHTTPError(400, "Product ID must be a valid integer")
		}
		prodID := uint(prodID_conv)

		quantity := uint(1)
		qParam := c.QueryParam("quantity")
		if qParam != "" {
			q_conv, err := strconv.Atoi(qParam)
			if err != nil {
				c.Logger().Errorf("Quantity conversion failed: %v", err)
				return echo.NewHTTPError(400, "Quantity must be a valid integer")
			}
			quantity = uint(q_conv)
		}

		var item models.Basket
		err := db.Transaction(func(tx *gorm.DB) error {
			err := tx.Scopes(models.ByUserID(userID)).Scopes(models.ByProductID(prodID)).First(&item).Error
			if err == nil {
				item.Quantity += quantity
				return tx.Save(&item).Error
			} else {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					item = models.Basket{
						UserID:    userID,
						ProductID: prodID,
						Quantity:  quantity,
					}
					return tx.Create(&item).Error
				} else {
					return err
				}
			}
		})

		if err != nil {
			c.Logger().Errorf("Transaction failed for User %d, Product %d: %v", userID, prodID, err)
			return echo.NewHTTPError(500, "Internal server error while processing cart item")
		}

		return c.JSON(200, item)
	}
}

func UpdateItem(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId_conv, err_conv := strconv.Atoi(c.Param("user_id"))
		if err_conv != nil {
			c.Logger().Errorf("UserID conversion failed: %v", err_conv)
			return echo.NewHTTPError(400, "Invalid User ID")
		}
		userID := uint(userId_conv)

		prodID_conv, err_conv := strconv.Atoi(c.QueryParam("prod_id"))
		if err_conv != nil {
			c.Logger().Errorf("ProductID conversion failed: %v", err_conv)
			return echo.NewHTTPError(400, "Invalid Product ID")
		}
		prodID := uint(prodID_conv)

		q_conv, err_conv := strconv.Atoi(c.QueryParam("quantity"))
		if err_conv != nil {
			c.Logger().Errorf("Quantity conversion failed: %v", err_conv)
			return echo.NewHTTPError(400, "Quantity must be an integer")
		}
		quantity := uint(q_conv)

		var item models.Basket
		err := db.Transaction(func(tx *gorm.DB) error {
			result := tx.Scopes(models.ByUserID(userID)).Scopes(models.ByProductID(prodID)).First(&item)
			if result.Error != nil {
				return result.Error
			}

			if quantity == 0 {
				return tx.Delete(&item).Error
			}

			item.Quantity = quantity
			return tx.Save(&item).Error
		})

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return echo.NewHTTPError(404, fmt.Sprintf("Product %d not found in user %d's cart", prodID, userID))
			}
			c.Logger().Errorf("Update transaction failed for User %d: %v", userID, err)
			return echo.NewHTTPError(500, "Database error while updating cart quantity")
		}

		return c.String(200, fmt.Sprintf("Updated cart for user %d: product %d now has quantity %d", userID, prodID, quantity))
	}
}

func DeleteItem(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId_conv, err_conv := strconv.Atoi(c.Param("user_id"))
		if err_conv != nil {
			c.Logger().Errorf("UserID conversion failed: %v", err_conv)
			return echo.NewHTTPError(400, "Invalid User ID")
		}
		userID := uint(userId_conv)

		prodID_conv, err_conv := strconv.Atoi(c.QueryParam("prod_id"))
		if err_conv != nil {
			c.Logger().Errorf("ProductID conversion failed: %v", err_conv)
			return echo.NewHTTPError(400, "Invalid Product ID")
		}
		prodID := uint(prodID_conv)

		result := db.Scopes(models.ByUserID(userID)).Scopes(models.ByProductID(prodID)).Delete(&models.Basket{})

		if result.Error != nil {
			c.Logger().Errorf("Delete failed for User %d, Product %d: %v", userID, prodID, result.Error)
			return echo.NewHTTPError(500, "Internal database error during deletion")
		}

		if result.RowsAffected == 0 {
			return echo.NewHTTPError(404, "Item not found in cart; nothing to delete")
		}

		return c.String(200, fmt.Sprintf("Successfully removed product %d from user %d's cart", prodID, userID))
	}
}