
import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CurrencyPipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Api, Product, ProductRequest } from '../../services/api';

@Component({
  selector: 'app-inventory',
  standalone: true,
  imports: [CurrencyPipe, FormsModule],
  templateUrl: './inventory.html',
  styleUrl: './inventory.css'
})
export class Inventory implements OnInit {
  private api = inject(Api);
  private cdr = inject(ChangeDetectorRef);

  products: Product[] = [];

  loading = true;
  saving = false;
  error = '';
  success = '';

  showForm = false;
  selectedProduct: Product | null = null;
  newStock = 0;

  ngOnInit(): void {
    this.loadInventory();
  }

  loadInventory(): void {
    this.loading = true;
    this.cdr.detectChanges();

    this.api.getProducts().subscribe({
      next: (data) => {
        this.products = data ?? [];
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.error = 'Unable to load inventory.';
        this.loading = false;
        this.cdr.detectChanges();
      }
    });
  }

  get totalProducts(): number {
    return this.products.length;
  }

  get totalUnits(): number {
    return this.products.reduce(
      (total, product) => total + product.stock_quantity, 0
    );
  }

  get lowStockProducts(): Product[] {
    return this.products.filter(
      product => product.stock_quantity > 0 && product.stock_quantity <= 5
    );
  }

  get outOfStockProducts(): Product[] {
    return this.products.filter(
      product => product.stock_quantity === 0
    );
  }

  get inventoryValue(): number {
    return this.products.reduce(
      (total, product) => total + product.price * product.stock_quantity, 0
    );
  }

  getStockClass(quantity: number): string {
    if (quantity === 0) return 'out-of-stock';
    if (quantity <= 5) return 'low-stock';
    return 'in-stock';
  }

  openStockForm(product: Product): void {
    this.selectedProduct = product;
    this.newStock = product.stock_quantity;
    this.error = '';
    this.success = '';
    this.showForm = true;
  }

  closeStockForm(): void {
    this.showForm = false;
    this.selectedProduct = null;
    this.newStock = 0;
  }

  updateStock(): void {
    if (!this.selectedProduct || this.saving) return;

    if (!Number.isInteger(this.newStock) || this.newStock < 0) {
      this.error = 'Stock must be a whole number of zero or more.';
      return;
    }

    const product = this.selectedProduct;

    const request: ProductRequest = {
      name: product.name,
      sku: product.sku,
      price: product.price,
      stock_quantity: this.newStock
    };

    this.saving = true;
    this.error = '';
    this.success = '';

    this.api.updateProduct(product.id, request).subscribe({
      next: () => {
        this.saving = false;
        this.closeStockForm();
        this.success = `Stock updated successfully for ${product.name}.`;
        this.loadInventory();
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.saving = false;
        this.error =
          err.error?.error ||
          err.error?.message ||
          'Unable to update stock. Please try again.';
        this.cdr.detectChanges();
      }
    });
  }
}
