import axios from 'axios';
import type { Product, Category, Basket, Payments } from "../interfaces/definitions.ts";

// backend
const API_URL = 'http://127.0.0.1:13000';

// 🔥 axios instance
const api = axios.create({
    baseURL: API_URL,
    headers: {
        "Content-Type": "application/json"
    },
    withCredentials: false
});


// ===== PRODUCTS =====

export async function getProducts(): Promise<Product[]> {
    const res = await api.get<Product[]>('/products');
    return res.data;
}

export async function getProductById(id: number): Promise<Product> {
    const res = await api.get<Product>(`/products/${id}`);
    return res.data;
}


// ===== CATEGORY =====

export async function getCategories(): Promise<Category[]> {
    const res = await api.get<Category[]>('/category');
    return res.data;
}


// ===== BASKET =====

export async function getBasket(userID: number): Promise<Basket[]> {
    const res = await api.get<Basket[]>(`/cart/${userID}`);
    return res.data;
}

export async function addToBasket(product: Product, userID: number, quantity?: number): Promise<number> {

    const res = await api.post(`/cart/${userID}`, null, {
        params: {
            prod_id: product.id,
            quantity: quantity
        }
    });

    return res.status;
}

export async function updateBasket(product: Product, userID: number, quantity: number): Promise<[Basket, number]> {
    const res = await api.patch(`/cart/${userID}`, null, {
        params: {
            prod_id: product.id,
            quantity: quantity
        }
    });
    return [res.data, res.status];
}

export async function deleteBasket(product: Product, userID: number): Promise<number> {
    const res = await api.delete(`/cart/${userID}`, {
        params: {
            prod_id: product.id
        }
    });

    return res.status;
}

export async function clearBasket(userID: number): Promise<number> {
    const res = await api.delete(`/cart/${userID}`);
    return res.status;
}


// ===== PAYMENTS =====

export async function addPayment(userID: number): Promise<[Payments, number]> {
    const res = await api.post(`/payments/${userID}`);
    return [res.data, res.status];
}