package org.zadanie

class Product(var prod_name: String, var price: Double,var cat_id: Int){
    var prod_id: Int = nextGlobalId()

    companion object {
        private var counter = 0;

        private fun nextGlobalId(): Int {
            counter++;
            return counter
        }
    }
}