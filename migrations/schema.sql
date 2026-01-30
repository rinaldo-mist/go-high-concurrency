CREATE database dev;

CREATE TABLE IF NOT EXISTS coupons (
  coupon_name VARCHAR(255) UNIQUE,
  balance INT NOT NULL DEFAULT 0,
  version INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS claim_history (
  user_id VARCHAR(255),
  coupon_name VARCHAR(255)
);