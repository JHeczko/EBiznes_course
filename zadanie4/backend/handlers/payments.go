package handlers

import (
	"strconv"
	"zadanie4/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetPayments(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userIDConv, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.Logger().Errorf(errUserIDConv, err)
			return echo.NewHTTPError(400, "Invalid User ID format")
		}
		userID := uint(userIDConv)

		var payments []models.Payments

		if err := db.Table("Payments").Scopes(models.ByUserID(userID)).Find(&payments).Error; err != nil {
			c.Logger().Errorf("Couldn't get payments: %v", err)
			return echo.NewHTTPError(500, "Couldn't get payments")
		}

		return c.JSON(200, payments)
	}
}

func AddPayment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userIDConv, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.Logger().Errorf(errUserIDConv, err)
			return echo.NewHTTPError(400, "Invalid User ID format")
		}
		userID := uint(userIDConv)

		var total float64
		if err := db.Scopes(models.SumBasketForUserID(userID)).Scan(&total).Error; err != nil {
			c.Logger().Errorf("Basket sum failed: %v", err)
			return echo.NewHTTPError(500, "Cannot calculate basket total")
		}

		if total == 0 {
			return echo.NewHTTPError(500, "Cannot pay for empty basket. Total == 0")
		}

		payment := models.Payments{
			UserID:      userID,
			TotalAmount: total,
			Status:      "paid",
		}

		if err := db.Create(&payment).Error; err != nil {
			c.Logger().Errorf("Payment create failed: %v", err)
			return echo.NewHTTPError(500, "Cannot create payment")
		}

		return c.JSON(200, payment)
	}
}