package entity

import "time"

type MerchantEntity struct {
	ID              int64      `db:"id"`
	Name            string     `db:"name"`
	ServerKey       *string    `db:"server_key"`
	KeyEnvirontment string     `db:"key_environment"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}
