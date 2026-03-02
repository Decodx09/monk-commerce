package models

type CartPayload struct {
	Cart Cart `json:"cart"`
}

type Cart struct {
	Items []CartItem `json:"items"`
}

type CartItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type ApplicableCouponsResponse struct {
	ApplicableCoupons []ApplicableCoupon `json:"applicable_coupons"`
}

type ApplicableCoupon struct {
	CouponID int        `json:"coupon_id"`
	Type     CouponType `json:"type"`
	Discount float64    `json:"discount"`
}

type ApplyCouponResponse struct {
	UpdatedCart UpdatedCart `json:"updated_cart"`
}

type UpdatedCart struct {
	Items         []UpdatedCartItem `json:"items"`
	TotalPrice    float64           `json:"total_price"`
	TotalDiscount float64           `json:"total_discount"`
	FinalPrice    float64           `json:"final_price"`
}

type UpdatedCartItem struct {
	ProductID     int     `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
	TotalDiscount float64 `json:"total_discount"`
}
