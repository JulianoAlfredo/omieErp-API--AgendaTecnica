package models

type ContaReceberRequest struct {
	CodigoLancamentoOmie       int    `json:"codigo_lancamento_omie"`
	CodigoLancamentoIntegracao string `json:"codigo_lancamento_integracao"`
}

type GerarBoletoConta struct {
	NCodTitulo    int    `json:"nCodTitulo"`
	CCodIntTitulo string `json:"cCodIntTitulo, string"`
}
