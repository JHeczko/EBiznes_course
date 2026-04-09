package models

type Category struct{
	CategoryID uint `gorm:"primaryKey;column:CategoryID" json:"id"`
	CategoryName string `gorm:"column:CategoryName" json:"name"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
func (Product) TableName() string{
	return "Products"
}
	
type Product struct{
	ProductID uint `gorm:"primaryKey;column:ProductID" json:"id"`
	ProductName string `gorm:"column:ProductName" json:"name"`
	Price float32 `gorm:"column:Price" json:"price"`
	CategoryID *uint `gorm:"column:CategoryID" json:"cat_id"`
	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category,omitempty"`
}
func (Category) TableName() string {
	return  "Category"
}

type Basket struct{

}