package controllers

import (
	"strconv"

	"monk-commerce/models"
	"monk-commerce/repository"
	"monk-commerce/service"

	"github.com/gofiber/fiber/v2"
)

type CouponController struct {
	Repo    repository.CouponRepo
	Service *service.CouponService
}

func NewCouponController(repo repository.CouponRepo, service *service.CouponService) *CouponController {
	return &CouponController{
		Repo:    repo,
		Service: service,
	}
}

func (c *CouponController) CreateCoupon(ctx *fiber.Ctx) error {
	var coupon models.Coupon
	if err := ctx.BodyParser(&coupon); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	if err := c.Repo.Create(&coupon); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(coupon)
}

func (c *CouponController) GetAllCoupons(ctx *fiber.Ctx) error {
	coupons, err := c.Repo.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(coupons)
}

func (c *CouponController) GetCouponByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id format"})
	}

	coupon, err := c.Repo.GetByID(id)
	if err != nil {
		if err == repository.ErrCouponNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(coupon)
}

func (c *CouponController) UpdateCoupon(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id format"})
	}

	var coupon models.Coupon
	if err := ctx.BodyParser(&coupon); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	if err := c.Repo.Update(id, &coupon); err != nil {
		if err == repository.ErrCouponNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(coupon)
}

func (c *CouponController) DeleteCoupon(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id format"})
	}

	if err := c.Repo.Delete(id); err != nil {
		if err == repository.ErrCouponNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *CouponController) GetApplicableCoupons(ctx *fiber.Ctx) error {
	var payload models.CartPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	applicable, err := c.Service.GetApplicableCoupons(payload.Cart)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(models.ApplicableCouponsResponse{ApplicableCoupons: applicable})
}

func (c *CouponController) ApplyCoupon(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id format"})
	}

	var payload models.CartPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	updatedCart, err := c.Service.ApplyCoupon(payload.Cart, id)
	if err != nil {
		if err == repository.ErrCouponNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(models.ApplyCouponResponse{UpdatedCart: *updatedCart})
}
