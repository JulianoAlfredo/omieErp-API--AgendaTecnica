// curl -s https://app.omie.com.br/api/v1/geral/contacorrente/ \
//  -H 'Content-type: application/json' \
//  -d '{"call":"IncluirContaCorrente","param":[{"cCodCCInt":"MyCC0001","tipo_conta_corrente":"CX","codigo_banco":"999","descricao":"Caixinha","saldo_inicial":0}],"app_key":"#APP_KEY#","app_secret":"#APP_SECRET#"}'

// {
//   "nCodCC": 2902849401,
//   "cCodCCInt": "MyCC0001",
//   "cCodStatus": "0",
//   "cDesStatus": "Conta corrente incluída com sucesso!"
// }

package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repositories"
)

func (s *OmieService) ListarContasCorrente() ([]models.ContaCorrente, error) {
	url := s.BaseURL + "/api/v1/geral/contacorrente/"

	payload := strings.NewReader(`{
				"call": "ListarContasCorrentes",
				"param":[{"pagina":1,"registros_por_pagina":50,"apenas_importado_api":"N"}],
				"app_key": "` + s.AppKey + `",
				"app_secret": "` + s.AppSecret + `"
		}`)
	httpReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw struct {
		ListarContasCorrentes []models.ContaCorrente `json:"ListarContasCorrentes"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	return raw.ListarContasCorrentes, nil
}

func parseValorMonetario(valor string) float64 {
	v := strings.TrimSpace(valor)
	v = strings.ReplaceAll(v, "R$", "")
	v = strings.ReplaceAll(v, " ", "")

	temPonto := strings.Contains(v, ".")
	temVirgula := strings.Contains(v, ",")

	if temPonto && temVirgula {
		// Formato BR com milhar e decimal: 1.234,56
		v = strings.ReplaceAll(v, ".", "")
		v = strings.ReplaceAll(v, ",", ".")
	} else if temVirgula {
		// Formato BR sem milhar: 123,45
		v = strings.ReplaceAll(v, ",", ".")
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}
	return f
}

func formatarNumeroBR(valor float64) string {
	negativo := valor < 0
	if negativo {
		valor = math.Abs(valor)
	}

	base := fmt.Sprintf("%.2f", valor)
	partes := strings.Split(base, ".")
	inteiro := partes[0]
	decimais := "00"
	if len(partes) > 1 {
		decimais = partes[1]
	}

	n := len(inteiro)
	if n > 3 {
		var b strings.Builder
		resto := n % 3
		if resto > 0 {
			b.WriteString(inteiro[:resto])
			if resto < n {
				b.WriteString(".")
			}
		}
		for i := resto; i < n; i += 3 {
			b.WriteString(inteiro[i : i+3])
			if i+3 < n {
				b.WriteString(".")
			}
		}
		inteiro = b.String()
	}

	resultado := inteiro + "," + decimais
	if negativo {
		return "-" + resultado
	}

	return resultado
}

func formatarBRL(valor float64) string {
	return "R$ " + formatarNumeroBR(valor)
}

func (s *OmieService) ExtratoCompleto() (*models.ExtratoCompletoResponse, error) {
	certFile := os.Getenv("INTER_CERT_FILE")
	keyFile := os.Getenv("INTER_KEY_FILE")
	clientID := os.Getenv("INTER_CLIENT_ID")
	clientSecret := os.Getenv("INTER_CLIENT_SECRET")
	contaCorrente := os.Getenv("INTER_CONTA_CORRENTE")

	if certFile == "" || keyFile == "" || clientID == "" || clientSecret == "" || contaCorrente == "" {
		return nil, fmt.Errorf("variaveis obrigatorias ausentes: INTER_CERT_FILE, INTER_KEY_FILE, INTER_CLIENT_ID, INTER_CLIENT_SECRET, INTER_CONTA_CORRENTE")
	}

	oauthURL := os.Getenv("INTER_OAUTH_URL")
	if oauthURL == "" {
		oauthURL = "https://cdpj.partners.bancointer.com.br/oauth/v2/token"
	}

	extratoURL := os.Getenv("INTER_EXTRATO_URL")
	if extratoURL == "" {
		extratoURL = "https://cdpj.partners.bancointer.com.br/banking/v2/extrato/completo"
	}

	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar certificado/chave da Inter: %w", err)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{clientCert},
		},
	}

	client := &http.Client{Transport: transport}

	oauthValues := url.Values{}
	oauthValues.Set("client_id", clientID)
	oauthValues.Set("client_secret", clientSecret)
	oauthValues.Set("scope", "extrato.read")
	oauthValues.Set("grant_type", "client_credentials")

	oauthReq, err := http.NewRequest("POST", oauthURL, strings.NewReader(oauthValues.Encode()))
	if err != nil {
		return nil, err
	}
	oauthReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	oauthResp, err := client.Do(oauthReq)
	if err != nil {
		return nil, err
	}
	defer oauthResp.Body.Close()

	oauthBody, err := io.ReadAll(oauthResp.Body)
	if err != nil {
		return nil, err
	}

	if oauthResp.StatusCode < http.StatusOK || oauthResp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("erro ao obter token Inter: status=%d body=%s", oauthResp.StatusCode, string(oauthBody))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(oauthBody, &tokenResp); err != nil {
		return nil, err
	}
	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("resposta da Inter sem access_token")
	}

	now := time.Now()
	primeiroDiaMes := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	ultimoDiaMes := primeiroDiaMes.AddDate(0, 1, -1)

	dataInicio := primeiroDiaMes.Format("2006-01-02")
	dataFim := ultimoDiaMes.Format("2006-01-02")

	fmt.Println(dataInicio)
	fmt.Println(dataFim)

	extratoReq, err := http.NewRequest("GET", extratoURL, nil)
	if err != nil {
		return nil, err
	}

	query := extratoReq.URL.Query()
	query.Set("dataInicio", dataInicio)
	query.Set("dataFim", dataFim)
	query.Set("pagina", "0")
	query.Set("tamanhoPagina", "9999")
	extratoReq.URL.RawQuery = query.Encode()

	extratoReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	extratoReq.Header.Set("x-conta-corrente", contaCorrente)
	extratoReq.Header.Set("Content-Type", "application/json")

	extratoResp, err := client.Do(extratoReq)
	if err != nil {
		return nil, err
	}
	defer extratoResp.Body.Close()

	extratoBody, err := io.ReadAll(extratoResp.Body)
	if err != nil {
		return nil, err
	}

	if extratoResp.StatusCode < http.StatusOK || extratoResp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("erro ao consultar extrato Inter: status=%d body=%s", extratoResp.StatusCode, string(extratoBody))
	}

	var envelope models.ExtratoInterEnvelope
	if err := json.Unmarshal(extratoBody, &envelope); err != nil {
		return nil, err
	}

	transacoes := envelope.Transacoes
	for i := range transacoes {
		if transacoes[i].DetalheBoleto == nil && transacoes[i].Detalhes != nil {
			transacoes[i].DetalheBoleto = transacoes[i].Detalhes
		}

		if transacoes[i].SeuNumero == "" {
			if transacoes[i].DetalheBoleto != nil && transacoes[i].DetalheBoleto.SeuNumero != "" {
				transacoes[i].SeuNumero = transacoes[i].DetalheBoleto.SeuNumero
			} else if transacoes[i].Detalhes != nil && transacoes[i].Detalhes.SeuNumero != "" {
				transacoes[i].SeuNumero = transacoes[i].Detalhes.SeuNumero
			} else if strings.Contains(strings.ToLower(transacoes[i].TipoTransacao), "boleto") && transacoes[i].NumeroDocumento != "" {
				transacoes[i].SeuNumero = transacoes[i].NumeroDocumento
			}
		}
	}

	var totalRecebidoPix float64
	var totalRecebidoBoleto float64
	var totalPagoTarifas float64

	for _, transacao := range transacoes {
		textoAnalise := strings.ToLower(transacao.TipoTransacao + " " + transacao.Descricao + " " + transacao.Titulo)
		tipoOperacao := strings.ToLower(strings.TrimSpace(transacao.TipoOperacao))
		valorBruto := parseValorMonetario(transacao.Valor)
		valor := math.Abs(valorBruto)

		ehCredito := tipoOperacao == "c" || strings.Contains(tipoOperacao, "credito") || strings.Contains(tipoOperacao, "entrada")
		ehDebito := tipoOperacao == "d" || strings.Contains(tipoOperacao, "debito") || strings.Contains(tipoOperacao, "saida")

		if !ehCredito && !ehDebito {
			ehCredito = valorBruto > 0
			ehDebito = valorBruto < 0
		}

		ehPix := transacao.DetalhePixRecebido != nil || strings.Contains(textoAnalise, "pix")
		ehBoleto := transacao.DetalheBoleto != nil || transacao.Detalhes != nil || strings.Contains(textoAnalise, "boleto")
		ehTarifa := strings.Contains(textoAnalise, "tarifa")

		if ehPix && ehCredito {
			totalRecebidoPix += valor
		}

		if ehBoleto && ehCredito {
			totalRecebidoBoleto += valor
		}

		if ehTarifa && ehDebito {
			totalPagoTarifas += valor
		}
	}

	resumo := []models.ExtratoResumoItem{
		{Indicador: "Transacoes total", Valor: strconv.Itoa(len(transacoes))},
		{Indicador: "Valores recebidos: Pix", Valor: formatarBRL(totalRecebidoPix)},
		{Indicador: "Valores recebidos: boleto", Valor: formatarBRL(totalRecebidoBoleto)},
		{Indicador: "Valores pagos: tarifas", Valor: formatarBRL(totalPagoTarifas)},
	}

	return &models.ExtratoCompletoResponse{
		Resumo:     resumo,
		Transacoes: transacoes,
	}, nil
}

func (s *OmieService) SincronizarBaixasOmie(nCodCC int64, dataInicial string, dataFinal string) (*models.SincronizarBaixasResult, error) {
	endpoint := s.BaseURL + "/api/v1/financas/extrato/"

	resultado := &models.SincronizarBaixasResult{
		NaoEncontrados: []string{},
	}

	db := database.ConnectToDB()
	defer db.Close()

	payload := fmt.Sprintf(`{"call":"ListarExtrato","param":[{"nCodCC":%d,"cCodIntCC":"","dPeriodoInicial":"%s","dPeriodoFinal":"%s"}],"app_key":"%s","app_secret":"%s"}`,
		nCodCC, dataInicial, dataFinal, s.AppKey, s.AppSecret)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição OMIE: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar API OMIE: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta OMIE: %w", err)
	}

	var extratoResp models.ExtratoOmieResponse
	if err := json.Unmarshal(body, &extratoResp); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta OMIE: %w, body=%s", err, string(body))
	}

	for _, mov := range extratoResp.Movimentos {
		if mov.NCodLancamento == 0 || mov.CNumero == "" {
			continue
		}
		if mov.CNatureza != "R" {
			continue
		}

		resultado.Total++

		rows, err := repositories.UpdateBaixaPorNumeroRps(db, mov.CNumero, mov.NCodLancamento, mov.NValorDocumento, mov.CObservacoes, mov.CDataLancamento)
		if err != nil {
			fmt.Printf("[SincronizarBaixas] Erro ao atualizar RPS %s: %v\n", mov.CNumero, err)
			resultado.NaoEncontrados = append(resultado.NaoEncontrados, mov.CNumero)
			continue
		}
		if rows == 0 {
			fmt.Printf("[SincronizarBaixas] RPS não encontrado na base: %s\n", mov.CNumero)
			resultado.NaoEncontrados = append(resultado.NaoEncontrados, mov.CNumero)
		} else {
			resultado.Sincronizados++
			fmt.Printf("[SincronizarBaixas] RPS %s sincronizado (lancamento=%d, valor=%.2f)\n", mov.CNumero, mov.NCodLancamento, mov.NValorDocumento)
		}
	}

	return resultado, nil
}
