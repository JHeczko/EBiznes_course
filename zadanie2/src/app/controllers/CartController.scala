package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import models.{DataBase, CartItem}

@Singleton
class CartController @Inject()(cc: ControllerComponents) extends AbstractController(cc) {

  // GET /cart/all
  def all = Action {
    Ok(Json.toJson(DataBase.cart))
  }

  // POST /cart/create?productId=1&quantity=2
  def create(productId: Int, quantity: Int) = Action {
    if (DataBase.products.exists(_.id == productId)) {
      val index = DataBase.cart.indexWhere(_.productId == productId)

      if (index != -1) {
        val old = DataBase.cart(index)
        val updated_item = CartItem(productId, old.quantity + quantity)
        DataBase.cart.update(index, updated_item)
        Ok(Json.toJson(updated_item))
      } else {
        val new_item = CartItem(productId, quantity)
        DataBase.cart.addOne(new_item)
        Created(Json.toJson(new_item))
      }
    } else {
      NotFound(Json.obj("error" -> s"Product $productId not found in store"))
    }
  }

  // GET /cart/read?productId=1
  def read(productId: Int) = Action {
    DataBase.cart.find(_.productId == productId) match {
      case Some(item) => Ok(Json.toJson(item))
      case None       => NotFound(Json.obj("error" -> s"Product $productId not in cart"))
    }
  }

  // PATCH /cart/update?productId=1&quantity=5
  def update(productId: Int, quantity: Int) = Action {
    val index = DataBase.cart.indexWhere(_.productId == productId)

    if (index != -1) {
      val updated_item = CartItem(productId, quantity)
      DataBase.cart.update(index, updated_item)
      Ok(Json.toJson(updated_item))
    } else {
      NotFound(Json.obj("error" -> s"Product $productId not in cart"))
    }
  }

  // DELETE /cart/delete?productId=1
  def delete(productId: Int) = Action {
    if (DataBase.cart.exists(_.productId == productId)) {
      DataBase.cart.filterInPlace(_.productId != productId)
      Ok(Json.obj("status" -> "Deleted", "productId" -> productId))
    } else {
      NotFound(Json.obj("error" -> s"Product $productId not in cart"))
    }
  }
}