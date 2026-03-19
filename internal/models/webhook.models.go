package models

type WebhookOsFaturadaResponse struct {
	NumeroOS      string `json:"numero_os"`
	CodigoIntegra string `json:"codigo_integracao"`
	IdOs          int64  `json:"id_os"`
}

type WebhookContaReceberResponseInclude struct {
	CodigoCliente         int64  `json:"codigo_cliente_fornecedor"`
	CodigoConta           int64  `json:"codigo_lancamento_omie"`
	NumeroDocumento       string `json:"numero_documento"`
	NumeroDocumentoFiscal string `json:"numero_documento_fiscal"`
	NumeroPedido          string `json:"numero_pedido"`
}

type WebhookOsIncluidaResponse struct {
	CodigoIntegra string `json:"codigoIntegracao"`
	IdOs          int64  `json:"idOrdemServico"`
	IdCliente     int64  `json:"idCliente"`
	NumeroOs      string `json:"numeroOrdemServico"`
}
