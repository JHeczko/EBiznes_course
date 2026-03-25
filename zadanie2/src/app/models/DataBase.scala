package models

import scala.collection.mutable.ListBuffer
import play.api.libs.json._

case class Category(id: Int, name: String)
object Category { implicit val format = Json.format[Category] }

// Produkt teraz wie, do jakiej kategorii należy (categoryId)
case class Product(
                    id: Int,
                    name: String,
                    price: Double,
                    categoryId: Int)
object Product { implicit val format = Json.format[Product] }

case class CartItem(
                     productId: Int,
                     quantity: Int)
object CartItem { implicit val format = Json.format[CartItem] }

object DataBase {
  val categories = ListBuffer[Category](
    Category(1, "Elektronika"),
    Category(2, "Akcesoria")
  )

  val products = ListBuffer[Product](
    // 1 -> Elektronika, 2 -> Akcesoria
    Product(1, "Laptop", 3000.0, 1),
    Product(2, "Mouse", 100.0, 2)
  )

  val cart = ListBuffer[CartItem]()
}