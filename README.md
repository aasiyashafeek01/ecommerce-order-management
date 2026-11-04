# E-Commerce Order Management and Fulfillment Platform

A full-stack web application developed to manage products, customers, orders, inventory, and order fulfillment.

## Tech Stack

- **Backend:** Go
- **Frontend:** Angular, TypeScript, HTML, CSS
- **Database:** MySQL
- **Communication:** REST APIs

## Features

- **Dashboard:** View order statistics, sales value, customer and product counts, and inventory information.
- **Product Management:** Add, view, update, and delete products.
- **Customer Management:** Manage customer details and records.
- **Order Management:** Create orders with multiple items, calculate totals, and view order history.
- **Order Fulfillment:** Confirm orders, update order statuses, reserve inventory, and handle cancellations.
- **Inventory Management:** Monitor stock levels, update quantities, and view inventory value.
- **Search and Filtering:** Search orders by order number or customer name and filter by order status.

## Project Structure

```text
ecommerce-order-management/
├── backend/
│   ├── database/
│   ├── handlers/
│   ├── models/
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── database/
│   └── schema.sql
├── frontend/
│   └── src/
├── screenshots/
└── README.md
```

## Prerequisites

- Go
- Node.js and npm
- Angular CLI
- MySQL Server

## Setup and Installation

### 1. Clone the Repository

```bash
git clone https://github.com/aasiyashafeek01/ecommerce-order-management.git
cd ecommerce-order-management
```

### 2. Database Setup

Create the database tables using the provided SQL schema.

Run the following command from the project root using MySQL:

```bash
mysql -u root -p < database/schema.sql
```

This creates the `ecommerce_db` database and its required tables.

### 3. Run the Backend

Open a terminal in the `backend` directory:

```bash
cd backend
```

Configure the database connection using environment variables.

For Windows PowerShell:

```powershell
$env:MYSQL_USER = "root"
$env:MYSQL_DATABASE = "ecommerce_db"

$securePassword = Read-Host "Enter MySQL password" -AsSecureString
$env:MYSQL_PASSWORD = [System.Net.NetworkCredential]::new("", $securePassword).Password
```

Start the backend:

```bash
go run .
```

The backend runs at `http://localhost:8080`.

### 4. Run the Frontend

Open another terminal:

```bash
cd frontend
npm install
ng serve
```

Open `http://localhost:4200` in your browser.

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/products` | List products |
| POST | `/api/products` | Create product |
| PUT | `/api/products/{id}` | Update product |
| DELETE | `/api/products/{id}` | Delete product |
| GET | `/api/customers` | List customers |
| POST | `/api/customers` | Create customer |
| PUT | `/api/customers/{id}` | Update customer |
| DELETE | `/api/customers/{id}` | Delete customer |
| GET | `/api/orders` | List orders |
| GET | `/api/orders/{id}` | View order details |
| POST | `/api/orders` | Create order |
| POST | `/api/orders/{id}/confirm` | Confirm order |
| PATCH | `/api/orders/{id}/status` | Update order status |
| GET | `/api/dashboard` | Dashboard statistics |

## Screenshots

Refer to the `screenshots/` directory for application screenshots.

---

**Developed as part of the Beinex Full Stack Web Development Internship Assignment.**
