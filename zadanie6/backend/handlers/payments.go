package handlers

import (
	"strconv"
	"zadanie4/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetPayments(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID_conv, err_conv := strconv.Atoi(c.Param("user_id"))
		if err_conv != nil {
			c.Logger().Errorf("UserID conversion failed: %v", err_conv)
			return echo.NewHTTPError(400, "Invalid User ID format")
		}
		userID := uint(userID_conv)

		var payments []models.Payments

		err := db.Table("Payments").Scopes(models.ByUserID(userID)).Find(&payments).Error

		if err != nil {
			c.Logger().Errorf("Couldn't get payments: %v", err)
			return echo.NewHTTPError(500, "Couldn't get payments")
		}

		return c.JSON(200, payments)
	}
}

func AddPayment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		userID_conv, err_conv := strconv.Atoi(c.Param("user_id"))
		if err_conv != nil {
			c.Logger().Errorf("UserID conversion failed: %v", err_conv)
			return echo.NewHTTPError(400, "Invalid User ID format")
		}
		userID := uint(userID_conv)

		var total float64
		err := db.Scopes(models.SumBasketForUserID(userID)).Scan(&total).Error
		if err != nil {
			c.Logger().Errorf("Basket sum failed: %v", err)
			return echo.NewHTTPError(500, "Cannot calculate basket total")
		}

		if total==0{
			c.Logger().Errorf("Cannot pay for empty basket: %v", err)
			return echo.NewHTTPError(500, "Cannot pay for empty basket. Total == 0")
		}

		payment := models.Payments{
			UserID:      userID,
			TotalAmount: total,
			Status:      "paid",
		}

		err = db.Create(&payment).Error
		if err != nil {
			c.Logger().Errorf("Payment create failed: %v", err)
			return echo.NewHTTPError(500, "Cannot create payment")
		}

		// 3. response
		return c.JSON(200, payment)
	}
}