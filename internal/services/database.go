package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jhenriquem/gongo-wabot/internal/types"
	_ "github.com/lib/pq"
)

func InitDatabase() {
	var err error

	SVC.Database, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Database online")

	// Criar a tabela caso não exista
	_, err = SVC.Database.Exec("CREATE TABLE IF NOT EXISTS books (id SERIAL PRIMARY KEY, title TEXT, author TEXT , value INT, amount INT)")
	if err != nil {
		log.Fatal(err)
	}
}

func GetBook(book *types.Book, msg string) error {
	// Inicia uma transação
	tx, err := SVC.Database.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	// Busca o livro no banco de dados
	err = tx.QueryRow("SELECT * FROM books WHERE LOWER(title) LIKE LOWER($1) || '%' ", msg).Scan(&book.ID, &book.Title, &book.Author, &book.Value, &book.Amount)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao buscar livro: %w", err)
	}

	// Confirma a transação
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}

// Atualiza a quantidade de um livro no estoque
func SetBookStock(bookId string, amount int) error {
	var bookName string
	var oldAmount int

	// Inicia uma transação
	tx, err := SVC.Database.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	err = tx.QueryRow("SELECT title, amount FROM books WHERE id = $1", bookId).Scan(&bookName, &oldAmount)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao buscar livro: %w", err)
	}

	_, err = tx.Exec("UPDATE books SET amount = $1 WHERE id = $2", oldAmount-amount, bookId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao atualizar estoque: %w", err)
	}

	// Confirma a transação
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}
