-- Kompletny cykl testu
DROP TABLE Category;
DROP TABLE Products;

-- Włączenie wsparcia dla kluczy obcych (obowiązkowe w SQLite!)
PRAGMA foreign_keys = ON;

-- Tabele
CREATE TABLE IF NOT EXISTS Category(CategoryID int primary key, CategoryName text UNIQUE);

CREATE TABLE IF NOT EXISTS Products (ProductID INTEGER PRIMARY KEY,ProductName TEXT, Price INTEGER, CategoryID INTEGER, FOREIGN KEY (CategoryID) REFERENCES Category(CategoryID) ON DELETE SET NULL ON UPDATE CASCADE);

-- Czyścimy tabele na wypadek, gdybyś już coś tam miał
DELETE FROM Products;
DELETE FROM Category;

-- 1. Dodawanie kategorii (7 kategorii)
INSERT INTO Category (CategoryID, CategoryName) VALUES 
(1, 'Elektronika'),
(2, 'Dom i Ogród'),
(3, 'Sport i Turystyka'),
(4, 'Moda'),
(5, 'Książki'),
(6, 'Zdrowie i Uroda'),
(7, 'Motoryzacja');

-- 2. Dodawanie produktów (30 produktów rozrzuconych po kategoriach)
INSERT INTO Products (ProductName, Price, CategoryID) VALUES 
-- 1. Elektronika
('Smartfon Galaxy', 2999.99, 1), 
('Laptop Pro', 5499.00, 1), 
('Słuchawki BT', 199.50, 1), 
('Monitor 4K', 1200.00, 1), 
('Klawiatura Mechaniczna', 350.00, 1),

-- 2. Dom i Ogród
('Wiertarka Akumulatorowa', 450.00, 2), 
('Lampa Stojąca', 120.00, 2), 
('Zestaw Noży', 299.00, 2), 
('Kosiarka', 899.99, 2), 
('Poduszka Dekoracyjna', 45.00, 2),

-- 3. Sport i Turystyka
('Namiot 3-osobowy', 600.00, 3), 
('Rower Górski', 2100.00, 3), 
('Plecak Trekkingowy', 320.00, 3), 
('Hantle 5kg', 80.00, 3), 
('Mata do Jogi', 55.00, 3),

-- 4. Moda
('Koszulka Bawełniana', 49.99, 4), 
('Jeansy Slim', 129.00, 4), 
('Kurtka Przeciwdeszczowa', 180.00, 4), 
('Skarpetki Sportowe', 15.00, 4), 
('Czapka z Daszkiem', 35.00, 4),

-- 5. Książki
('Wiedźmin: Ostatnie Życzenie', 39.90, 5), 
('Finansowy Ninja', 69.00, 5), 
('Kuchnia Polska - Przepisy', 25.00, 5), 
('Atlas Świata', 110.00, 5), 
('Kryminał pod napięciem', 32.00, 5),

-- 6. Zdrowie i Uroda
('Krem Nawilżający', 28.50, 6), 
('Szampon do Włosów', 18.00, 6), 
('Szczoteczka Soniczna', 149.99, 6),

-- 7. Motoryzacja
('Olej Silnikowy 5W30', 160.00, 7), 
('Wycieraczki Samochodowe', 45.00, 7);