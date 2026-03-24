package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import scala.collection.mutable
case class Product(id: Int, name: String, price: Double)

object Product {
  implicit val format = Json.format[Product]
}

@Singleton
class ProductController @Inject()(cc: ControllerComponents) extends AbstractController(cc) {

  var db_products = mutable.ListBuffer(
    Product(1, "Laptop", 3000),
    Product(2, "Mouse", 100)
  )

  def getAll = Action {
    Ok(Json.toJson(db_products))
  }

  def create(name: String, count: Int) = Action {
    // Bezpieczne generowanie ID
    val new_id = if (db_products.isEmpty) 1 else db_products.map(_.id).max + 1
    val new_prod = Product(new_id, name, count)

    db_products.addOne(new_prod)
    Ok(Json.toJson(new_prod))
  }

  def read(id: Int) = Action {
    val productOpt = db_products.find(_.id == id)

    if (productOpt.isDefined) {
      Ok(Json.toJson(productOpt.get))
    } else {
      NotFound("No such product")
    }
  }

  def update(id: Int, name: Option[String], count: Option[Int]) = Action {
    val index_to_update = db_products.indexWhere(_.id == id)

    if (index_to_update != -1) {
      val stary_prod = db_products(index_to_update)

      // Logika "jeśli podano, to weź nowe, inaczej stare"
      val final_name = if (name.isDefined) name.get else stary_prod.name
      val final_count = if (count.isDefined) count.get else stary_prod.price.toInt // uwaga na typy!

      val updated_prod = Product(id, final_name, final_count.toDouble)

      db_products.update(index_to_update, updated_prod)
      Ok(Json.toJson(updated_prod))
    } else {
      NotFound("No such product")
    }
  }

  def delete(id: Int) = Action {
    // Usuwamy ten konkretny produkt (zostawiamy te, które mają inne ID)
    db_products.filterInPlace(_.id != id)
    Ok(Json.toJson("Deleted " + id))
  }
}