#!/bin/sh

# Ścieżka do bazy danych
DB_PATH="/root/app/database.db"

# Sprawdzamy, czy baza danych już istnieje
if [ ! -f "$DB_PATH" ]; then
    echo "Baza nie istnieje. Inicjalizacja za pomocą sqlite3..."
    # Tworzymy plik bazy i wstrzykujemy skrypt SQL
    sqlite3 "$DB_PATH" < /root/app/data_init.sql
    echo "Baza zainicjalizowana pomyślnie."
else
    echo "Baza już istnieje, pomijam inicjalizację."
fi