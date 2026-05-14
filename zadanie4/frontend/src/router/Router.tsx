import {Route, Routes} from "react-router-dom";
import ProductsPage from "../pages/ProductsPage.tsx";
import CartPage from "../pages/CartPage.tsx";
import MainPage from "../pages/MainPage.tsx";
import * as React from "react";
import type {Product, Basket} from "../interfaces/definitions.ts";
import InfoPage from "../pages/InfoPage.tsx";

interface RouterProps {
    products: Product[];
    setProducts: React.Dispatch<React.SetStateAction<Product[]>>;
    basket: Basket[];
    setBasket: React.Dispatch<React.SetStateAction<Basket[]>>;
}

function Router({products, setProducts, basket, setBasket}: Readonly<RouterProps>) {
    return (
        <Routes>
            <Route path="/" element={<MainPage />}/>
            <Route path="/products" element={<ProductsPage
                products={products}
                setProducts={setProducts}
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