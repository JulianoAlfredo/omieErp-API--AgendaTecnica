package models

type ClienteRequest struct {
	CodCliIntegra string `json:"codigo_cliente_integracao"`
	Email         string `json:"email"`
	RazaoSocial   string `json:"razao_social"`
	NomeFantasia  string `json:"nome_fantasia"`
	CnpjCpf       string `json:"cnpj_cpf"`
}

type ClienteConsulta struct {
	CodCliIntegra string `json:"codigo_cliente_integracao"`
}

type ClienteImporta struct {
	Id int `json:"id"`
}

type OmieClienteCadastro struct {
	CodigoClienteIntegracao string `json:"codigo_cliente_integracao"`
	CodigoClienteOmie       int64  `json:"codigo_cliente_omie"`
	CnpjCpf                 string `json:"cnpj_cpf"`
}

type OmieListarClientesResponse struct {
	Pagina           int                   `json:"pagina"`
	TotalDePaginas   int                   `json:"total_de_paginas"`
	Registros        int                   `json:"registros"`
	TotalDeRegistros int                   `json:"total_de_registros"`
	ClientesCadastro []OmieClienteCadastro `json:"clientes_cadastro"`
}
