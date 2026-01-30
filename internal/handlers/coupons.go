package handlers

import (
	"database/sql"
	"fmt"
	"highconcurrency/internal/models"

	"net/http"

	"github.com/labstack/echo/v5"
)

func GetCoupons(c *echo.Context, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	couponName := c.Param("name")

	query := fmt.Sprintf(`
		SELECT c.coupon_name, balance, user_id
		FROM coupons c left join claim_history ch
		ON c.coupon_name = ch.coupon_name
		WHERE c.coupon_name = '%s'
	`, couponName)

	rows, err := tx.Query(query)

	if err != nil {
		return err
	}
	defer rows.Close()

	var coupons []map[string]interface{}

	var name string
	var amount int
	var claimedByArr []string

	for rows.Next() {
		var claimedBy sql.NullString

		if err := rows.Scan(&name, &amount, &claimedBy); err != nil {
			fmt.Println("err:", err)
			return err
		}

		if !claimedBy.Valid {
			continue
		}

		claimedByArr = append(claimedByArr, claimedBy.String)
	}

	if len(name) == 0 {
		return c.JSON(http.StatusNotFound, "coupon not found")
	}

	coupons = append(coupons, map[string]interface{}{
		"name":             name,
		"amount":           amount + len(claimedByArr),
		"remaining_amount": amount,
		"claimed_by":       claimedByArr,
	})
	tx.Commit()
	return c.JSON(http.StatusOK, coupons)
}

func CreateCoupon(c *echo.Context, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	fmt.Println("CreateCoupon called")
	var req *models.Coupon
	if err := c.Bind(&req); err != nil {
		fmt.Println("CreateCoupon err binding", err)
		return err
	}

	fmt.Println("CreateCoupon binding success")

	query := fmt.Sprintf(`INSERT INTO coupons (coupon_name, balance, version) 
		VALUES ('%s', %d, 1)`, req.Name, req.Amount)

	_, err = tx.Exec(query)
	if err != nil {
		fmt.Println("CreateCoupon err insert", err)
		return err
	}
	fmt.Println("CreateCoupon success insert")
	tx.Commit()

	fmt.Println("CreateCoupon committed")

	return c.JSON(http.StatusCreated, "coupon created")
}

// ✅ Optimistic locking with version column
func UpdateWithOptimisticLock(c *echo.Context, db *sql.DB) error {
	fmt.Println("UpdateWithOptimisticLock called")
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var req *models.ClaimTrx

	if err := c.Bind(&req); err != nil {
		fmt.Println("UpdateWithOptimisticLock binding error:", err)
		return err
	}
	fmt.Println("UpdateWithOptimisticLock success binding")

	for retries := 0; retries < 3; retries++ {
		var balance int
		var version int
		query := fmt.Sprintf(`SELECT balance, version
            FROM coupons c inner join claim_history ch
			ON c.coupon_name = ch.coupon_name
            WHERE c.coupon_name = '%s' AND ch.user_id = '%s'`, req.CouponName, req.UserID)

		err := tx.QueryRow(query).Scan(&balance, &version)

		if err != sql.ErrNoRows {
			fmt.Println("UpdateWithOptimisticLock - err already claimed :", err)
			return c.JSON(http.StatusNotFound, fmt.Sprintf("%s already claimed %s", req.UserID, req.CouponName))
		}

		var couponName string
		var claimedCouponCount int

		query = fmt.Sprintf(`SELECT c.coupon_name, c.balance, c.version, count(ch.*) coupon_claimed
            FROM coupons c left join claim_history ch 
            on c.coupon_name = ch.coupon_name
            WHERE c.coupon_name = '%s'
            group by 1,2,3`, req.CouponName)

		err = tx.QueryRow(query).Scan(&couponName, &balance, &version, &claimedCouponCount)

		if err == sql.ErrNoRows {
			fmt.Println("UpdateWithOptimisticLock - err coupon not found :", err)
			return c.JSON(http.StatusNotFound, "coupon not found")
		}

		if balance <= 0 {
			fmt.Println("UpdateWithOptimisticLock - err insufficient balance coupon :", err)
			return c.JSON(http.StatusBadRequest, "insufficient balance")
		}

		fmt.Printf("Coupon name : %s , Update version: %d\n", req.CouponName, version)

		balance -= 1

		query = fmt.Sprintf(`UPDATE coupons
            SET balance = %d, version = version + 1
            WHERE coupon_name = '%s' AND version = %d`, balance, req.CouponName, version)
		res, err := tx.Exec(query)

		if err != nil {
			fmt.Println("UpdateWithOptimisticLock - err update:", err)
			return err
		}
		fmt.Printf("Coupon name : %s , Update version: %d\n", req.CouponName, version+1)

		rows, err := res.RowsAffected()

		fmt.Println("Current version:", version, "Rows affected:", rows)

		errClaim := CreateClaimHistory(tx, req)
		if errClaim != nil {
			return errClaim
		}
		fmt.Println("Rows affected:", rows)

		if rows == 1 {
			fmt.Println("Sukses ! Current version:", version, "Rows affected:", rows)
			tx.Commit()
			return c.JSON(http.StatusOK, "claimed (success)")
		}
		// else: conflict detected → retry
	}
	return c.JSON(http.StatusConflict, "high contention, retry")
}

func CreateClaimHistory(tx *sql.Tx, req *models.ClaimTrx) error {
	fmt.Println("Req:", req)

	query := fmt.Sprintf(`INSERT INTO claim_history (user_id, coupon_name)
		VALUES ('%s', '%s')`, req.UserID, req.CouponName)

	_, err := tx.Exec(query)

	if err != nil {
		fmt.Println("err claim:", err)
		return err
	}

	return nil
}
