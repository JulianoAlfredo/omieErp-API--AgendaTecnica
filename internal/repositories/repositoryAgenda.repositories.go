package repositories

import (
	"database/sql"
	"log"
)

func SearchClients(db *sql.DB, idClient string) []map[string]any {
	rows, err := db.Query("SELECT id, nome_fantasia, razao_social, emails, cnpj FROM amm_clientes WHERE id = ? ORDER BY id DESC", idClient)
	if err != nil {
		log.Fatal("Error querying database: ", err.Error())
	}
	defer rows.Close()
	employees := []map[string]any{}
	for rows.Next() {
		var id int
		var nomeFantasia string
		var cnpj string
		var razao_social string
		var emails string
		err := rows.Scan(&id, &nomeFantasia, &razao_social, &emails, &cnpj)
		if err != nil {
			log.Fatal("Error scanning row: ", err.Error())
		}
		employees = append(employees, map[string]any{
			"id":            id,
			"nome_fantasia": nomeFantasia,
			"razao_social":  razao_social,
			"emails":        emails,
			"cnpj":          cnpj,
		})

	}
	return employees
}
