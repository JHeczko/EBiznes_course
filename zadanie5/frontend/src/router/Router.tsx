import {Route, Routes} from "react-router-dom";
import ProductsPage from "../pages/ProductsPage.tsx";
import CartPage from "../pages/CartPage.tsx";
import MainPage from "../pages/MainPage.tsx";
import * as React from "react";
import type {Product, Category, Basket} from "../interfaces/definitions.ts";
import InfoPage from "../pages/InfoPage.tsx";

interface RouterProps {
    products: Product[];
    setProducts: React.Dispatch<React.SetStateAction<Product[]>>;
    categories: Category[];
    setCategories: React.Dispatch<React.SetStateAction<Category[]>>;
    basket: Basket[];
    setBasket: React.Dispatch<React.SetStateAction<Basket[]>>;
    //payments: Payments[];
    //setPayments: React.Dispatch<React.SetStateAction<Payments[]>>;
}

function Router({products, setProducts, categories, setCategories, basket, setBasket}: RouterProps) {
    return (
        <Routes>
            <Route path="/" element={<MainPage />}/>
            <Route path="/products" element={<ProductsPage
                products={products}
                setProducts={setProducts}
                categories={categories}
                setCategories={setCategories}
            />}/>
            <Route path="/cart" element={<CartPage
                basket={basket}
                setBasket={setBasket}
            />}/>
            <Route path="/info" element={<InfoPage />}></Route>
        </Routes>
    )
}

export default Router