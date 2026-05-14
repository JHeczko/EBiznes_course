import "./cartPage.css"
import type {Basket, Product} from "../interfaces/definitions.ts";
import * as React from "react";
import {useEffect} from "react";
import {addPayment, clearBasket, deleteBasket, getBasket, updateBasket} from "../services/api.ts";
import toast from "react-hot-toast";

type CartPageProps = Readonly<{
    basket: Basket[];
    setBasket: React.Dispatch<React.SetStateAction<Basket[]>>;
}>

function CartPage({basket, setBasket}: CartPageProps) {

    const handleQuantityChange = async (product: Product, userID: number, newQuantity: number) => {
        if (newQuantity < 1) return;

        try {
            const [updatedItem, status] = await updateBasket(product, userID, newQuantity);

            if (status === 200) {
                setBasket((prevBasket) =>
                    prevBasket.map((item) => {
                        if (item.product.id === product.id) {
                            return { ...item, quantity: updatedItem.quantity };
                        }
                        return item;
                    })
                );
            }
        } catch (error) {
            console.error("Błąd aktualizacji:", error);
        }
    };

    const deleteItem = async (product: Product, userID: number) => {
        try {
            const status = await deleteBasket(product, userID);

            if (status === 200) {
                setBasket(prev => prev.filter(x => x.product.id !== product.id));
                toast.success(`Deleted "${product.name}" from cart`);
            }
        } catch (error) {
            console.error("Error while deleting item from basket: ", error)
            toast.error("Could not delete item :(");
        }
    }

    const handlePayment = async (userID: number) => {
        try {
            const [, status] = await addPayment(userID)
            if (status === 200) {
                const statusClear = await clearBasket(userID)
                if (statusClear) {
                    setBasket([])
                }
                toast.success(`Payment was successful.`)
            }
        } catch (error) {
            console.error("Error while finalizing payment: ", error)
            toast.error("Couldn't finalize payment!")
        }
    }

    useEffect(() => {
        const loadBasketFromDb = async (userID: number) => {
            const basketFromDb = await getBasket(userID);
            setBasket(basketFromDb);
        }
        loadBasketFromDb(1);
    }, [])

    return (
        <main className="cart-page">
            <h1>Cart</h1>
            <div className="cart-container">
                {basket.length === 0
                    ? (<h4 className="empty-cart-message">Empty cart :(</h4>)
                    : basket.map((item: Basket) => (
                        <div className="cart-item" key={item.id}>
                            <div className="item-upper">
                                <h3>{item.product.name}</h3>
                            </div>

                            <div className="item-lower">
                                <div className="control-quantity">
                                    <button
                                        type="button"
                                        className="quantity-control-button"
                                        onClick={() => handleQuantityChange(item.product, item.user_id, item.quantity - 1)}
                                        aria-label="Decrease quantity"
                                    >-</button>
                                    <h3>{item.quantity}</h3>
                                    <button
                                        type="button"
                                        className="quantity-control-button"
                                        onClick={() => handleQuantityChange(item.product, item.user_id, item.quantity + 1)}
                                        aria-label="Increase quantity"
                                    >+</button>
                                    <button
                                        type="button"
                                        className="delete-button"
                                        onClick={() => deleteItem(item.product, item.user_id)}
                                        aria-label="Remove item from cart"
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                            <path d="M3 6h18"></path>
                                            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"></path>
                                            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"></path>
                                            <line x1="10" y1="11" x2="10" y2="17"></line>
                                            <line x1="14" y1="11" x2="14" y2="17"></line>
                                        </svg>
                                    </button>
                                </div>
                                <h5>{(item.quantity * item.product.price).toFixed(2)} zł</h5>
                            </div>
                        </div>
                    ))
                }
            </div>
            {basket.length === 0 ? null : (
                <button
                    type="button"
                    className="checkout-button"
                    onClick={() => handlePayment(1)}
                >Checkout</button>
            )}
        </main>
    )
}

export default CartPage;
