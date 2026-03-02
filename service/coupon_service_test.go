package service

import (
	"testing"

	"monk-commerce/models"
)

func floatPtr(v float64) *float64 { return &v }
func intPtr(v int) *int           { return &v }

func TestCalculateCartWithCoupon_CartWise(t *testing.T) {
	svc := &CouponService{}

	cart := models.Cart{
		Items: []models.CartItem{
			{ProductID: 1, Quantity: 2, Price: 50.0},
		},
	}

	coupon := models.Coupon{
		Type: models.CartWise,
		Details: models.CouponDetails{
			Threshold: floatPtr(50),
			Discount:  floatPtr(10),
		},
	}

	uc := svc.CalculateCartWithCoupon(cart, coupon)
	if uc.TotalPrice != 100 {
		t.Errorf("expected 100, got %f", uc.TotalPrice)
	}
	if uc.TotalDiscount != 10 {
		t.Errorf("expected 10, got %f", uc.TotalDiscount)
	}
	if uc.FinalPrice != 90 {
		t.Errorf("expected 90, got %f", uc.FinalPrice)
	}
}

func TestCalculateCartWithCoupon_BxGy(t *testing.T) {
	svc := &CouponService{}

	cart := models.Cart{
		Items: []models.CartItem{
			{ProductID: 1, Quantity: 6, Price: 50.0},
			{ProductID: 2, Quantity: 3, Price: 30.0},
			{ProductID: 3, Quantity: 2, Price: 25.0},
		},
	}
	coupon := models.Coupon{
		Type: models.BxGy,
		Details: models.CouponDetails{
			BuyProducts: []models.ProductCondition{
				{ProductID: 1, Quantity: 3},
				{ProductID: 2, Quantity: 3},
			},
			GetProducts: []models.ProductCondition{
				{ProductID: 3, Quantity: 1},
			},
			RepLimit: intPtr(2),
		},
	}
	uc := svc.CalculateCartWithCoupon(cart, coupon)
	if uc.FinalPrice < 0 {
		t.Fail()
	}
}
