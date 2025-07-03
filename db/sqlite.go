package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./data.db") // Cria ou abre o arquivo do banco de dados
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS dados (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		usuario_que_registrou VARCHAR(255) NOT NULL,
		data_registrada DATE NOT NULL,
		descricao_curta VARCHAR(255) NOT NULL,
		descricao_longa TEXT,
		data_de_expiracao DATE
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Tabela 'dados' verificada/criada com sucesso!")
}