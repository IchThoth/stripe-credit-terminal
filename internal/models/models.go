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

// GopherImages is our product structure
type GopherImages struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Image          string    `json:"image"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

type Order struct {
	ID            int       `json:"id"`
	WidgetID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID string    `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

	row := m.DB.QueryRowContext(ctx, `select 
		id, name, description, inventory_level, price, coalesce(image, ''), 
		created_at, updated_at
	 from 
	 	GopherImages 
	 where id = ?`, id)
	err := row.Scan(
		&image.Id,
		&image.Name,
		&image.Description,
		&image.InventoryLevel,
		&image.Price,
		&image.Image,
		&image.CreatedAt,
		&image.UpdatedAt)
	if err != nil {
		return image, err
	}

	return image, nil
}
