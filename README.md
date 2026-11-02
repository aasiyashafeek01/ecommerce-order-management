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
- **Order Management:** Create orders, add multiple items, calculate order totals, and view order history.
- **Order Fulfillment:** Update order statuses, confirm orders, reserve inventory, and handle cancellations.
- **Inventory Management:** Monitor stock levels, update quantities, and view inventory value.
- **Search and Filtering:** Search orders by order number or customer name and filter by order status.

## Project Structure

ecommerce-order-management/
├── backend/
│   ├── database/
│   ├── handlers/
│   ├── models/
│   ├── main.go
│   └── go.mod
├── frontend/
│   └── src/
├── screenshots/
└── README.md

## Prerequisites

- Go
- Node.js and npm
- Angular CLI
- MySQL Server

## Setup and Installation

### 1. Database Setup

- Create a MySQL database named `ecommerce_db`.
- Create the required tables: `products`, `customers`, `orders`, and `order_items`.
- Configure the database connection using environment variables: `MYSQL_USER`, `MYSQL_PASSWORD`, `MYSQL_HOST`, `MYSQL_PORT`, and `MYSQL_DATABASE`.

### 2. Run the Backend

cd backend
go run .

The backend runs at `http://localhost:8080`.

### 3. Run the Frontend

Open another terminal:

cd frontend
npm install
ng serve

Open `http://localhost:4200` in your browser.

## API Endpoints

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
- **Order Management:** Create orders, add multiple items, calculate order totals, and view order history.
- **Order Fulfillment:** Update order statuses, confirm orders, reserve inventory, and handle cancellations.
- **Inventory Management:** Monitor stock levels, update quantities, and view inventory value.
- **Search and Filtering:** Search orders by order number or customer name and filter by order status.

## Project Structure

ecommerce-order-management/
├── backend/
│   ├── database/
│   ├── handlers/
│   ├── models/
│   ├── main.go
│   └── go.mod
├── frontend/
│   └── src/
├── screenshots/
└── README.md

## Prerequisites

- Go
- Node.js and npm
- Angular CLI
- MySQL Server

## Setup and Installation

### 1. Database Setup

- Create a MySQL database named `ecommerce_db`.
- Create the required tables: `products`, `customers`, `orders`, and `order_items`.
- Configure the database connection using environment variables: `MYSQL_USER`, `MYSQL_PASSWORD`, `MYSQL_HOST`, `MYSQL_PORT`, and `MYSQL_DATABASE`.

### 2. Run the Backend

cd backend
go run .

The backend runs at `http://localhost:8080`.

### 3. Run the Frontend

Open another terminal:

cd frontend
npm install
ng serve

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

**Developed as part of the Beinex Project Associate Internship Assignment.**

## Screenshots

Refer to the `screenshots/` directory for application screenshots.

---

**Developed as part of the Beinex Full Stack Web Development Internship Assignment.**