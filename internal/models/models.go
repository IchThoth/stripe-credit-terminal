package models

import (
	"database/sql"
	"time"
)

type DBmodels struct {
	DB *sql.DB
}

type Models struct {
	DB DBmodels
}

//GopherImages is our product structure
type GopherImages struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel string    `json:"inventory_level"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBmodels{
			DB: db,
		},
	}
}
