package service

import (
	"monk-commerce/models"
	"monk-commerce/repository"
)

type CouponService struct {
	Repo repository.CouponRepo
}

func NewCouponService(repo repository.CouponRepo) *CouponService {
	return &CouponService{Repo: repo}
}

func (s *CouponService) GetApplicableCoupons(cart models.Cart) ([]models.ApplicableCoupon, error) {
	coupons, err := s.Repo.GetAll()
	if err != nil {
		return nil, err
	}

	var applicable []models.ApplicableCoupon
	for _, c := range coupons {
		uc := s.CalculateCartWithCoupon(cart, c)
		if uc.TotalDiscount > 0 {
			applicable = append(applicable, models.ApplicableCoupon{
				CouponID: c.ID,
				Type:     c.Type,
				Discount: uc.TotalDiscount,
			})
		}
	}

	return applicable, nil
}

func (s *CouponService) ApplyCoupon(cart models.Cart, couponID int) (*models.UpdatedCart, error) {
	coupon, err := s.Repo.GetByID(couponID)
	if err != nil {
		return nil, err
	}

	uc := s.CalculateCartWithCoupon(cart, coupon)
	if uc.TotalDiscount == 0 {
		return nil, repository.ErrCouponNotFound
	}
	return &uc, nil
}

func (s *CouponService) CalculateCartWithCoupon(cart models.Cart, coupon models.Coupon) models.UpdatedCart {
	var uc models.UpdatedCart

	for _, item := range cart.Items {
		uc.Items = append(uc.Items, models.UpdatedCartItem{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			Price:         item.Price,
			TotalDiscount: 0,
		})
	}

	for _, item := range uc.Items {
		uc.TotalPrice += float64(item.Quantity) * item.Price
	}
	switch coupon.Type {
	case models.CartWise:
		if coupon.Details.Threshold != nil && uc.TotalPrice > *coupon.Details.Threshold {
			if coupon.Details.Discount != nil {
				discountValue := (uc.TotalPrice * *coupon.Details.Discount) / 100.0
				uc.TotalDiscount = discountValue
			}
		}

	case models.ProductWise:
		if coupon.Details.ProductID != nil && coupon.Details.Discount != nil {
			pid := *coupon.Details.ProductID
			for i, item := range uc.Items {
				if item.ProductID == pid {
					itemDiscount := (float64(item.Quantity) * item.Price * *coupon.Details.Discount) / 100.0
					uc.Items[i].TotalDiscount = itemDiscount
					uc.TotalDiscount += itemDiscount
				}
			}
		}

	case models.BxGy:
		repLimit := 1
		if coupon.Details.RepLimit != nil {
			repLimit = *coupon.Details.RepLimit
		}
		cartQty := make(map[int]int)
		for _, item := range cart.Items {
			cartQty[item.ProductID] = item.Quantity
		}

		buyNeededCount := 0

		for _, b := range coupon.Details.BuyProducts {
			buyNeededCount += b.Quantity
		}

		eligibleBuyItemsCartQty := 0
		buySetIds := make(map[int]bool)
		for _, b := range coupon.Details.BuyProducts {
			buySetIds[b.ProductID] = true
			eligibleBuyItemsCartQty += cartQty[b.ProductID]
		}
		if buyNeededCount > 0 {
			possibleApplications := eligibleBuyItemsCartQty / buyNeededCount
			if possibleApplications > repLimit {
				possibleApplications = repLimit
			}
			if possibleApplications > 0 {
				for _, g := range coupon.Details.GetProducts {
					freeQty := g.Quantity * possibleApplications

					found := false
					for i, item := range uc.Items {
						if item.ProductID == g.ProductID {
							uc.Items[i].Quantity += freeQty
							addedValue := float64(freeQty) * item.Price
							uc.Items[i].TotalDiscount += addedValue
							uc.TotalPrice += addedValue
							uc.TotalDiscount += addedValue
							found = true
							break
						}
					}

					if !found {
					}
				}
			}
		}
	}

	uc.FinalPrice = uc.TotalPrice - uc.TotalDiscount
	if uc.FinalPrice < 0 {
		uc.FinalPrice = 0
	}

	return uc
}
