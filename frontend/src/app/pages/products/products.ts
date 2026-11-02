import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';import { CurrencyPipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Api, Product, ProductRequest } from '../../services/api';

@Component({
  selector: 'app-products',
  standalone: true,
  imports: [CurrencyPipe, FormsModule],
  templateUrl: './products.html',
  styleUrl: './products.css'
})
export class Products implements OnInit {
  private api = inject(Api);

  products: Product[] = [];
  loading = true;
  error = '';
  success = '';

  saving = false;
  deletingId: number | null = null;
  private cdr = inject(ChangeDetectorRef);

  showForm = false;
  editingId: number | null = null;

  form: ProductRequest = this.emptyForm();

  emptyForm(): ProductRequest {
    return {
      name: '',
      sku: '',
      price: 0,
      stock_quantity: 0
    };
  }

  ngOnInit(): void {
    this.loadProducts();
  }

  loadProducts(): void {
    this.loading = true;
    this.cdr.detectChanges();

    this.api.getProducts().subscribe({
      next: (data) => {
        this.products = data ?? [];
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.error = 'Unable to load products.';
        this.loading = false;
        this.cdr.detectChanges();
      }
    });
  }

  openAddForm(): void {
    this.editingId = null;
    this.form = this.emptyForm();
    this.error = '';
    this.success = '';
    this.showForm = true;
  }

  openEditForm(product: Product): void {
    this.editingId = product.id;
    this.form = {
      name: product.name,
      sku: product.sku,
      price: product.price,
      stock_quantity: product.stock_quantity
    };
    this.error = '';
    this.success = '';
    this.showForm = true;
  }

  closeForm(): void {
    this.showForm = false;
    this.editingId = null;
    this.form = this.emptyForm();
  }

  saveProduct(): void {
    if (this.saving) return;

    const name = String(this.form.name ?? '').trim();
    const sku = String(this.form.sku ?? '').trim();
    const price = Number(this.form.price);
    const stock = Number(this.form.stock_quantity);

    if (!name) {
      this.error = 'Please enter a product name.';
      return;
    }

    if (!sku) {
      this.error = 'Please enter a SKU.';
      return;
    }

    if (!Number.isFinite(price) || price <= 0) {
      this.error = 'Please enter a valid price greater than zero.';
      return;
    }

    if (!Number.isInteger(stock) || stock < 0) {
      this.error = 'Stock quantity must be a whole number of zero or more.';
      return;
    }

    const request: ProductRequest = {
      name,
      sku,
      price,
      stock_quantity: stock
    };

    const isEditing = this.editingId !== null;

    const request$ = isEditing
      ? this.api.updateProduct(this.editingId!, request)
      : this.api.createProduct(request);

    this.saving = true;
    this.error = '';
    this.success = '';

    request$.subscribe({
      next: () => {
        this.saving = false;
        this.closeForm();

        this.success = isEditing
          ? 'Product updated successfully.'
          : 'Product added successfully.';

        this.loadProducts();
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.saving = false;
        this.error =
          typeof err.error === 'string'
            ? err.error
            : err.error?.error ||
              err.error?.message ||
              'The product could not be saved. Please try again.';
        this.cdr.detectChanges();
      }
    });
  }

 deleteProduct(product: Product): void {
    if (this.deletingId !== null) return;

    const confirmed = confirm(`Delete "${product.name}"?`);

    if (!confirmed) return;

    this.deletingId = product.id;
    this.error = '';
    this.success = '';

    this.api.deleteProduct(product.id).subscribe({
      next: () => {
        this.deletingId = null;
        this.success = `"${product.name}" deleted successfully.`;

        this.loadProducts();
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.deletingId = null;

        this.error =
          typeof err.error === 'string'
            ? err.error
            : err.error?.error ||
              err.error?.message ||
              'Unable to delete this product. Please check whether it is linked to an order.';
        this.cdr.detectChanges();
      }
    });
  }
}