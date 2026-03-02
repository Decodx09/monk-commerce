package repository

import (
	"errors"

	"monk-commerce/models"
)

var (
	ErrCouponNotFound = errors.New("coupon not found")
)

type CouponRepo interface {
	Create(coupon *models.Coupon) error
	GetAll() ([]models.Coupon, error)
	GetByID(id int) (models.Coupon, error)
	Update(id int, coupon *models.Coupon) error
	Delete(id int) error
}
