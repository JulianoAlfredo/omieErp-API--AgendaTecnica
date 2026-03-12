package models

type ClienteRequest struct {
	CodCliIntegra string `json:"codigo_cliente_integracao"`
	Email         string `json:"email"`
	RazaoSocial   string `json:"razao_social"`
	NomeFantasia  string `json:"nome_fantasia"`
	CnpjCpf       string `json:"cnpj_cpf"`
}

type ClienteImporta struct {
	Id string `json:"id"`
}
