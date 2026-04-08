package org.zadanie

class Category(var cat_name: String) {
    val cat_id: Int = nextId()

    companion object {
        private var counter = 0

        private fun nextId(): Int {
            counter++
            return counter
        }
    }
}