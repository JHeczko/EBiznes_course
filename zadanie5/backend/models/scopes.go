package models

import (
	"gorm.io/gorm"
)

func ByUserID(userID uint) func(db *gorm.DB) *gorm.DB{
	return func(db *gorm.DB) *gorm.DB{
		return db.Where("UserID = ?", userID)
	}
}

func ByProductID(prodID uint) func(db *gorm.DB) *gorm.DB{
	return func(db *gorm.DB) *gorm.DB{
		return db.Where("ProductID = ?", prodID)
	}
}


func ByCategoryID(catID uint) func(db *gorm.DB) *gorm.DB{
	return func(db *gorm.DB) *gorm.DB{
		return db.Where("CategoryID = ?", catID)
	}
}

func FilterCheapProduct(threshold float32) func(db *gorm.DB) *gorm.DB{
	return func(db *gorm.DB) *gorm.DB{
		return db.Where("Price < ?", threshold)
	}
}

func FilterBasketByPrice(threshold float32) func (db *gorm.DB) *gorm.DB{
	return func (db *gorm.DB) *gorm.DB{
		return db.Joins("JOIN ProductsPage ON ProductsPage.ProductID=Basket.ProductID").
		Scopes(FilterCheapProduct(threshold))
	}
}