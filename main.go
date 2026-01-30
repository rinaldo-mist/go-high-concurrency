package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	_ "github.com/lib/pq"

	cfg "highconcurrency/internal/config"
	"highconcurrency/internal/handlers"
)

var db *sql.DB

func main() {
	fmt.Println("starting app .....")
	var err error
	db, err = cfg.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Recover())

	// Coupon routes
	e.POST("/api/coupons", func(c *echo.Context) error {
		return handlers.CreateCoupon(c, db)
	})
	e.POST("/api/coupons/claim", func(c *echo.Context) error {
		return handlers.UpdateWithOptimisticLock(c, db)
	})
	e.GET("/api/coupons/:name", func(c *echo.Context) error {
		return handlers.GetCoupons(c, db)
	})

	fmt.Println("app started !")

	log.Println("server running on :8080")
	e.Start(":8080")
}
