package repositories

import (
	"database/sql"
	"log"
)

func SearchClients(db *sql.DB) {
	rows, err := db.Query("SELECT TOP 10 id, nome_fantasia FROM amm_clientes ORDER BY id DESC")
	if err != nil {
		log.Fatal("Error querying database: ", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var nomeFantasia string
		err := rows.Scan(&id, &nomeFantasia)
		if err != nil {
			log.Fatal("Error scanning row: ", err.Error())
		}
		log.Printf("ID: %d, Nome Fantasia: %s\n", id, nomeFantasia)
	}

}
