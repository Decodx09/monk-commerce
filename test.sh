#!/bin/bash

# Base URL
URL="http://localhost:3000"

echo "--------------------------------------------------------"
echo "1. Create Cart-wise Coupon"
CART_WISE_JSON='{
    "type": "cart-wise",
    "details": {
        "threshold": 100,
        "discount": 10
    }
}'
CART_COUPON_ID=$(curl -s -X POST $URL/coupons \
  -H "Content-Type: application/json" \
  -d "$CART_WISE_JSON" | jq -r '.coupon_id')
echo "Created Cart-wise Coupon ID: $CART_COUPON_ID"


echo "--------------------------------------------------------"
echo "2. Create Product-wise Coupon"
PRODUCT_WISE_JSON='{
    "type": "product-wise",
    "details": {
        "product_id": 1,
        "discount": 20
    }
}'
PROD_COUPON_ID=$(curl -s -X POST $URL/coupons \
  -H "Content-Type: application/json" \
  -d "$PRODUCT_WISE_JSON" | jq -r '.coupon_id')
echo "Created Product-wise Coupon ID: $PROD_COUPON_ID"


echo "--------------------------------------------------------"
echo "3. Create BxGy Coupon"
BXGY_JSON='{
    "type": "bxgy",
    "details": {
        "buy_products": [
            {"product_id": 1, "quantity": 3},
            {"product_id": 2, "quantity": 1}
        ],
        "get_products": [
            {"product_id": 3, "quantity": 1}
        ],
        "repition_limit": 2
    }
}'
BXGY_COUPON_ID=$(curl -s -X POST $URL/coupons \
  -H "Content-Type: application/json" \
  -d "$BXGY_JSON" | jq -r '.coupon_id')
echo "Created BxGy Coupon ID: $BXGY_COUPON_ID"


echo "--------------------------------------------------------"
echo "4. Get All Coupons"
curl -s -X GET $URL/coupons | jq .


echo "--------------------------------------------------------"
echo "5. Get Specific Coupon (ID: $CART_COUPON_ID)"
curl -s -X GET $URL/coupons/$CART_COUPON_ID | jq .


echo "--------------------------------------------------------"
echo "6. Update Coupon (ID: $CART_COUPON_ID)"
UPDATE_JSON='{
    "type": "cart-wise",
    "details": {
        "threshold": 150,
        "discount": 15
    }
}'
curl -s -X PUT $URL/coupons/$CART_COUPON_ID \
  -H "Content-Type: application/json" \
  -d "$UPDATE_JSON" | jq .


echo "--------------------------------------------------------"
echo "7. Applicable Coupons (Check applicable coupons for a cart)"
CART_JSON='{
    "cart": {
        "items": [
            {"product_id": 1, "quantity": 6, "price": 50},
            {"product_id": 2, "quantity": 3, "price": 30},
            {"product_id": 3, "quantity": 1, "price": 25}
        ]
    }
}'
curl -s -X POST $URL/applicable-coupons \
  -H "Content-Type: application/json" \
  -d "$CART_JSON" | jq .


echo "--------------------------------------------------------"
echo "8. Apply Coupon: Cart-wise (ID: $CART_COUPON_ID)"
curl -s -X POST $URL/apply-coupon/$CART_COUPON_ID \
  -H "Content-Type: application/json" \
  -d "$CART_JSON" | jq .

echo "--------------------------------------------------------"
echo "9. Apply Coupon: Product-wise (ID: $PROD_COUPON_ID)"
curl -s -X POST $URL/apply-coupon/$PROD_COUPON_ID \
  -H "Content-Type: application/json" \
  -d "$CART_JSON" | jq .

echo "--------------------------------------------------------"
echo "10. Apply Coupon: BxGy (ID: $BXGY_COUPON_ID)"
curl -s -X POST $URL/apply-coupon/$BXGY_COUPON_ID \
  -H "Content-Type: application/json" \
  -d "$CART_JSON" | jq .


echo "--------------------------------------------------------"
echo "11. Delete Coupon (ID: $PROD_COUPON_ID)"
curl -s -X DELETE $URL/coupons/$PROD_COUPON_ID | jq .


echo "--------------------------------------------------------"
echo "12. Get All Coupons After Deletion"
curl -s -X GET $URL/coupons | jq .
