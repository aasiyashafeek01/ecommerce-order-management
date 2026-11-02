import { Routes } from '@angular/router';
import { Dashboard } from './pages/dashboard/dashboard';
import { Products } from './pages/products/products';
import { Customers } from './pages/customers/customers';
import { Orders } from './pages/orders/orders';
import { Inventory } from './pages/inventory/inventory';

export const routes: Routes = [
  { path: '', component: Dashboard },
  { path: 'products', component: Products },
  { path: 'customers', component: Customers },
  { path: 'orders', component: Orders },
  { path: 'inventory', component: Inventory }
];