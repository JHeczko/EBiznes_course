package models

// ====== CATEGORY ======
type Category struct{
	CategoryID uint `gorm:"primaryKey;column:CategoryID" json:"id"`
	CategoryName string `gorm:"column:CategoryName" json:"name"`

	// Relations
	ProductsPage []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
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
type Users struct {
    UserID    uint   `gorm:"primaryKey;column:UserID" json:"user_id"`
    UserName  string `gorm:"column:UserName" json:"user_name"`
    Email     string `gorm:"column:Email" json:"mail"`
    
    Password      string `gorm:"column:Password" json:"-"` // "-" bo nie chcemy go wysyłać do frontu
    Provider      string `gorm:"column:Provider;default:'local'" json:"provider"`
    ProviderID    string `gorm:"column:ProviderID" json:"provider_id"`
    ProviderToken string `gorm:"column:ProviderToken" json:"-"` // Tutaj trzymasz token od Google/GitHub

    // Relations
    Basket []Basket `gorm:"foreignKey:UserID;references:UserID" json:"baskets,omitempty"` 
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

type Payments struct {
	ID          uint      `gorm:"primaryKey;column:ID" json:"id"`
	UserID      uint      `gorm:"column:UserID" json:"user_id"`
	TotalAmount float64   `gorm:"column:TotalAmount" json:"total_amount"`
	Status      string    `gorm:"column:Status" json:"status"`
	CreatedAt   string `gorm:"column:CreatedAt" json:"created_at"`

	// relations
	User *Users `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
}

func (Payments) TableName() string{
	return "Payments"
}