package models

type ContaCorrente struct {
	NCodCC              int64  `json:"nCodCC"`
	CCodCCInt           string `json:"cCodCCInt"`
	Descricao           string `json:"descricao"`
	Tipo                string `json:"tipo"`
	CodigoBanco         string `json:"codigo_banco"`
	CodigoAgencia       string `json:"codigo_agencia"`
	NumeroContaCorrente string `json:"numero_conta_corrente"`
	Inativo             string `json:"inativo"`
	Bloqueado           string `json:"bloqueado"`
}
