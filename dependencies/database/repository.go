package database

import (
	"context"
	"generate-prompt-service/app/domain/agentnumber"
	"os"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	Conn *pgx.Conn
}

func NewConnection() (*pgx.Conn, error) {
    dbURL := os.Getenv("DATABASE_URL")
    return pgx.Connect(context.Background(), dbURL)
}

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{Conn: conn}
}

func (r *Repository) InsertAgentNumber(ctx context.Context, an agentnumber.AgentNumber) error {
	query := `
		INSERT INTO agent_number (client_id, whatsapp_number, description, customer_id)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.Conn.Exec(ctx, query, an.ClientID, an.WhatsappNumber, an.Description, an.CustomerID)
	return err
}