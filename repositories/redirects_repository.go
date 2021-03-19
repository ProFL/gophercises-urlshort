package repositories

import (
	"database/sql"
	"log"

	"github.com/ProFL/gophercises-urlshort/helpers"
	"github.com/ProFL/gophercises-urlshort/models"
)

type RedirectRepository struct {
	Db *sql.DB
}

func (repo *RedirectRepository) CreateTableIfNotExists() {
	_, err := repo.Db.Exec(
		`CREATE TABLE IF NOT EXISTS redirects(
		   path TEXT NOT NULL PRIMARY KEY,
		   url TEXT NOT NULL
		)`,
	)
	if err != nil {
		log.Panic("Failed to create redirects table")
	}
}

func (repo *RedirectRepository) Seed() {
	transaction, err := repo.Db.Begin()
	if err != nil {
		log.Println("Failed to create seed redirects transaction", err.Error())
		return
	}
	defer func() {
		commitErr := transaction.Commit()
		helpers.LogError("Failed to commit seed redirects transaction", commitErr)
	}()
	prepStmt, err := transaction.Prepare(`INSERT INTO redirects(path, url) VALUES(?, ?) ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Println("Failed to prepare seed redirects transaction statement", err.Error())
		return
	}
	defer func() {
		closeErr := prepStmt.Close()
		helpers.LogError("Failed to close redirects seed prepared statement", closeErr)
	}()

	_, err = prepStmt.Exec("/google", "https://google.com.br/")
	helpers.LogError("Failed to insert /google redirects seed", err)
	_, err = prepStmt.Exec("/golang", "https://golang.org/")
	helpers.LogError("Failed to insert /golang redirects seed", err)
}

func (repo *RedirectRepository) FindByPath(path string) (*models.Redirect, error) {
	rows, err := repo.Db.Query(`SELECT path, url FROM redirects WHERE path = ?`, path)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var path string
		var url string
		rows.Scan(&path, &url)
		return &models.Redirect{Path: path, Url: url}, nil
	}
	return nil, nil
}
