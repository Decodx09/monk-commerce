package models

import "time"

type CouponType string

const (
	CartWise    CouponType = "cart-wise"
	ProductWise CouponType = "product-wise"
	BxGy        CouponType = "bxgy"
)

type ProductCondition struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CouponDetails struct {
	Threshold *float64 `json:"threshold,omitempty"`
	Discount  *float64 `json:"discount,omitempty"`

	ProductID *int `json:"product_id,omitempty"`

	BuyProducts []ProductCondition `json:"buy_products,omitempty"`
	GetProducts []ProductCondition `json:"get_products,omitempty"`
	RepLimit    *int               `json:"repition_limit,omitempty"`
}

type Coupon struct {
	ID             int           `json:"coupon_id"`
	Type           CouponType    `json:"type"`
	Details        CouponDetails `json:"details"`
	ExpirationDate *time.Time    `json:"expiration_date,omitempty"`
}
