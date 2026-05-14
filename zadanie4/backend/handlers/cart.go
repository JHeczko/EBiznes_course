package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"zadanie4/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	errUserIDConv    = "UserID conversion failed: %v"
	errProductIDConv = "ProductID conversion failed: %v"
)

func parseUserID(c echo.Context) (uint, error) {
	conv, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.Logger().Errorf(errUserIDConv, err)
		return 0, echo.NewHTTPError(400, "Invalid User ID format. It must be an integer")
	}
	return uint(conv), nil
}

func parseProductID(c echo.Context) (uint, error) {
	conv, err := strconv.Atoi(c.QueryParam("prod_id"))
	if err != nil {
		c.Logger().Errorf(errProductIDConv, err)
		return 0, echo.NewHTTPError(400, "Product ID must be a valid integer")
	}
	return uint(conv), nil
}

func GetItems(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := parseUserID(c)
		if err != nil {
			return err
		}

		var basket []models.Basket
		query := db.Model(&models.Basket{}).Preload("Product.Category")

		if c.QueryParam("cheap") == "true" {
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

func createItemTransaction(db *gorm.DB, userID, prodID, quantity uint) (models.Basket, error) {
	var item models.Basket
	err := db.Transaction(func(tx *gorm.DB) error {
		out := tx.Scopes(models.ByUserID(userID)).
			Scopes(models.ByProductID(prodID)).
			Take(&item)

		if out.Error != nil && !errors.Is(out.Error, gorm.ErrRecordNotFound) {
			return out.Error
		}

		if out.RowsAffected > 0 {
			item.Quantity += quantity
			return tx.Save(&item).Error
		}

		item = models.Basket{
			UserID:    userID,
			ProductID: prodID,
			Quantity:  quantity,
		}
		return tx.Create(&item).Error
	})
	return item, err
}

func CreateItem(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := parseUserID(c)
		if err != nil {
			return err
		}

		prodID, err := parseProductID(c)
		if err != nil {
			return err
		}

		quantity := uint(1)
		if qParam := c.QueryParam("quantity"); qParam != "" {
			qConv, qErr := strconv.Atoi(qParam)
			if qErr != nil {
				c.Logger().Errorf("Quantity conversion failed: %v", qErr)
				return echo.NewHTTPError(400, "Quantity must be a valid integer")
			}
			quantity = uint(qConv)
		}

		item, err := createItemTransaction(db, userID, prodID, quantity)
		if err != nil {
			c.Logger().Errorf("Transaction failed for User %d, Product %d: %v", userID, prodID, err)
			return echo.NewHTTPError(500, "Internal server error while processing cart item")
		}

		return c.JSON(200, item)
	}
}

func updateItemTransaction(db *gorm.DB, userID, prodID, quantity uint) (models.Basket, error) {
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
	return item, err
}

func UpdateItem(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := parseUserID(c)
		if err != nil {
			return err
		}

		prodID, err := parseProductID(c)
		if err != nil {
			return err
		}

		qConv, err := strconv.Atoi(c.QueryParam("quantity"))
		if err != nil {
			c.Logger().Errorf("Quantity conversion failed: %v", err)
			return echo.NewHTTPError(400, "Quantity must be an integer")
		}
		quantity := uint(qConv)

		item, err := updateItemTransaction(db, userID, prodID, quantity)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return echo.NewHTTPError(404, fmt.Sprintf("Product %d not found in user %d's cart", prodID, userID))
			}
			c.Logger().Errorf("Update transaction failed for User %d: %v", userID, err)
			return echo.NewHTTPError(500, "Database error while updating cart quantity")
		}

		return c.JSON(200, item)
	}
}

func deleteAllItems(db *gorm.DB, c echo.Context, userID uint) error {
	result := db.Scopes(models.ByUserID(userID)).Delete(&models.Basket{})
	if result.Error != nil {
		c.Logger().Errorf("Delete ALL failed for User %d: %v", userID, result.Error)
		return echo.NewHTTPError(500, "Internal database error during deletion")
	}

	if result.RowsAffected == 0 {
		return c.String(204, fmt.Sprintf("No items for user %d. Deleted nothing", userID))
	}

	return c.JSON(200, map[string]interface{}{
		"user_id": userID,
		"deleted": result.RowsAffected,
		"type":    "all",
	})
}

func deleteSingleItem(db *gorm.DB, c echo.Context, userID uint, prodParam string) error {
	prodIDConv, err := strconv.Atoi(prodParam)
	if err != nil {
		c.Logger().Errorf(errProductIDConv, err)
		return echo.NewHTTPError(400, "Invalid Product ID")
	}
	prodID := uint(prodIDConv)

	result := db.Scopes(models.ByUserID(userID)).Scopes(models.ByProductID(prodID)).Delete(&models.Basket{})
	if result.Error != nil {
		c.Logger().Errorf("Delete failed for User %d, Product %d: %v", userID, prodID, result.Error)
		return echo.NewHTTPError(500, "Internal database error during deletion")
	}

	if result.RowsAffected == 0 {
		return echo.NewHTTPError(404, "Item not found in cart; nothing to delete")
	}

	return c.JSON(200, map[string]interface{}{
		"user_id": userID,
		"prod_id": prodID,
		"deleted": result.RowsAffected,
		"type":    "single",
	})
}

func DeleteItem(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := parseUserID(c)
		if err != nil {
			return err
		}

		prodParam := c.QueryParam("prod_id")
		if prodParam == "" {
			return deleteAllItems(db, c, userID)
		}
		return deleteSingleItem(db, c, userID, prodParam)
	}
}