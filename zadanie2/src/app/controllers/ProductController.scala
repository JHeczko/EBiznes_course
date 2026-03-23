package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._

case class Product(id: Int, name: String, price: Double)

object Product {
  implicit val format = Json.format[Product]
}

@Singleton
class ProductController @Inject()(cc: ControllerComponents) extends AbstractController(cc) {

  var products = List(
    Product(1, "Laptop", 3000),
    Product(2, "Mouse", 100)
  )

  def getAll = {
    Action {Ok(Json.toJson(products))}
  }
}