import type {Category, Product} from "../interfaces/definitions.ts";
import * as React from "react";
import {addToBasket, getProducts} from "../services/api.ts";

import "./productPage.css"
import {useEffect} from "react";
import toast from "react-hot-toast";

type ProductsPageParams = {
    products: Product[];
    setProducts: React.Dispatch<React.SetStateAction<Product[]>>;
    categories: Category[];
    setCategories: React.Dispatch<React.SetStateAction<Category[]>>;
}

function ProductsPage({products, setProducts}: ProductsPageParams) {
    useEffect(() => {
        const init_data = async () =>{
            const data = await getProducts();
            setProducts(data);
        }
        init_data();
    }, []);

    if (!products.length) return <div>Loading...</div>;

    const addItemToBasket = async (product: Product) => {
        try {
            await addToBasket(product, 1);
            toast.success(`Added "${product.name}" to cart`);
        }catch(err){
            toast.error(`Failed to add "${product.name} to basket"`);
            console.error(err);
        }
    }

    return (
        <main className="products-page">
            <h1 className="title">Our products</h1>
            <div className="filters">

            </div>
            <div className="product-container">
                {products.map((product) => (
                    <div key={product.id} className="product-card">
                        {/* GÓRA: Info o produkcie */}
                        <div className="card-info">
                            <h3>{product.name}</h3>
                            <p className="category">Category: {product.category.name}</p>
                        </div>

                        {/* DÓŁ: Cena i guzik (to musi być razem!) */}
                        <div className="card-footer">
                            <div className="price-tag">{product.price} zł</div>
                            <div className="add-to-cart-btn" onClick={() => addItemToBasket(product)}>
                                Dodaj
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </main>
    )
}

export default ProductsPage;