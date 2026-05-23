-- Full reset
PRAGMA foreign_keys = OFF;
DROP TABLE IF EXISTS Basket;
DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Products;
DROP TABLE IF EXISTS Category;
DROP TABLE IF EXISTS Payments;
PRAGMA foreign_keys = ON;

-- 1. Categories
CREATE TABLE Category (
    CategoryID INTEGER PRIMARY KEY AUTOINCREMENT, 
    CategoryName TEXT UNIQUE
);

-- 2. Products
CREATE TABLE Products (
    ProductID INTEGER PRIMARY KEY AUTOINCREMENT,
    ProductName TEXT UNIQUE, 
    Price REAL,
    CategoryID INTEGER, 
    FOREIGN KEY (CategoryID) REFERENCES Category(CategoryID) ON DELETE SET NULL ON UPDATE CASCADE
);

-- 3. Users
CREATE TABLE Users (
    UserID INTEGER PRIMARY KEY AUTOINCREMENT, 
    UserName TEXT UNIQUE,
    Email TEXT UNIQUE,
    Password TEXT,    
    Provider TEXT DEFAULT 'local',
    ProviderID TEXT,
    ProviderToken TEXT
);

-- 4. Basket
CREATE TABLE Basket (
    ID INTEGER PRIMARY KEY AUTOINCREMENT, 
    UserID INTEGER, 
    ProductID INTEGER, 
    Quantity INTEGER DEFAULT 1,
    FOREIGN KEY(UserID) REFERENCES Users(UserID) ON DELETE CASCADE ON UPDATE CASCADE, 
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX idx_user_product on Basket(UserID, ProductID);

-- 5. Payments
CREATE TABLE Payments (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID INTEGER,
    TotalAmount REAL,
    Status TEXT,
    CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(UserID) REFERENCES Users(UserID) ON DELETE SET NULL ON UPDATE CASCADE
);

-- ===== DATA =====

-- Categories (furniture)
INSERT INTO Category (CategoryName) VALUES 
('Sofas & Couches'),
('Tables & Desks'),
('Chairs & Armchairs'),
('Beds'),
('Wardrobes & Dressers'),
('Outdoor Furniture'),
('Shelves & Bookcases'),
('Kitchen Furniture'),
('Office Furniture'),
('Decor & Accessories');

-- Products
INSERT INTO Products (ProductName, Price, CategoryID) VALUES 
-- Sofas
('3-Seater Sofa Oslo', 2499.99, 1),
('Corner Sofa Milano', 3999.00, 1),
('Sofa Bed Luna', 1899.50, 1),

-- Tables
('Oak Dining Table Classic', 1299.00, 2),
('Gaming Desk RGB', 899.99, 2),
('Extendable Table Family', 1599.00, 2),

-- Chairs
('Wooden Chair Nordic', 199.99, 3),
('Office Chair ErgoPro', 799.00, 3),
('Velvet Armchair Comfort', 999.00, 3),

-- Beds
('King Size Bed Dream', 2999.00, 4),
('Storage Bed SmartSleep', 1899.00, 4),
('Kids Bed Bunny', 799.00, 4),

-- Wardrobes
('Sliding Wardrobe Modern', 2199.00, 5),
('6-Drawer Dresser Simple', 899.00, 5),
('Nightstand Mini', 199.00, 5),

-- Outdoor
('Rattan Garden Set', 1499.00, 6),
('Wooden Sun Lounger Relax', 399.00, 6),
('Patio Table Outdoor', 699.00, 6),

-- Shelves
('Industrial Shelf Loft', 599.00, 7),
('Wall Shelf Cube', 149.00, 7),
('Classic Bookcase', 899.00, 7),

-- Kitchen
('Kitchen Cabinet Basic', 499.00, 8),
('Kitchen Island Premium', 1799.00, 8),
('Compact Kitchen Table', 699.00, 8),

-- Office
('Corner Desk OfficeMax', 1199.00, 9),
('Mobile Drawer Unit', 299.00, 9),
('Conference Chair Simple', 249.00, 9),

-- Decor
('Wall Mirror Loft', 199.00, 10),
('Floor Lamp ModernLight', 349.00, 10),
('Scandinavian Rug Soft', 499.00, 10);

-- Users
INSERT INTO Users (UserName, Email, Password, Provider) VALUES 
('Wrex', 'wrex@example.com', '', 'local'),
('Jakub', 'jakub@test.pl', '', 'local'),
('Admin', 'admin@store.com', '', 'local');