package models

type ServicoRequest struct {
	CodInterno    string `json:"codInterno"`
	Descricao     string `json:"descricao"`
	CCodigo       string `json:"cCodigo"`
	PrecoUnitario string `json:"precoUnitario"`
}
