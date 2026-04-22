import {Route, Routes} from "react-router-dom";
import ProductsPage from "../pages/ProductsPage.tsx";
import PaymentPage from "../pages/PaymentPage.tsx";
import CartPage from "../pages/CartPage.tsx";
import MainPage from "../pages/MainPage.tsx";

function Router() {
    return (
        <Routes>
            <Route path="/" element={<MainPage />}/>
            <Route path="/products" element={<ProductsPage />}/>
            <Route path="/payment" element={<PaymentPage />}/>
            <Route path="/cart" element={<CartPage />}/>
        </Routes>
    )
}

export default Router