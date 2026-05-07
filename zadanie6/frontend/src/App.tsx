import './App.css'
import Router from './router/Router.tsx'
import PageHeader from "./components/PageHeader.tsx";
import {BrowserRouter} from "react-router-dom";
import {useState} from "react";
import type {Basket, Category, Product} from "./interfaces/definitions.ts";
import {Toaster} from "react-hot-toast";


function App() {
    const [products, setProducts] = useState<Product[]>([]);
    const [categories, setCategories] = useState<Category[]>([]);
    const [basket, setBasket] = useState<Basket[]>([]);


    return (
        <>
            <BrowserRouter>
                <PageHeader/>
                <Toaster
                    position="top-right"
                    reverseOrder={false}
                    gutter={6}
                    toastOptions={{
                        duration: 1500,
                        style: {
                            maxWidth: '30rem',
                            maxLines: 6,
                            fontSize: '14px',
                            padding: '10px 14px',
                        },
                    }}
                />                <Router
                    products={products}
                    setProducts={setProducts}
                    categories={categories}
                    setCategories={setCategories}
                    basket={basket}
                    setBasket={setBasket}
                />
            </BrowserRouter>
        </>
    )
}

export default App
