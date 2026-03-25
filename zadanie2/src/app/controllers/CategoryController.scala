package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import models.{DataBase, Category}

@Singleton
class CategoryController @Inject()(cc: ControllerComponents) extends AbstractController(cc) {

  def getAll = Action {
    Ok(Json.toJson(DataBase.categories))
  }

  def create(name: String) = Action {
    val new_id = if (DataBase.categories.isEmpty) 1 else DataBase.categories.map(_.id).max + 1
    val new_category = Category(new_id, name)

    DataBase.categories.addOne(new_category)
    Created(Json.toJson(new_category))
  }

  def read(id: Int) = Action {
    DataBase.categories.find(_.id == id) match {
      case Some(cat) => Ok(Json.toJson(cat))
      case None      => NotFound(Json.obj("error" -> s"Category $id not found"))
    }
  }

  def update(id: Int, name: Option[String]) = Action {
    val index = DataBase.categories.indexWhere(_.id == id)

    if (index != -1) {
      val old = DataBase.categories(index)
      val name_updated = name.getOrElse(old.name)
      val new_cat = Category(old.id, name_updated)

      DataBase.categories.update(index, new_cat)
      Ok(Json.toJson(new_cat))
    } else {
      NotFound(Json.obj("error" -> s"Category $id not found"))
    }
  }

  def delete(id: Int) = Action {
    if (DataBase.categories.exists(_.id == id)) {
      DataBase.categories.filterInPlace(_.id != id)
      Ok(Json.obj("status" -> "Deleted", "id" -> id))
    } else {
      NotFound(Json.obj("error" -> s"Category $id not found"))
    }
  }
}