package db

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Cotacao struct {
	ID  int    `json:"id"`
	Bid string `json:"bid"`
}

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		panic("Erro tentar conectar com o banco de dados.")
	}

	// Cria a tabela de cotações
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT NOT NULL
	);
	`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}
}

func AdicionarCotacao(ctx context.Context, bid string) error {
	query := "INSERT INTO cotacoes(bid) VALUES(?)"

	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bid)
	if err != nil {
		return err
	}

	return err

}

func BuscarCotacao() ([]Cotacao, error) {
	query := "SELECT * FROM cotacoes"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cotacoes []Cotacao

	for rows.Next() {
		var cotacao Cotacao
		err := rows.Scan(&cotacao.ID, &cotacao.Bid)
		if err != nil {
			return nil, err
		}
		cotacoes = append(cotacoes, cotacao)
	}

	return cotacoes, nil
}
