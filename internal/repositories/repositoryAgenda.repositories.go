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

func WebhookUpdateOsIncluida(db *sql.DB, idOs string, CodigoIntegra string) (sql.Result, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM amm_contas_omie_x_agenda WHERE id_conta_agenda = ?", CodigoIntegra).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar se o registro existe: %v", err)
		return nil, err
	}
	if count == 0 {
		insertDb, err := db.Exec("INSERT INTO amm_contas_omie_x_agenda (id_conta_agenda, id_os) VALUES (?, ?)", CodigoIntegra, idOs)
		if err != nil {
			log.Printf("Erro ao inserir novo registro: %v", err)
			return nil, err
		} else {
			rowsAffected, _ := insertDb.RowsAffected()
			log.Printf("Novo registro inserido com sucesso. Linhas afetadas: %d", rowsAffected)
			return insertDb, nil
		}
	}
	result, err := db.Exec("UPDATE amm_contas_omie_x_agenda SET id_os = ? WHERE id_conta_agenda = ?", idOs, CodigoIntegra)
	if err != nil {
		log.Printf("Erro ao atualizar o banco de dados: %v", err)
		return nil, err
	}
	return result, nil
}

func WebhookUpdateOsFaturada(db *sql.DB, idOs string, CodigoIntegra string) (sql.Result, error) {
	result, err := db.Exec("UPDATE amm_contas_omie_x_agenda SET faturada = 1 WHERE  id_os = ?", idOs)
	if err != nil {
		log.Printf("Erro ao atualizar o banco de dados: %v", err)
		return nil, err
	}
	return result, err
}
