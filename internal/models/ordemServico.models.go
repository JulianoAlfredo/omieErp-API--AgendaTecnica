package models

type OrdemServicoRequest struct {
	Cabecalho struct {
		CCodIntOS   string `json:"cCodIntOS"`
		CCodParc    string `json:"cCodParc"`
		CEtapa      string `json:"cEtapa"`
		DDtPrevisao string `json:"dDtPrevisao"`
		NCodCli     int    `json:"nCodCli"`
		NQtdeParc   int    `json:"nQtdeParc"`
	} `json:"Cabecalho"`
	Departamentos []interface{} `json:"Departamentos"`
	Email         struct {
		CEnvBoleto  string `json:"cEnvBoleto"`
		CEnvLink    string `json:"cEnvLink"`
		CEnviarPara string `json:"cEnviarPara"`
	} `json:"Email"`
	InformacoesAdicionais struct {
		CCidPrestServ string `json:"cCidPrestServ"`
		CCodeCateg    string `json:"cCodCateg"`
		CDadosAdicNF  string `json:"cDadosAdicNF"`
		NCodCC        int    `json:"nCodCC"`
	} `json:"InformacoesAdicionais"`
	ServicosPrestados []struct {
		CCodeServLC116 string  `json:"cCodServLC116"`
		CCodeServMun   string  `json:"cCodServMun"`
		CDadosAdicItem string  `json:"cDadosAdicItem"`
		CDescServ      string  `json:"cDescServ"`
		CRetemISS      string  `json:"cRetemISS"`
		CTribServ      string  `json:"cTribServ"`
		NQtde          float64 `json:"nQtde"`
		NValUnit       float64 `json:"nValUnit"`
		NCodServico    int     `json:"nCodServico"`
		Impostos       struct {
			CRetemIRRF  string  `json:"cRetemIRRF"`
			CRetemPIS   string  `json:"cRetemPIS"`
			NAliqCOFINS float64 `json:"nAliqCOFINS"`
			NAliqCSLL   float64 `json:"nAliqCSLL"`
			NAliqIRRF   float64 `json:"nAliqIRRF"`
			NAliqISS    float64 `json:"nAliqISS"`
			NAliqPIS    float64 `json:"nAliqPIS"`
		} `json:"Impostos"`
	} `json:"ServicosPrestados"`
}

type FaturaOrdemServicoRequest struct {
	CCodIntOS string `json:"cCodIntOS"`
	NCodOS    int    `json:"nCodOS"`
}
