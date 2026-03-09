package models

type ContaReceberRequest struct {
	CodigoLancamentoOmie       int    `json:"codigo_lancamento_omie"`
	CodigoLancamentoIntegracao string `json:"codigo_lancamento_integracao"`
}

type GerarBoletoConta struct {
	NCodTitulo    int64  `json:"nCodTitulo" binding:"required"`
	CCodIntTitulo string `json:"cCodIntTitulo"`
}
