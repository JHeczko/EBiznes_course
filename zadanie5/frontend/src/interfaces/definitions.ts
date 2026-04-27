// Database structures
export interface Category{
    id: number;
    name: string;
};

export interface Product{
    id: number;
    name: string;
    price: number;
    cat_id: number;
    category: Category;
};

export interface User{
    user_id: number;
    user_name: string;
    mail: string;
};

export interface Basket{
    id: number;
    user_id : number;
    prod_id: number;
    quantity: number;
    product: Product;
};

export interface Payments{
    id: number;
    user_id: number;
    total_amount: number;
    status: string;
    created_at: Date;
}


// Responses structures
export interface DeleteResponse{
    user_id: number;
    prod_id: number;
    deleted: number;
    type: "all" | "single";
}

