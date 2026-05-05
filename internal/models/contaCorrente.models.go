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

type DetalhePix struct {
	ChavePixRecebedor  string `json:"chavePixRecebedor"`
	CpfCnpjPagador     string `json:"cpfCnpjPagador"`
	DescricaoPix       string `json:"descricaoPix"`
	EndToEndId         string `json:"endToEndId"`
	NomeEmpresaPagador string `json:"nomeEmpresaPagador"`
	NomePagador        string `json:"nomePagador"`
	TipoDetalhe        string `json:"tipoDetalhe"`
	TxId               string `json:"txId"`
}

type DetalheBoletoCobanca struct {
	Abatimento     string `json:"abatimento"`
	CodBarras      string `json:"codBarras"`
	CpfCnpj        string `json:"cpfCnpj"`
	DataEmissao    string `json:"dataEmissao"`
	DataLimite     string `json:"dataLimite"`
	DataTransacao  string `json:"dataTransacao"`
	DataVencimento string `json:"dataVencimento"`
	Desconto1      string `json:"desconto1"`
	Desconto2      string `json:"desconto2"`
	Desconto3      string `json:"desconto3"`
	Juros          string `json:"juros"`
	Multa          string `json:"multa"`
	Nome           string `json:"nome"`
	NossoNumero    string `json:"nossoNumero"`
	SeuNumero      string `json:"seuNumero"`
	TipoDetalhe    string `json:"tipoDetalhe"`
}

type ExtratoInterResponse struct {
	DataInclusao       string                `json:"dataInclusao"`
	DataTransacao      string                `json:"dataTransacao"`
	Descricao          string                `json:"descricao"`
	IdTransacao        string                `json:"idTransacao"`
	NumeroDocumento    string                `json:"numeroDocumento"`
	SeuNumero          string                `json:"seuNumero,omitempty"`
	TipoOperacao       string                `json:"tipoOperacao"`
	TipoTransacao      string                `json:"tipoTransacao"`
	Titulo             string                `json:"titulo"`
	Valor              string                `json:"valor"`
	Detalhes           *DetalheBoletoCobanca `json:"detalhes,omitempty"`
	DetalhePixRecebido *DetalhePix           `json:"detalhesPix,omitempty"`
	DetalheBoleto      *DetalheBoletoCobanca `json:"detalhesBoleto,omitempty"`
}

type ExtratoInterEnvelope struct {
	Transacoes []ExtratoInterResponse `json:"transacoes"`
}

type ExtratoResumoItem struct {
	Indicador string `json:"indicador"`
	Valor     string `json:"valor"`
}

type ExtratoCompletoResponse struct {
	Resumo     []ExtratoResumoItem    `json:"resumo"`
	Transacoes []ExtratoInterResponse `json:"transacoes"`
}
