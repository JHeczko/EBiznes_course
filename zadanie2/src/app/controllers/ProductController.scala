package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import models.{DataBase, Product}

@Singleton
class ProductController @Inject()(cc: ControllerComponents) extends AbstractController(cc) {

  def getAll = Action {
    Ok(Json.toJson(DataBase.products))
  }

  // POST: /products/create?name=X&price=1.0&catId=1
  def create(name: String, price: Double, catId: Int) = Action {
    // Walidacja: czy kategoria istnieje w bazie?
    if (DataBase.categories.exists(_.id == catId)) {
      val new_id = if (DataBase.products.isEmpty) 1 else DataBase.products.map(_.id).max + 1
      val new_prod = Product(new_id, name, price, catId)
      DataBase.products.addOne(new_prod)
      Ok(Json.toJson(new_prod))
    } else {
      BadRequest(Json.obj("error" -> "Category does not exist"))
    }
  }

  def read(id: Int) = Action {
    DataBase.products.find(_.id == id) match {
      case Some(p) => Ok(Json.toJson(p))
      case None => NotFound("No such product")
    }
  }

  // PATCH: /products/update?id=1&name=X&price=2.0&catId=1
  def update(id: Int, name: Option[String], price: Option[Double], catId: Option[Int]) = Action {
    val index = DataBase.products.indexWhere(_.id == id)

    if (index != -1) {
      val old = DataBase.products(index)

      // Jeśli podano catId, sprawdź czy kategoria istnieje
      val validCatId = catId match {
        case Some(cid) if DataBase.categories.exists(_.id == cid) => cid
        case _ => old.categoryId
      }

      val updated = Product(
        id,
        name.getOrElse(old.name),
        price.getOrElse(old.price),
        validCatId
      )

      DataBase.products.update(index, updated)
      Ok(Json.toJson(updated))
    } else {
      NotFound("No such product")
    }
  }

  def delete(id: Int) = Action {
    DataBase.products.filterInPlace(_.id != id)
    Ok(Json.obj("message" -> s"Deleted $id"))
  }
}