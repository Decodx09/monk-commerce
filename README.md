# Coupon Management System

A RESTful API built with Go and Fiber to manage coupons and apply them to shopping carts. This repository fulfills the backend assignment requirements, allowing for the application of different discount logic including cart-wise, product-wise, and "Buy X, Get Y" (BxGy) conditions.

## Supported Coupon Types & Scenarios

1. **Cart-wise**: Applies a percentage discount to the entire cart if the total amount exceeds a configured threshold.
   - *Scenario*: A user has $120 worth of items in their cart. A coupon requires a min total of $100 and offers 10% off. The system will discount the cart's final price by $12.
2. **Product-wise**: Applies a percentage discount only on a specific product ID if it exists in the cart.
   - *Scenario*: A user buys 3 units of Product A and 2 of Product B. A 20% discount coupon applies specifically to Product A. The system discounts only the total value of Product A by 20%. 
3. **BxGy (Buy X, Get Y)**: Buy specified total quantities from a "Buy" array to get free quantities from a "Get" array. Supports a `repition_limit`.
   - *Scenario*: "Buy 3 of X or Y, Get 1 Z Free" (with repetition limit = 2). The user has 6 units of X and 1 unit of Y in their cart. The system calculates they fulfill the condition $7 / 3 = 2.33$ times. Capped by the repetition limit, the user gets 2 units of Z for free.

## Assumptions

- **Discounts format**: For Cart and Product-wise coupons, the `discount` field always acts as a percentage (e.g., a value of `10` is 10%).
- **BxGy Logic**: The application interprets the required BxGy "buy conditions" as the *sum* of required counts across all matching `buy_products` criteria. 
- **BxGy "Get" application**: Free items in the BxGy tier are evaluated by matching them against the existing cart. The prompt examples implied the free items are physically "added" to the cart (since "Product Z qty becomes 4" while initially 2). We mark their exact retail value as `total_discount`, leaving the base `total_price` as it originally was with the extra items factored in.
- **Overlapping coupons**: This API assumes that only **one active coupon** can be applied at a time, calculating the absolute best available discount and selecting logic correspondingly on a single evaluation flow.

## Edge Cases and Potential Issues

- **BxGy with missing "Get" products**: The "Get" product must already be present in the initial cart request to identify its unit price. Without it, the application doesn't know the monetary value the free item provides (discount given would be 0).
- **Cart price dropping below zero**: Handled securely. If the cumulative discount value somehow exceeds the cart base total, the system hard-clamps the cart's final price to $0.
- **Large item quantities**: Handled safely since product quantities are processed dynamically, returning the correct scaled integers as evaluated per repeating application limit blocks.

## Limitations & Suggestions for Improvement

- **Database Extensibility**: For BxGy items with zero cart prerequisites (giving them out strictly for free without prior cart inclusion), we would ideally attach the standard unit price dynamically through a related Product-Catalog-Database pull.
- **Coupon Stacking**: Reusing independent coupons overlapping each other in a chained configuration could be supported by modifying the base calculation strategy into an ordered chained architecture sequence.
- **Authentication**: Implementing JWT-based middleware authorization is highly recommended to protect CRUD administration of active system coupons.

---

## Getting Started

### 1. Using Docker Compose (Recommended)
Automatically spins up the required Fiber application and a connected MySQL database.
```bash
docker-compose up --build -d
```
The application will be running on `http://localhost:3000`. 

You can run the fully automated API testing suite simulating the exact assignment scenarios using:
```bash
bash test.sh
```

### 2. Running Locally (Without Docker)
Requires a local MySQL server running on port `3306`. Customize environment variables as needed:
```bash
DB_USER=root DB_PASSWORD=my-pass DB_HOST=127.0.0.1 DB_NAME=monk_commerce go run main.go
```

### Running Tests
Execute unit tests validating the service layer calculations via:
```bash
go test ./... -v
```
