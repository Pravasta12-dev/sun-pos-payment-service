package entity

import "time"

type MerchantEntity struct {
	ID              string      `db:"id"`
	Name            string     `db:"name"`
	ServerKey       *string    `db:"server_key"`
	KeyEnvironment  string     `db:"key_environment"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}

func (m *MerchantEntity) TableName() string {
	return "merchants"
}