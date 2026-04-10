package models

// ====== CATEGORY ======
type Category struct{
	CategoryID uint `gorm:"primaryKey;column:CategoryID" json:"id"`
	CategoryName string `gorm:"column:CategoryName" json:"name"`

	// Relations
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

func (Category) TableName() string {
	return  "Category"
}

// ====== PRODUCTS ======
type Product struct{
	ProductID uint `gorm:"primaryKey;column:ProductID" json:"id"`
	ProductName string `gorm:"column:ProductName" json:"name"`
	Price float32 `gorm:"column:Price" json:"price"`
	CategoryID *uint `gorm:"column:CategoryID" json:"cat_id"`

	// Relations
	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category,omitempty"`
}

func (Product) TableName() string{
	return "Products"
}

// ====== USERS ======
type Users struct{
	UserID uint `gorm:"primaryKey;column:UserID" json:"user_id"`
	UserName string `gorm:"column:UserName" json:"user_name"`
	Email string `gorm:"column:Email" json:"mail"`

	// Relations
	Basket []Basket `gorm:"foreignKey:UserID;references:UserID" json:"baskets,omitempty"` 
}

func (Users) TableName() string{
	return "Users"
}

// ====== BASKET ======
type Basket struct{
	ID uint `gorm:"primaryKey;column:ID" json:"id"` 
    UserID uint `gorm:"column:UserID;index:idx_user_id" json:"user_id"`
    ProductID uint `gorm:"column:ProductID" json:"prod_id"`
    Quantity uint `gorm:"column:Quantity;default:1" json:"quantity"` 

	// Relations
	User      *Users    `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
    Product   *Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (Basket) TableName() string{
	return "Basket"
}
