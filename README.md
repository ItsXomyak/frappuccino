# frappuccino: Coffee Shop Management System

## Table of Contents

- [Introduction](#introduction)
- [Project Overview](#project-overview)
- [Features](#features)
- [Architecture](#architecture)
  - [Layered Architecture](#layered-architecture)
  - [Data Models](#data-models)
- [API Endpoints](#api-endpoints)
  - [Orders](#orders)
  - [Menu Items](#menu-items)
  - [Inventory](#inventory)
  - [Aggregations](#aggregations)
- [Data Storage](#data-storage)
- [Logging](#logging)
- [Error Handling and Validation](#error-handling-and-validation)
- [Usage](#usage)
- [Getting Started](#getting-started)
- [Dependencies](#dependencies)
- [Contributing](#contributing)
- [License](#license)

---

## Introduction

Managing a coffee shop involves juggling multiple tasks: taking orders, tracking inventory, updating the menu, and ensuring customer satisfaction. The `frappuccino` project is a simplified coffee shop management system designed to provide hands-on experience with backend development, RESTful API design, and data management using Go.

## Project Overview

`frappuccino` is a command-line application that starts an HTTP server to manage a coffee shop's operations. It allows staff to:

- **Manage Orders**: Create, retrieve, update, and delete customer orders.
- **Oversee Inventory**: Track ingredient stock levels to prevent shortages.
- **Update the Menu**: Add new menu items, adjust prices, and manage offerings.
- **Generate Reports**: View total sales and popular menu items.

This project is a practical exploration of how order management and inventory systems operate under the hood, including how they handle data processing, manage concurrent requests, and ensure data integrity.

## Features

- **RESTful API**: Provides endpoints to manage orders, menu items, and inventory.
- **JSON Data Handling**: Encodes and decodes JSON data for seamless data transmission.
- **Data Storage with JSON Files**: Stores data locally in JSON files without the need for a database.
- **Layered Architecture**: Implements a three-layered architecture for clean code and scalability.
- **Logging**: Uses Go's `log/slog` package to record significant events and errors.
- **Aggregations**: Provides endpoints to retrieve total sales and popular menu items.

## Architecture

The application is built using a three-layered architecture:

1. **Presentation Layer (Handlers)**
   - **Responsibilities**:
     - Handle HTTP requests and responses.
     - Parse input data and format output data.
     - Invoke appropriate methods from the Business Logic Layer.
   - **Implementation Details**:
     - Handlers are organized based on entities (e.g., `order_handler.go`, `menu_handler.go`, `inventory_handler.go`).
     - Uses Go's `net/http` package to set up routes and handle requests.
     - Validates input data and returns meaningful error messages.

2. **Business Logic Layer (Services)**
   - **Responsibilities**:
     - Implement core business logic and rules.
     - Define interfaces for services to promote decoupling.
     - Perform data processing and call methods from the Data Access Layer.
     - Handle aggregations and computations based on business requirements.
   - **Implementation Details**:
     - Defines service interfaces (e.g., `OrderService`, `MenuService`, `InventoryService`) in separate files.
     - Implements methods for aggregations (e.g., `GetTotalSales`, `GetPopularMenuItems`).
     - Ensures that services are independent and can be tested in isolation.

3. **Data Access Layer (Repositories)**
   - **Responsibilities**:
     - Manage data storage and retrieval operations.
     - Interact with JSON files to persist and read data.
     - Ensure data integrity and consistency.
     - Provide interfaces for repositories to allow flexibility.
   - **Implementation Details**:
     - Creates repository interfaces for each entity (e.g., `OrderRepository`, `MenuRepository`, `InventoryRepository`).
     - Implements these interfaces.
     - Organizes data in separate JSON files for each entity, stored in the `data/` directory.

### Project Structure

```plaintext
frappuccino/
├── cmd/
│   └── main.go
├── internal/
│   ├── handler/
│   │   ├── order_handler.go
│   │   ├── menu_handler.go
│   │   └── inventory_handler.go
│   ├── service/
│   │   ├── order_service.go
│   │   ├── menu_service.go
│   │   ├── inventory_service.go
│   │   └── ...
│   └── dal/
│       ├── order_repository.go
│       └── ...
├── models/
│   ├── order.go
│   ├── menu_item.go
│   ├── inventory_item.go
│   └── ...
├── go.mod
├── go.sum
└── ...
```

### Orders

- **GET /orders**: Retrieve all orders.
- **GET /orders/{id}**: Retrieve a specific order by ID.
- **PUT /orders/{id}**: Update an existing order.
- **DELETE /orders/{id}**: Delete an order.
- **POST /orders/{id}/close**: Close an order.

### Menu Items

- **POST /menu**: Add a new menu item.
- **GET /menu**: Retrieve all menu items.
- **GET /menu/{id}**: Retrieve a specific menu item.
- **PUT /menu/{id}**: Update a menu item.
- **DELETE /menu/{id}**: Delete a menu item.

### Inventory

- **POST /inventory**: Add a new inventory item.
- **GET /inventory**: Retrieve all inventory items.
- **GET /inventory/{id}**: Retrieve a specific inventory item.
- **PUT /inventory/{id}**: Update an inventory item.
- **DELETE /inventory/{id}**: Delete an inventory item.

### Aggregations

- **GET /reports/total-sales**: Get the total sales amount.
- **GET /reports/popular-items**: Get a list of popular menu items.

## Data Storage with JSON Files

### No Database Usage

The application stores all data locally in JSON files without using a database system. Data is persisted in the following files within the `data/` directory:

- **orders.json**: Stores information about customer orders.
- **menu_items.json**: Stores details about the products available in the coffee shop.
- **inventory.json**: Stores inventory of ingredients required to prepare menu items.

### Examples of JSON Files

#### orders.json

```json
[
  {
    "order_id": "order123",
    "customer_name": "Alice Smith",
    "items": [
      {
        "product_id": "latte",
        "quantity": 2
      },
      {
        "product_id": "muffin",
        "quantity": 1
      }
    ],
    "status": "open",
    "created_at": "2023-10-01T09:00:00Z"
  },
  {
    "order_id": "order124",
    "customer_name": "Bob Johnson",
    "items": [
      {
        "product_id": "espresso",
        "quantity": 1
      }
    ],
    "status": "closed",
    "created_at": "2023-10-01T09:30:00Z"
  }
]