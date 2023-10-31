package models

import (
	"context"
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
	InventoryLevel int       `json:"inventory_level"`
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

func (m *DBmodels) GetGopherImages(id int) (GopherImages, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	var image GopherImages

	row := m.DB.QueryRowContext(ctx, "select id, name from GopherImages where id = ?", id)
	err := row.Scan(&image.Id, &image.Name)
	if err != nil {
		return image, err
	}

	return image, nil
}
