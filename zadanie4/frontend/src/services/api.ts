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

function validateNumericId(id: number): number {
    if (!Number.isInteger(id) || id < 0) {
        throw new Error("Invalid ID");
    }

    return id;
}

function sanitizePathId(id: number): string {
    const validatedId = validateNumericId(id);

    return encodeURIComponent(String(validatedId));
}

function validateQuantity(quantity: number | undefined): number | undefined {
    if (quantity === undefined) {
        return undefined
    }

    if (!Number.isInteger(quantity) || quantity < 1) {
        throw new Error("Invalid quantity");
    }

    return quantity;
}

// ===== PRODUCTS =====

export async function getProducts(): Promise<Product[]> {
    const res = await api.get<Product[]>('/products');
    return res.data;
}

export async function getProductById(id: number): Promise<Product> {
    const safeId = sanitizePathId(id);

    const res = await api.get<Product>(`/products/${safeId}`);
    return res.data;
}


// ===== CATEGORY =====

export async function getCategories(): Promise<Category[]> {
    const res = await api.get<Category[]>('/category');
    return res.data;
}


// ===== BASKET =====

export async function getBasket(userID: number): Promise<Basket[]> {
    const safeId = sanitizePathId(userID);
    const res = await api.get<Basket[]>(`/cart/${safeId}`);
    return res.data;
}

export async function addToBasket(product: Product, userID: number, quantity?: number): Promise<number> {
    const safeId = sanitizePathId(userID);
    const safeQuantity = validateQuantity(quantity);

    const res = await api.post(`/cart/${safeId}`, null, {
        params: {
            prod_id: product.id,
            quantity: safeQuantity
        }
    });

    return res.status;
}

export async function updateBasket(product: Product, userID: number, quantity: number): Promise<[Basket, number]> {
    const validUserID = Number(userID);
    if (Number.isNaN(validUserID)) throw new Error("Invalid User ID");
    
    const safeId = sanitizePathId(validUserID);
    const safeQuantity = validateQuantity(quantity);
    
    const res = await api.patch(`/cart/${safeId}`, null, {
        params: {
            prod_id: product.id,
            quantity: safeQuantity
        }
    });
    return [res.data, res.status];
}

export async function deleteBasket(product: Product, userID: number): Promise<number> {
    const validUserID = Number(userID);
    if (Number.isNaN(validUserID)) throw new Error("Invalid User ID");
    
    const safeId = sanitizePathId(validUserID);
    const res = await api.delete(`/cart/${safeId}`, {
        params: {
            prod_id: product.id
        }
    });

    return res.status;
}

export async function clearBasket(userID: number): Promise<number> {
    const safeId = sanitizePathId(userID);
    const res = await api.delete(`/cart/${safeId}`);
    return res.status;
}


// ===== PAYMENTS =====

export async function addPayment(userID: number): Promise<[Payments, number]> {
    const safeId = sanitizePathId(userID);
    const res = await api.post(`/payments/${safeId}`);
    return [res.data, res.status];
}