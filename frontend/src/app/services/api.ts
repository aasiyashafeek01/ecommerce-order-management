import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface DashboardData {
  total_orders: number;
  total_order_value: number;
  pending_orders: number;
  total_products: number;
  total_customers: number;
  low_stock_products: number;
  inventory_value: number;
}

export interface Product {
  id: number;
  name: string;
  sku: string;
  price: number;
  stock_quantity: number;
  created_at: string;
}

export interface ProductRequest {
  name: string;
  sku: string;
  price: number;
  stock_quantity: number;
}

export interface Customer {
  id: number;
  name: string;
  email: string;
  phone: string | null;
  created_at: string;
}

export interface CustomerRequest {
  name: string;
  email: string;
  phone: string | null;
}


export interface OrderItem {
  id: number;
  order_id: number;
  product_id: number;
  quantity: number;
  unit_price: number;
  subtotal: number;
}

export interface Order {
  id: number;
  order_number: string;
  customer_id: number;
  total_amount: number;
  status: string;
  created_at: string;
  items?: OrderItem[];
}

export interface CreateOrderItem {
  product_id: number;
  quantity: number;
}

export interface CreateOrderRequest {
  customer_id: number;
  items: CreateOrderItem[];
}

export interface UpdateOrderStatusRequest {
  status: string;
}

@Injectable({
  providedIn: 'root'
})
export class Api {
  private http = inject(HttpClient);
  private baseUrl = 'http://localhost:8080/api';

  getDashboard(): Observable<DashboardData> {
    return this.http.get<DashboardData>(
      `${this.baseUrl}/dashboard`
    );
  }

  getProducts(): Observable<Product[]> {
    return this.http.get<Product[]>(`${this.baseUrl}/products`);
  }

  createProduct(product: ProductRequest): Observable<Product> {
    return this.http.post<Product>(`${this.baseUrl}/products`, product);
  }

  updateProduct(id: number, product: ProductRequest): Observable<Product> {
    return this.http.put<Product>(`${this.baseUrl}/products/${id}`, product);
  }

  deleteProduct(id: number): Observable<any> {
    return this.http.delete(`${this.baseUrl}/products/${id}`);
  }

  getCustomers(): Observable<Customer[]> {
    return this.http.get<Customer[]>(`${this.baseUrl}/customers`);
  }

  createCustomer(customer: CustomerRequest): Observable<Customer> {
    return this.http.post<Customer>(`${this.baseUrl}/customers`, customer);
  }

  updateCustomer(id: number, customer: CustomerRequest): Observable<Customer> {
    return this.http.put<Customer>(`${this.baseUrl}/customers/${id}`, customer);
  }

  deleteCustomer(id: number): Observable<any> {
    return this.http.delete(`${this.baseUrl}/customers/${id}`);
  }

  
    getOrders(): Observable<Order[]> {
    return this.http.get<Order[]>(`${this.baseUrl}/orders`);
    }

    getOrder(id: number): Observable<Order> {
    return this.http.get<Order>(`${this.baseUrl}/orders/${id}`);
    }

    createOrder(order: CreateOrderRequest): Observable<Order> {
    return this.http.post<Order>(`${this.baseUrl}/orders`, order);
    }

    confirmOrder(id: number): Observable<any> {
    return this.http.post(`${this.baseUrl}/orders/${id}/confirm`, {});
    }

    updateOrderStatus(
    id: number,
    request: UpdateOrderStatusRequest
    ): Observable<any> {
    return this.http.patch(
        `${this.baseUrl}/orders/${id}/status`,
        request
    );
    }

}
