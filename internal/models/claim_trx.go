package models

import "time"

type ClaimTrx struct {
	UserID      string `json:"user_id"`
	CouponName  string `json:"coupon_name"`
	CreatedDate time.Time
}
