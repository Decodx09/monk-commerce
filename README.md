# Coupon Management System

A RESTful API built with Go and Fiber to manage coupons and apply them to shopping carts. 

### Supported Coupon Types
1. **Cart-wise**: Percentage discount on the whole cart above a defined threshold.
2. **Product-wise**: Percentage discount on a specific product.
3. **BxGy (Buy X, Get Y)**: Buy a specified quantity of products to get another set of products for free (supports repetition limits).

## Features
- **CRUD Operations**: Manage coupons via REST API endpoints.
- **Cart Application**: Evaluate carts, calculate possible discounts, and apply applicable coupons automatically.
- **Testing**: A full suite of testing using `test.sh` as well as Go unit tests.
- **Expiration Tracking**: Validates expiration dates (Bonus).

*Note: For BxGy coupons, the "Get" product must already be present in the initial cart request to identify its unit price.*

## Getting Started

### 1. Using Docker Compose (Recommended)
Automatically spins up the required Fiber application and a connected MySQL database.
```bash
docker-compose up --build -d
```
The application will be running on `http://localhost:3000`. 

You can run the fully automated API testing suite using:
```bash
bash test.sh
```

### 2. Running Locally (Without Docker)
Requires a local MySQL server running on port `3306`. Customize environment variables as needed:
```bash
DB_USER=root DB_PASSWORD=my-pass DB_HOST=127.0.0.1 DB_NAME=monk_commerce go run main.go
```

### Running Tests
Execute unit tests via:
```bash
go test ./... -v
```
# monk-commerce
