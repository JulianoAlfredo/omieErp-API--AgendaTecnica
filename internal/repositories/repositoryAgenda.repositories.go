package repositories

import (
	"database/sql"
	"encoding/json"
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
func WebhookUpdateOsIncluida(db *sql.DB, idOs string, CodigoIntegra string, NumeroOs string) (sql.Result, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM amm_contas_omie_x_agenda WHERE id_conta_agenda = ?", CodigoIntegra).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar se o registro existe: %v", err)
		return nil, err
	}
	if count == 0 {
		insertDb, err := db.Exec("INSERT INTO amm_contas_omie_x_agenda (id_conta_agenda, id_os, numero_os) VALUES (?, ?, ?)", CodigoIntegra, idOs, NumeroOs)
		if err != nil {
			log.Printf("Erro ao inserir novo registro: %v", err)
			return nil, err
		} else {
			rowsAffected, _ := insertDb.RowsAffected()
			log.Printf("Novo registro inserido com sucesso. Linhas afetadas: %d", rowsAffected)
			return insertDb, nil
		}
	}
	result, err := db.Exec("UPDATE amm_contas_omie_x_agenda SET id_os = ?, numero_os = ? WHERE id_conta_agenda = ?", idOs, NumeroOs, CodigoIntegra)
	if err != nil {
		log.Printf("Erro ao atualizar o banco de dados: %v", err)
		return nil, err
	}
	return result, nil
}

func CriarTabelaRelacaoClientes(db *sql.DB) error {
	query := `
	
	CREATE TABLE IF NOT EXISTS amm_omie_relaciona_clientes (
		cliente_agenda NVARCHAR(200),
		cliente_omie   BIGINT,
		cnpj           NVARCHAR(30)
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Erro ao criar tabela amm_omie_relaciona_clientes: %v", err)
		return err
	}
	return nil
}

func UpsertRelacaoCliente(db *sql.DB, clienteAgenda string, clienteOmie int64, cnpj string) error {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM amm_omie_relaciona_clientes WHERE cliente_agenda = ?", clienteAgenda).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar relacao cliente %s: %v", clienteAgenda, err)
		return err
	}
	if count == 0 {
		_, err = db.Exec(
			"INSERT INTO amm_omie_relaciona_clientes (cliente_agenda, cliente_omie, cnpj) VALUES (?, ?, ?)",
			clienteAgenda, clienteOmie, cnpj,
		)
		if err != nil {
			log.Printf("Erro ao inserir relacao cliente %s: %v", clienteAgenda, err)
			return err
		}
	} else {
		_, err = db.Exec(
			"UPDATE amm_omie_relaciona_clientes SET cliente_omie = ?, cnpj = ? WHERE cliente_agenda = ?",
			clienteOmie, cnpj, clienteAgenda,
		)
		if err != nil {
			log.Printf("Erro ao atualizar relacao cliente %s: %v", clienteAgenda, err)
			return err
		}
	}
	return nil
}
func WebhookUpdateOsFaturada(db *sql.DB, idOs string, CodigoIntegra string) (sql.Result, error) {
	result, err := db.Exec("UPDATE amm_contas_omie_x_agenda SET faturada = 1 WHERE  id_os = ?", idOs)
	if err != nil {
		log.Printf("Erro ao atualizar o banco de dados: %v", err)
		return nil, err
	}
	return result, err
}
func WebhookInsertContaReceber(db *sql.DB, CodigoCliente int64, CodigoConta int64, NumeroDocumento string, NumeroDocumentoFiscal string, NumeroPedido string) (sql.Result, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM amm_contas_omie_x_agenda WHERE  numero_os = ?", NumeroPedido).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar se o registro existe: %v", err)
		return nil, err
	}
	var insertDb sql.Result
	if count == 0 {
		insertDb, err = db.Exec("INSERT INTO amm_contas_omie_x_agenda (id_cliente, id_conta_omie, numero_nf, numero_rps, numero_os) VALUES (?, ?, ?, ?, ?)", CodigoCliente, CodigoConta, NumeroDocumentoFiscal, NumeroDocumento, NumeroPedido)
	} else {
		insertDb, err = db.Exec("UPDATE amm_contas_omie_x_agenda SET id_cliente = ?, id_conta_omie = ?, numero_nf = ?, numero_rps = ?, numero_os = ? WHERE  numero_os = ?", CodigoCliente, CodigoConta, NumeroDocumentoFiscal, NumeroDocumento, NumeroPedido, NumeroPedido)
	}
	if err != nil {
		log.Printf("Erro ao inserir conta a receber: %v", err)
		return nil, err
	}
	return insertDb, nil
}
func WebhookInsertBoletoGerado(db *sql.DB, CodigoCliente int64, CodigoConta int64, NumeroPedido string, BoletoGerado string, CodigoBarras string, BoletoNumero string) (sql.Result, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM amm_contas_omie_x_agenda WHERE  numero_os = ?", NumeroPedido).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar se o registro existe: %v", err)
		return nil, err
	}

	var insertDb sql.Result
	if count == 0 {
		insertDb, err = db.Exec("INSERT INTO amm_contas_omie_x_agenda (id_cliente, id_conta_omie, numero_os, boleto_gerado, codigo_barras) VALUES (?, ?, ?, ?, ?)", CodigoCliente, CodigoConta, NumeroPedido, BoletoGerado, CodigoBarras)
	} else {
		insertDb, err = db.Exec("UPDATE amm_contas_omie_x_agenda SET  id_conta_omie = ?, boleto_gerado = ?, codigo_barras_boleto = ? , boleto_numero = ? WHERE  numero_os = ? AND id_cliente = ?", CodigoConta, BoletoGerado, CodigoBarras, BoletoNumero, NumeroPedido, CodigoCliente)
	}
	if err != nil {
		log.Printf("Erro ao inserir boleto gerado: %v", err)
		return nil, err
	}
	return insertDb, nil
}

func InsertLinkBoletoGerado(db *sql.DB, CodigoConta int64, LinkBoletoGerado string) (sql.Result, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM amm_contas_omie_x_agenda WHERE  id_conta_omie = ?", CodigoConta).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar se o registro existe: %v", err)
		return nil, err
	}
	var insertDb sql.Result
	if count == 0 {
		insertDb, err = db.Exec("INSERT INTO amm_contas_omie_x_agenda (link_boleto) VALUES (?)", LinkBoletoGerado)
	} else {
		insertDb, err = db.Exec("UPDATE amm_contas_omie_x_agenda SET link_boleto = ? WHERE id_conta_omie = ?", LinkBoletoGerado, CodigoConta)
	}
	if err != nil {
		log.Printf("Erro ao inserir link do boleto gerado: %v", err)
		return nil, err
	}
	return insertDb, nil
}
func InserirLogFaturamento(db *sql.DB, codIntOS string, etapa string, status string, mensagem string, dados any) error {
	var dadosJSON *string
	if dados != nil {
		b, err := json.Marshal(dados)
		if err == nil {
			s := string(b)
			dadosJSON = &s
		}
	}
	_, err := db.Exec(
		`INSERT INTO amm_omie_faturamento_log (cod_int_os, etapa, status, mensagem, dados) VALUES (?, ?, ?, ?, ?)`,
		codIntOS, etapa, status, mensagem, dadosJSON,
	)
	if err != nil {
		log.Printf("[FaturamentoLog] Erro ao inserir log (etapa=%s): %v", etapa, err)
	}
	return err
}
