-- Kompletny cykl testu
PRAGMA foreign_keys = OFF;
DROP TABLE IF EXISTS Basket;
DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Products;
DROP TABLE IF EXISTS Category;
PRAGMA foreign_keys = ON;

-- 1. Kategorie
CREATE TABLE Category (
    CategoryID INTEGER PRIMARY KEY AUTOINCREMENT, 
    CategoryName TEXT UNIQUE
);

-- 2. Produkty
CREATE TABLE Products (
    ProductID INTEGER PRIMARY KEY AUTOINCREMENT,
    ProductName TEXT UNIQUE, 
    Price REAL,
    CategoryID INTEGER, 
    FOREIGN KEY (CategoryID) REFERENCES Category(CategoryID) ON DELETE SET NULL ON UPDATE CASCADE
);

-- 3. Użytkownicy (Dodane!)
CREATE TABLE Users (
    UserID INTEGER PRIMARY KEY AUTOINCREMENT, 
    UserName TEXT UNIQUE,
    Email TEXT
);

-- 4. Koszyk (Poprawiony klucz obcy!)
CREATE TABLE Basket (
    ID INTEGER PRIMARY KEY AUTOINCREMENT, 
    UserID INTEGER, 
    ProductID INTEGER, 
    Quantity INTEGER DEFAULT 1,
    FOREIGN KEY(UserID) REFERENCES Users(UserID) ON DELETE CASCADE ON UPDATE CASCADE, 
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX idx_user_product on Basket(UserID, ProductID);

-- WSTAWIANIE DANYCH --

-- Kategorie
INSERT INTO Category (CategoryName) VALUES 
('Elektronika'), ('Dom i Ogród'), ('Sport i Turystyka'), 
('Moda'), ('Książki'), ('Zdrowie i Uroda'), ('Motoryzacja');

-- Produkty
INSERT INTO Products (ProductName, Price, CategoryID) VALUES 
('Smartfon Galaxy', 2999.99, 1), ('Laptop Pro', 5499.00, 1), ('Słuchawki BT', 199.50, 1),
('Wiertarka Akumulatorowa', 450.00, 2), ('Lampa Stojąca', 120.00, 2),
('Namiot 3-osobowy', 600.00, 3), ('Rower Górski', 2100.00, 3),
('Koszulka Bawełniana', 49.99, 4), ('Jeansy Slim', 129.00, 4),
('Wiedźmin: Ostatnie Życzenie', 39.90, 5), ('Finansowy Ninja', 69.00, 5),
('Krem Nawilżający', 28.50, 6), ('Olej Silnikowy 5W30', 160.00, 7);

-- 5. Dodawanie Userów
INSERT INTO Users (UserName, Email) VALUES 
('Wrex', 'wrex@example.com'),
('Jakub', 'jakub@test.pl'),
('Admin', 'admin@sklep.pl');

-- 6. Dodawanie do Koszyka (Przykładowe zamówienia)
INSERT INTO Basket (UserID, ProductID, Quantity) VALUES 
(1, 1, 1), -- Wrex kupuje Smartfona
(1, 3, 2), -- Wrex kupuje 2 pary słuchawek
(2, 6, 1), -- Jakub kupuje Namiot
(2, 10, 1); -- Jakub kupuje Wiedźmina