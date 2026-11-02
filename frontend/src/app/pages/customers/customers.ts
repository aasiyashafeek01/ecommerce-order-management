
import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Api, Customer, CustomerRequest } from '../../services/api';

@Component({
  selector: 'app-customers',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './customers.html',
  styleUrl: './customers.css'
})
export class Customers implements OnInit {
  private api = inject(Api);
  private cdr = inject(ChangeDetectorRef);

  customers: Customer[] = [];

  loading = true;
  error = '';
  success = '';

  saving = false;
  deletingId: number | null = null;

  showForm = false;
  editingId: number | null = null;

  form: CustomerRequest = this.emptyForm();

  emptyForm(): CustomerRequest {
    return {
      name: '',
      email: '',
      phone: ''
    };
  }

  ngOnInit(): void {
    this.loadCustomers();
  }

  loadCustomers(): void {
    this.loading = true;
    this.cdr.detectChanges();

    this.api.getCustomers().subscribe({
      next: (data) => {
        this.customers = data ?? [];
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.error = 'Unable to load customers.';
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

  openEditForm(customer: Customer): void {
    this.editingId = customer.id;

    this.form = {
      name: customer.name,
      email: customer.email,
      phone: customer.phone ?? ''
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

  saveCustomer(): void {
    if (this.saving) return;

    const name = String(this.form.name ?? '').trim();
    const email = String(this.form.email ?? '').trim();
    const phone = String(this.form.phone ?? '').trim();

    if (!name) {
      this.error = 'Please enter a customer name.';
      return;
    }

    if (!email) {
      this.error = 'Please enter an email address.';
      return;
    }

    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    if (!emailPattern.test(email)) {
      this.error = 'Please enter a valid email address.';
      return;
    }

    const request: CustomerRequest = {
      name,
      email,
      phone: phone || null
    };

    const isEditing = this.editingId !== null;

    const request$ = isEditing
      ? this.api.updateCustomer(this.editingId!, request)
      : this.api.createCustomer(request);

    this.saving = true;
    this.error = '';
    this.success = '';

    request$.subscribe({
      next: () => {
        this.saving = false;
        this.closeForm();

        this.success = isEditing
          ? 'Customer updated successfully.'
          : 'Customer added successfully.';

        this.loadCustomers();
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.saving = false;

        this.error =
          typeof err.error === 'string'
            ? err.error
            : err.error?.error ||
              err.error?.message ||
              'The customer could not be saved. Please try again.';
        this.cdr.detectChanges();
      }
    });
  }

  deleteCustomer(customer: Customer): void {
    if (this.deletingId !== null) return;

    const confirmed = confirm(`Delete "${customer.name}"?`);

    if (!confirmed) return;

    this.deletingId = customer.id;
    this.error = '';
    this.success = '';

    this.api.deleteCustomer(customer.id).subscribe({
      next: () => {
        this.deletingId = null;
        this.success = `"${customer.name}" deleted successfully.`;

        this.loadCustomers();
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.deletingId = null;

        this.error =
          typeof err.error === 'string'
            ? err.error
            : err.error?.error ||
              err.error?.message ||
              'Unable to delete this customer. The customer may be linked to an order.';
                this.cdr.detectChanges();
      }
    });
  }
}
