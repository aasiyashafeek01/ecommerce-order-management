import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { Api, DashboardData } from '../../services/api';
import { CurrencyPipe } from '@angular/common';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CurrencyPipe],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css'
})
export class Dashboard implements OnInit {
  private api = inject(Api);
  private cdr = inject(ChangeDetectorRef);

  dashboard: DashboardData | null = null;
  loading = true;
  error = '';

  ngOnInit(): void {
    this.loadDashboard();
  }

  loadDashboard(): void {
    this.loading = true;
    this.error = '';
    this.cdr.detectChanges();

    this.api.getDashboard().subscribe({
      next: (data) => {
        this.dashboard = data;
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error('Dashboard API error:', err);
        this.error = 'Unable to load dashboard data. Please check the backend server.';
        this.loading = false;
        this.cdr.detectChanges();
      }
    });
  }
}