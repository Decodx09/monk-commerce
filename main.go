package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "github.com/go-sql-driver/mysql"

	"monk-commerce/controllers"
	"monk-commerce/repository"
	"monk-commerce/service"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "monk_commerce"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to MySQL Database!")
	repo := repository.NewMySQLCouponRepo(db)

	svc := service.NewCouponService(repo)
	ctrl := controllers.NewCouponController(repo, svc)

	app.Post("/coupons", ctrl.CreateCoupon)
	app.Get("/coupons", ctrl.GetAllCoupons)
	app.Get("/coupons/:id", ctrl.GetCouponByID)
	app.Put("/coupons/:id", ctrl.UpdateCoupon)
	app.Delete("/coupons/:id", ctrl.DeleteCoupon)

	app.Post("/applicable-coupons", ctrl.GetApplicableCoupons)
	app.Post("/apply-coupon/:id", ctrl.ApplyCoupon)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
