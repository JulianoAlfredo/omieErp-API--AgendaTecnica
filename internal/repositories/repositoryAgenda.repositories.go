package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
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

func UpsertRelacaoCliente(db *sql.DB, clienteAgenda string, clienteOmie int64, cnpj string) error {
	var count int
	fmt.Printf(clienteAgenda, clienteOmie, cnpj)
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

func UpsertNFSEGerada(db *sql.DB, nCodNF int64, codigoOs float64, cDataEmissao string, cXmlNFSe string, cUrlNFSe string, cLinkPortal string, cNumNFSe string, cPdfNfse string) (sql.Result, error) {
	var count int
	fmt.Printf("%f\n", codigoOs)
	err := db.QueryRow("SELECT COUNT(1) FROM amm_contas_omie_x_agenda WHERE id_os = ?", codigoOs).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar se o registro existe: %v", err)
		return nil, err
	}
	var result sql.Result
	if count == 0 {
		result, err = db.Exec(
			"INSERT INTO amm_contas_omie_x_agenda (id_nf, id_os, data_emissao, xml_nfe, link_portal, numero_nf, link_nf) VALUES (?, ?, ?, ?, ?, ?, ?)",
			nCodNF, codigoOs, cDataEmissao, cXmlNFSe, cLinkPortal, cNumNFSe, cPdfNfse,
		)
		fmt.Printf("inserindo novo registro")
	} else {
		result, err = db.Exec(
			"UPDATE amm_contas_omie_x_agenda SET id_nf = ?, data_emissao = ?, xml_nfe = ?, link_portal = ?, numero_nf = ?, link_nf = ? WHERE id_os = ?",
			nCodNF, cDataEmissao, cXmlNFSe, cLinkPortal, cNumNFSe, cPdfNfse, codigoOs,
		)
		fmt.Printf("atualizando registro existente\n")

	}

	if err != nil {
		log.Printf("Erro ao inserir/atualizar NFSE gerada: %v", err)
		return nil, err
	}
	return result, nil
}

func WebhookUpdateConferido(db *sql.DB, codigoLancamentoOmie int64, data string, dataCred string, observacao string, valor float64) (sql.Result, error) {
	parsedData, err := time.Parse(time.RFC3339, data)
	if err != nil {
		log.Printf("Erro ao parsear data_baixa '%s': %v", data, err)
		return nil, err
	}
	parsedDataCred, err := time.Parse(time.RFC3339, dataCred)
	if err != nil {
		log.Printf("Erro ao parsear data_cred '%s': %v", dataCred, err)
		return nil, err
	}
	result, err := db.Exec(
		"UPDATE amm_contas_omie_x_agenda SET conferido = 1, data_baixa = ?, data_cred = ?, observacao_baixa = ?, valor_baixa = ? WHERE id_conta_omie = ?",
		parsedData.UTC(), parsedDataCred.UTC(), observacao, valor, codigoLancamentoOmie,
	)
	if err != nil {
		log.Printf("Erro ao atualizar conferido: %v", err)
		return nil, err
	}
	return result, nil
}

func UpdateBaixaPorNumeroRps(db *sql.DB, numeroRps string, codigoLancamentoOmie int64, valor float64, observacao string, dataBaixa string) (int64, error) {
	parsed, err := time.Parse("02/01/2006", dataBaixa)
	if err != nil {
		log.Printf("Erro ao parsear data_baixa '%s': %v", dataBaixa, err)
		return 0, err
	}
	result, err := db.Exec(
		`UPDATE amm_contas_omie_x_agenda
		 SET conferido = 1, id_conta_omie = ?, valor_baixa = ?, observacao_baixa = ?, data_baixa = ?, data_cred = ?
		 WHERE numero_rps = ?`,
		codigoLancamentoOmie, valor, observacao, parsed.UTC(), parsed.UTC(), numeroRps,
	)
	if err != nil {
		log.Printf("Erro ao atualizar baixa por numero_rps '%s': %v", numeroRps, err)
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}

func SearchClientByField(db *sql.DB, id int) (map[string]any, error) {
	const query = `SELECT
		c.ID AS codigo_integracao,
		LEFT(CAST(c.RAZAO_SOCIAL AS VARCHAR(MAX)), 60) AS RAZAO_SOCIAL,
		CASE
			WHEN NULLIF(REPLACE(REPLACE(REPLACE(LTRIM(RTRIM(CAST(c.CNPJ AS VARCHAR(MAX)))), '.', ''), '-', ''), '/', ''), '') IS NULL THEN NULL
			WHEN LEN(REPLACE(REPLACE(REPLACE(LTRIM(RTRIM(CAST(c.CNPJ AS VARCHAR(MAX)))), '.', ''), '-', ''), '/', '')) < 11 THEN NULL
			ELSE LEFT(REPLACE(REPLACE(REPLACE(LTRIM(RTRIM(CAST(c.CNPJ AS VARCHAR(MAX)))), '.', ''), '-', ''), '/', ''), 14)
		END AS CNPJ,
		LEFT(COALESCE(NULLIF(LTRIM(RTRIM(CAST(c.NOME_FANTASIA AS VARCHAR(MAX)))), ''), CAST(c.RAZAO_SOCIAL AS VARCHAR(MAX))), 100) AS nome_fantasia,
		'Sim' AS cliente,
		'Não' AS fornecedor,
		'Não' AS transportadora,
		'Não' AS funcionario,
		c.NOME_CONTATOS AS contato,
		LEFT(LTRIM(RTRIM(ISNULL(CAST(c.ENDERECO AS VARCHAR(MAX)), ''))), 60) AS ENDERECO,
		c.endereco_numero AS numero,
		c.Bairro,
		LEFT(LTRIM(RTRIM(ISNULL(CAST(c.endereco_Complemento AS VARCHAR(MAX)), ''))), 60) AS complemento,
		NULLIF(
			CASE UPPER(LTRIM(RTRIM(ISNULL(CAST(c.Estado AS VARCHAR(MAX)), ''))))
				WHEN 'NULL'   THEN ''
				WHEN 'NONE'   THEN ''
				WHEN 'ESTADO' THEN ''
				ELSE LTRIM(RTRIM(CAST(c.Estado AS VARCHAR(MAX))))
			END
		, '') AS Estado,
		CASE
			WHEN UPPER(LTRIM(RTRIM(ISNULL(CAST(c.Cidade AS VARCHAR(MAX)), '')))) IN ('NULL', 'NONE', '0', 'RJ', '') THEN NULL
			WHEN LTRIM(RTRIM(CAST(c.Cidade AS VARCHAR(MAX)))) = 'Rio Janeiro' THEN 'Rio de Janeiro'
			WHEN UPPER(LTRIM(RTRIM(CAST(c.Cidade AS VARCHAR(MAX))))) IN ('SÃO MATHEUS', 'SAO MATHEUS') THEN 'São Mateus'
			ELSE LTRIM(RTRIM(CAST(c.Cidade AS VARCHAR(MAX))))
		END AS Cidade,
		CASE
			WHEN UPPER(LTRIM(RTRIM(ISNULL(CAST(c.Pais AS VARCHAR(MAX)), '')))) IN ('NULL', 'NONE', '0', 'BR', '') THEN NULL
			ELSE LTRIM(RTRIM(CAST(c.Pais AS VARCHAR(MAX))))
		END AS Pais,
		NULLIF(LEFT(REPLACE(REPLACE(LTRIM(RTRIM(ISNULL(CAST(c.CEP AS VARCHAR(MAX)), ''))), '-', ''), '.', ''), 8), '') AS CEP,
		'' AS ddd_telefone,
		'' AS telefone,
		NULL AS ddd_telefone2,
		NULL AS telefone2,
		NULL AS ddd_fax,
		NULL AS fax,
		CASE
			WHEN LTRIM(RTRIM(ISNULL(CAST(c.EMAILS AS VARCHAR(MAX)), ''))) IN ('', '0', 'NULL') THEN NULL
			WHEN CHARINDEX('@', LTRIM(RTRIM(CAST(c.EMAILS AS VARCHAR(MAX))))) = 0 THEN NULL
			WHEN CHARINDEX('.', SUBSTRING(LTRIM(RTRIM(CAST(c.EMAILS AS VARCHAR(MAX)))), CHARINDEX('@', LTRIM(RTRIM(CAST(c.EMAILS AS VARCHAR(MAX))))) + 1, LEN(LTRIM(RTRIM(CAST(c.EMAILS AS VARCHAR(MAX))))))) = 0 THEN NULL
			ELSE LTRIM(RTRIM(CAST(c.EMAILS AS VARCHAR(MAX))))
		END AS EMAILS,
		c.siteCliente,
		NULL AS banco,
		NULL AS agencia,
		NULL AS conta_corrente,
		NULL AS cnpj_titular,
		NULL AS nome_titular,
		NULL AS transferencia_padrao,
		c.INS_ESTADUAL,
		c.INS_MUNICIPAL,
		NULL AS inscricao_suframa,
		NULL AS tipo_atividade,
		c.CNAE,
		'Não' AS simples_nacional,
		'Não' AS produtor_rural,
		'Sim' AS contribuinte,
		NULL AS tags,
		c.OBS,
		c.txt_restricao,
		NULL AS parcelas_padrao,
		(select nome from amm_usuarios where id = c.ID_VENDEDOR) AS vendedor,
		CASE
			WHEN LTRIM(RTRIM(ISNULL(CAST(c.EMAILS AS VARCHAR(MAX)), ''))) IN ('', '0', 'NULL') THEN NULL
			WHEN CHARINDEX('@', LTRIM(RTRIM(CAST(c.EMAILS AS VARCHAR(MAX))))) = 0 THEN NULL
			WHEN CHARINDEX('.', SUBSTRING(LTRIM(RTRIM(CAST(c.EMAILS AS VARCHAR(MAX)))), CHARINDEX('@', LTRIM(RTRIM(CAST(c.EMAILS AS VARCHAR(MAX))))) + 1, LEN(LTRIM(RTRIM(CAST(c.EMAILS AS VARCHAR(MAX))))))) = 0 THEN NULL
			ELSE LTRIM(RTRIM(CAST(c.EMAILS AS VARCHAR(MAX))))
		END AS email_nf,
		'Sim' AS gerar_boleto,
		CASE
			WHEN NULLIF(REPLACE(REPLACE(REPLACE(LTRIM(RTRIM(CAST(c.CNPJ AS VARCHAR(MAX)))), '.', ''), '-', ''), '/', ''), '') IS NULL THEN NULL
			WHEN LEN(REPLACE(REPLACE(REPLACE(LTRIM(RTRIM(CAST(c.CNPJ AS VARCHAR(MAX)))), '.', ''), '-', ''), '/', '')) < 11 THEN NULL
			ELSE LEFT(REPLACE(REPLACE(REPLACE(LTRIM(RTRIM(CAST(c.CNPJ AS VARCHAR(MAX)))), '.', ''), '-', ''), '/', ''), 14)
		END AS cnpj_entrega,
		LEFT(CAST(c.RAZAO_SOCIAL AS VARCHAR(MAX)), 60) AS nome_entrega,
		c.INS_ESTADUAL AS ie_entrega,
		LEFT(LTRIM(RTRIM(ISNULL(CAST(c.ENDERECO AS VARCHAR(MAX)), ''))), 60) AS endereco_entrega,
		c.endereco_numero AS numero_entrega,
		c.Bairro AS bairro_entrega,
		LEFT(LTRIM(RTRIM(ISNULL(CAST(c.endereco_Complemento AS VARCHAR(MAX)), ''))), 60) AS complemento_entrega,
		NULLIF(
			CASE UPPER(LTRIM(RTRIM(ISNULL(CAST(c.Estado AS VARCHAR(MAX)), ''))))
				WHEN 'NULL'   THEN ''
				WHEN 'NONE'   THEN ''
				WHEN 'ESTADO' THEN ''
				ELSE LTRIM(RTRIM(CAST(c.Estado AS VARCHAR(MAX))))
			END
		, '') AS estado_entrega,
		CASE
			WHEN UPPER(LTRIM(RTRIM(ISNULL(CAST(c.Cidade AS VARCHAR(MAX)), '')))) IN ('NULL', 'NONE', '0', 'RJ', '') THEN NULL
			WHEN LTRIM(RTRIM(CAST(c.Cidade AS VARCHAR(MAX)))) = 'Rio Janeiro' THEN 'Rio de Janeiro'
			WHEN UPPER(LTRIM(RTRIM(CAST(c.Cidade AS VARCHAR(MAX))))) IN ('SÃO MATHEUS', 'SAO MATHEUS') THEN 'São Mateus'
			ELSE LTRIM(RTRIM(CAST(c.Cidade AS VARCHAR(MAX))))
		END AS cidade_entrega,
		NULLIF(LEFT(REPLACE(REPLACE(LTRIM(RTRIM(ISNULL(CAST(c.CEP AS VARCHAR(MAX)), ''))), '-', ''), '.', ''), 8), '') AS cep_entrega,
		NULLIF(LEFT(LTRIM(RTRIM(ISNULL(CAST(c.CONTATOS AS VARCHAR(MAX)), ''))), 15), '') AS telefone_entrega,
		NULL AS limite_credito,
		CASE WHEN c.bloqueiaFaturamento = 1 THEN 'Sim' ELSE 'Não' END AS bloquear_faturamento,
		NULL AS nome_transportadora,
		NULL AS chave_pix
	FROM amm_clientes c
	WHERE c.id = ?`

	var (
		codigoIntegracao    int
		razaoSocial         *string
		cnpj                *string
		nomeFantasia        *string
		cliente             string
		fornecedor          string
		transportadora      string
		funcionario         string
		contato             *string
		endereco            string
		numero              *string
		bairro              *string
		complemento         string
		estado              *string
		cidade              *string
		pais                *string
		cep                 *string
		dddTelefone         string
		telefone            string
		dddTelefone2        *string
		telefone2           *string
		dddFax              *string
		fax                 *string
		emails              *string
		siteCliente         *string
		banco               *string
		agencia             *string
		contaCorrente       *string
		cnpjTitular         *string
		nomeTitular         *string
		transferenciaPadrao *string
		insEstadual         *string
		insMunicipal        *string
		inscricaoSuframa    *string
		tipoAtividade       *string
		cnae                *string
		simplesNacional     string
		produtorRural       string
		contribuinte        string
		tags                *string
		obs                 *string
		txtRestricao        *string
		parcelasPadrao      *string
		vendedor            *string
		emailNf             *string
		gerarBoleto         string
		cnpjEntrega         *string
		nomeEntrega         *string
		ieEntrega           *string
		enderecoEntrega     string
		numeroEntrega       *string
		bairroEntrega       *string
		complementoEntrega  string
		estadoEntrega       *string
		cidadeEntrega       *string
		cepEntrega          *string
		telefoneEntrega     *string
		limiteCredito       *string
		bloquearFaturamento string
		nomeTransportadora  *string
		chavePix            *string
	)

	err := db.QueryRow(query, id).Scan(
		&codigoIntegracao, &razaoSocial, &cnpj, &nomeFantasia,
		&cliente, &fornecedor, &transportadora, &funcionario,
		&contato, &endereco, &numero, &bairro, &complemento,
		&estado, &cidade, &pais, &cep,
		&dddTelefone, &telefone, &dddTelefone2, &telefone2, &dddFax, &fax,
		&emails, &siteCliente,
		&banco, &agencia, &contaCorrente, &cnpjTitular, &nomeTitular, &transferenciaPadrao,
		&insEstadual, &insMunicipal, &inscricaoSuframa, &tipoAtividade, &cnae,
		&simplesNacional, &produtorRural, &contribuinte,
		&tags, &obs, &txtRestricao, &parcelasPadrao, &vendedor,
		&emailNf, &gerarBoleto,
		&cnpjEntrega, &nomeEntrega, &ieEntrega, &enderecoEntrega, &numeroEntrega,
		&bairroEntrega, &complementoEntrega, &estadoEntrega, &cidadeEntrega,
		&cepEntrega, &telefoneEntrega,
		&limiteCredito, &bloquearFaturamento, &nomeTransportadora, &chavePix,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Erro ao buscar cliente por id %d: %v", id, err)
		return nil, err
	}

	return map[string]any{
		"codigo_integracao":    codigoIntegracao,
		"RAZAO_SOCIAL":         razaoSocial,
		"CNPJ":                 cnpj,
		"nome_fantasia":        nomeFantasia,
		"cliente":              cliente,
		"fornecedor":           fornecedor,
		"transportadora":       transportadora,
		"funcionario":          funcionario,
		"contato":              contato,
		"ENDERECO":             endereco,
		"numero":               numero,
		"Bairro":               bairro,
		"complemento":          complemento,
		"Estado":               estado,
		"Cidade":               cidade,
		"Pais":                 pais,
		"CEP":                  cep,
		"ddd_telefone":         dddTelefone,
		"telefone":             telefone,
		"ddd_telefone2":        dddTelefone2,
		"telefone2":            telefone2,
		"ddd_fax":              dddFax,
		"fax":                  fax,
		"EMAILS":               emails,
		"siteCliente":          siteCliente,
		"banco":                banco,
		"agencia":              agencia,
		"conta_corrente":       contaCorrente,
		"cnpj_titular":         cnpjTitular,
		"nome_titular":         nomeTitular,
		"transferencia_padrao": transferenciaPadrao,
		"INS_ESTADUAL":         insEstadual,
		"INS_MUNICIPAL":        insMunicipal,
		"inscricao_suframa":    inscricaoSuframa,
		"tipo_atividade":       tipoAtividade,
		"CNAE":                 cnae,
		"simples_nacional":     simplesNacional,
		"produtor_rural":       produtorRural,
		"contribuinte":         contribuinte,
		"tags":                 tags,
		"OBS":                  obs,
		"txt_restricao":        txtRestricao,
		"parcelas_padrao":      parcelasPadrao,
		"vendedor":             vendedor,
		"email_nf":             emailNf,
		"gerar_boleto":         gerarBoleto,
		"cnpj_entrega":         cnpjEntrega,
		"nome_entrega":         nomeEntrega,
		"ie_entrega":           ieEntrega,
		"endereco_entrega":     enderecoEntrega,
		"numero_entrega":       numeroEntrega,
		"bairro_entrega":       bairroEntrega,
		"complemento_entrega":  complementoEntrega,
		"estado_entrega":       estadoEntrega,
		"cidade_entrega":       cidadeEntrega,
		"cep_entrega":          cepEntrega,
		"telefone_entrega":     telefoneEntrega,
		"limite_credito":       limiteCredito,
		"bloquear_faturamento": bloquearFaturamento,
		"nome_transportadora":  nomeTransportadora,
		"chave_pix":            chavePix,
	}, nil
}
