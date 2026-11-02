
import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CurrencyPipe, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { forkJoin } from 'rxjs';

import {
  Api,
  Order,
  Customer,
  Product,
  CreateOrderRequest
} from '../../services/api';

interface OrderLine {
  product_id: number;
  quantity: number;
}

@Component({
  selector: 'app-orders',
  standalone: true,
  imports: [FormsModule, CurrencyPipe, DatePipe],
  templateUrl: './orders.html',
  styleUrl: './orders.css'
})
export class Orders implements OnInit {
  private api = inject(Api);
  private cdr = inject(ChangeDetectorRef);

  orders: Order[] = [];
  searchTerm = '';
  statusFilter = 'ALL';
  customers: Customer[] = [];
  products: Product[] = [];

  loading = true;
  saving = false;
  error = '';
  success = '';

  showForm = false;
  selectedOrder: Order | null = null;

  customerId: number | null = null;

  orderLines: OrderLine[] = [
    { product_id: 0, quantity: 1 }
  ];

  ngOnInit(): void {
    this.loadData();
  }

  loadData(): void {
    this.loading = true;
    this.cdr.detectChanges();

    forkJoin({
      orders: this.api.getOrders(),
      customers: this.api.getCustomers(),
      products: this.api.getProducts()
    }).subscribe({
      next: (data) => {
        this.orders = data.orders ?? [];
        this.customers = data.customers ?? [];
        this.products = data.products ?? [];

        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.error = 'Unable to load order information.';
        this.loading = false;
        this.cdr.detectChanges();
      }
    });
  }

  getCustomerName(customerId: number): string {
    return this.customers.find(c => c.id === customerId)?.name
      ?? `Customer #${customerId}`;
  }

  getProductName(productId: number): string {
    return this.products.find(p => p.id === productId)?.name
      ?? `Product #${productId}`;
  }

  getProductPrice(productId: number): number {
    return this.products.find(p => p.id === Number(productId))?.price ?? 0;
  }

  getAvailableStock(productId: number): number {
    return this.products.find(p => p.id === productId)?.stock_quantity ?? 0;
  }

  getLineSubtotal(line: OrderLine): number {
    const product = this.products.find(p => p.id === Number(line.product_id));
    return (product?.price ?? 0) * Number(line.quantity || 0);
  }

  getOrderTotal(): number {
    return this.orderLines.reduce(
      (total, line) => total + this.getLineSubtotal(line),
      0
    );
  }

  addLine(): void {
    this.orderLines = [
      ...this.orderLines,
      { product_id: 0, quantity: 1 }
    ];
  }

  removeLine(index: number): void {
    if (this.orderLines.length === 1) return;

    this.orderLines = this.orderLines.filter((_, i) => i !== index);
  }

  openCreateForm(): void {
    this.customerId = null;
    this.orderLines = [{ product_id: 0, quantity: 1 }];
    this.error = '';
    this.success = '';
    this.showForm = true;
    this.selectedOrder = null;
  }

  closeCreateForm(): void {
    this.showForm = false;
    this.customerId = null;
    this.orderLines = [{ product_id: 0, quantity: 1 }];
  }

  createOrder(): void {
    if (this.saving) return;

    if (!this.customerId) {
      this.error = 'Please select a customer.';
      return;
    }

    if (this.orderLines.some(line => !line.product_id || line.quantity < 1)) {
      this.error = 'Please select a product and enter a valid quantity for every item.';
      return;
    }

    const productIds = this.orderLines.map(line => Number(line.product_id));

    if (new Set(productIds).size !== productIds.length) {
      this.error = 'Please select each product only once.';
      return;
    }

    for (const line of this.orderLines) {
      const product = this.products.find(
        p => p.id === Number(line.product_id)
      );

      if (!product) {
        this.error = 'One or more selected products could not be found.';
        return;
      }

      if (line.quantity > product.stock_quantity) {
        this.error = `Insufficient available stock for ${product.name}.`;
        return;
      }
    }

    const request: CreateOrderRequest = {
      customer_id: Number(this.customerId),
      items: this.orderLines.map(line => ({
        product_id: Number(line.product_id),
        quantity: Number(line.quantity)
      }))
    };

    this.saving = true;
    this.error = '';
    this.success = '';

    this.api.createOrder(request).subscribe({
      next: () => {
        this.saving = false;
        this.closeCreateForm();
        this.success = 'Order created successfully.';
        this.loadData();
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.saving = false;
        this.error =
          err.error?.error ||
          err.error?.message ||
          'Unable to create the order. Please try again.';

        this.cdr.detectChanges();
      }
    });
  }

  viewOrder(order: Order): void {
    this.error = '';
    this.success = '';
    this.selectedOrder = null;

    this.api.getOrder(order.id).subscribe({
      next: (data) => {
        this.selectedOrder = data;
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.error =
          err.error?.error ||
          'Unable to load order details.';

        this.cdr.detectChanges();
      }
    });
  }

  closeOrderDetails(): void {
    this.selectedOrder = null;
  }

  confirmOrder(order: Order): void {
    if (!confirm(`Confirm order ${order.order_number}? This will reserve its inventory.`)) {
      return;
    }

    this.api.confirmOrder(order.id).subscribe({
      next: () => {
        this.success = 'Order confirmed and inventory reserved successfully.';
        this.error = '';
        this.selectedOrder = null;
        this.loadData();
      },
      error: (err) => {
        this.error =
          err.error?.error ||
          'Unable to confirm this order. Please check stock availability.';

        this.cdr.detectChanges();
      }
    });
  }

  changeStatus(order: Order, status: string): void {
    if (status === 'CONFIRMED') {
      this.confirmOrder(order);
      return;
    }

    if (!confirm(`Change order ${order.order_number} to ${status}?`)) {
      return;
    }

    this.api.updateOrderStatus(order.id, { status }).subscribe({
      next: () => {
        this.success = `Order status updated to ${status}.`;
        this.error = '';
        this.selectedOrder = null;
        this.loadData();
      },
      error: (err) => {
        this.error =
          err.error?.error ||
          'Unable to update the order status.';

        this.cdr.detectChanges();
      }
    });
  }

  getNextStatuses(status: string): string[] {
    const transitions: Record<string, string[]> = {
      PLACED: ['CONFIRMED', 'CANCELLED'],
      CONFIRMED: ['PACKED', 'CANCELLED'],
      PACKED: ['SHIPPED', 'CANCELLED'],
      SHIPPED: ['DELIVERED'],
      DELIVERED: [],
      CANCELLED: []
    };

    return transitions[status] ?? [];
  }

  getStatusClass(status: string): string {
    return status.toLowerCase();
  }

  get filteredOrders(): Order[] {
  const search = this.searchTerm.trim().toLowerCase();

  return this.orders.filter(order => {
    const customerName = this.getCustomerName(order.customer_id).toLowerCase();

    const matchesSearch =
      order.order_number.toLowerCase().includes(search) ||
      customerName.includes(search);

    const matchesStatus =
      this.statusFilter === 'ALL' ||
      order.status === this.statusFilter;

    return matchesSearch && matchesStatus;
  });
}
}
