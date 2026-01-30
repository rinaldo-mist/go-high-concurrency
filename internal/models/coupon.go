package models

type Coupon struct {
	Name    string `json:"name"`
	Amount  int    `json:"amount"`
	Version int64  `json:"version"`
}
