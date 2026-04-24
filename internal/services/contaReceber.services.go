package services

import (
	"encoding/json"
	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repositories"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (s *OmieService) ConsultarContaReceber(req models.ContaReceberRequest) (string, error) {
	url := s.BaseURL + "/api/v1/financas/contareceber/"
	payload := strings.NewReader(`{	
	"call":"ConsultarContaReceber",
	"param":[{
		"codigo_lancamento_omie":` + fmt.Sprint(req.CodigoLancamentoOmie) + `,
		"codigo_lancamento_integracao":"` + req.CodigoLancamentoIntegracao + `"}], 
		"app_key":"` + s.AppKey + `",
		"app_secret":"` + s.AppSecret + `"
		}`)
	//codigo_lancamento_integracao não é obrigatório, pode ser vazio
	httpReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Resposta:", string(body))

	return string(body), nil
}
func (s *OmieService) ListarContasReceber() (map[string]any, error) {
	url := s.BaseURL + "/api/v1/financas/contareceber/"
	payload := strings.NewReader(`{
				"call": "ListarContasReceber",
				"param":[{"pagina":1,"registros_por_pagina":50,"apenas_importado_api":"N"}],
				"app_key": "` + s.AppKey + `",
				"app_secret": "` + s.AppSecret + `"
		}`)
	httpReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(httpReq)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	fmt.Println("Resposta:", result)
	return result, nil
}
func (s *OmieService) GerarBoletoConta(req models.GerarBoletoConta) (map[string]any, error) {
	url := s.BaseURL + "/api/v1/financas/contareceberboleto/"

	if req.NCodTitulo == 0 {
		return map[string]any{"erro": "nCodTitulo é obrigatório"}, fmt.Errorf("nCodTitulo é obrigatório")
	}

	payloadObj := map[string]any{
		"call": "GerarBoleto",
		"param": []map[string]any{
			{
				"nCodTitulo":    req.NCodTitulo,
				"cCodIntTitulo": req.CCodIntTitulo,
			},
		},
		"app_key":    s.AppKey,
		"app_secret": s.AppSecret,
	}

	payloadBytes, err := json.Marshal(payloadObj)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	httpReq, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(httpReq)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return map[string]any{"erro": fmt.Sprintf("omie retornou status %d: %s", resp.StatusCode, string(body))}, fmt.Errorf("omie retornou status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	return result, nil
}
func (s *OmieService) ConsultarBoletoGerado(req models.ConsultaBoletoGerado) (map[string]any, error) {
	url := s.BaseURL + "/api/v1/financas/contareceberboleto/"
	payloadObj := map[string]any{
		"call": "ObterBoleto",
		"param": []map[string]any{
			{
				"nCodTitulo": req.NCodTitulo,
			},
		},
		"app_key":    s.AppKey,
		"app_secret": s.AppSecret,
	}

	payloadBytes, err := json.Marshal(payloadObj)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	httpReq, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(httpReq)
	if err != nil {

		return map[string]any{"erro": err.Error()}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return map[string]any{"erro": fmt.Sprintf("omie retornou status %d: %s", resp.StatusCode, string(body))}, fmt.Errorf("omie retornou status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	fmt.Println(result)
	_, err = repositories.InsertLinkBoletoGerado(database.ConnectToDB(), req.NCodTitulo, result["cLinkBoleto"].(string))
	return result, nil
}

func (s *OmieService) UpsertNFSEGerada(req models.ConsultaNFSEGerada) (map[string]any, error) {
	urlListar := s.BaseURL + "/api/v1/servicos/nfse/"
	payloadListarObj := map[string]any{
		"call": "ListarNFSEs",
		"param": []map[string]any{
			{
				"nPagina":       1,
				"nRegPorPagina": 20,
				"nNumeroNFSe":   req.NNumeroNFSe,
			},
		},
		"app_key":    s.AppKey,
		"app_secret": s.AppSecret,
	}

	payloadListarBytes, err := json.Marshal(payloadListarObj)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	httpReqListar, err := http.NewRequest("POST", urlListar, strings.NewReader(string(payloadListarBytes)))
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	httpReqListar.Header.Set("Content-Type", "application/json")

	respListar, err := (&http.Client{}).Do(httpReqListar)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	defer respListar.Body.Close()

	bodyListar, err := io.ReadAll(respListar.Body)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	if respListar.StatusCode < 200 || respListar.StatusCode >= 300 {
		return map[string]any{"erro": fmt.Sprintf("omie retornou status %d: %s", respListar.StatusCode, string(bodyListar))}, fmt.Errorf("omie retornou status %d: %s", respListar.StatusCode, string(bodyListar))
	}

	var listarResp map[string]any
	if err := json.Unmarshal(bodyListar, &listarResp); err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	nfseEncontradas, ok := listarResp["nfseEncontradas"].([]any)
	if !ok || len(nfseEncontradas) == 0 {
		return map[string]any{"erro": "nenhuma NFSe encontrada para o nNumeroNFSe informado"}, fmt.Errorf("nenhuma NFSe encontrada para o nNumeroNFSe informado")
	}

	nfse0, ok := nfseEncontradas[0].(map[string]any)
	if !ok {
		return map[string]any{"erro": "formato inválido de nfseEncontradas[0]"}, fmt.Errorf("formato inválido de nfseEncontradas[0]")
	}

	cabecalho, ok := nfse0["Cabecalho"].(map[string]any)
	if !ok {
		return map[string]any{"erro": "Cabecalho não encontrado"}, fmt.Errorf("Cabecalho não encontrado")
	}

	emissao, ok := nfse0["Emissao"].(map[string]any)
	if !ok {
		return map[string]any{"erro": "Emissao não encontrada"}, fmt.Errorf("Emissao não encontrada")
	}

	ordermServico, ok := nfse0["OrdemServico"].(map[string]any)
	if !ok {
		return map[string]any{"erro": "OrdemServico não encontrada"}, fmt.Errorf("OrdemServico não encontrada")
	}

	nCodNFFloat, ok := cabecalho["nCodNF"].(float64)
	if !ok {
		return map[string]any{"erro": "nCodNF não encontrado ou inválido"}, fmt.Errorf("nCodNF não encontrado ou inválido")
	}
	nCodNF := int64(nCodNFFloat)

	cDataEmissao, _ := emissao["cDataEmissao"].(string)
	fmt.Printf("%v\n", ordermServico)
	codigoOs, _ := ordermServico["nCodigoOS"].(float64)

	urlObter := s.BaseURL + "/api/v1/servicos/osdocs/"
	payloadObterObj := map[string]any{
		"call": "ObterNFSe",
		"param": []map[string]any{
			{
				"nIdNf": nCodNF,
			},
		},
		"app_key":    s.AppKey,
		"app_secret": s.AppSecret,
	}

	payloadObterBytes, err := json.Marshal(payloadObterObj)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	httpReqObter, err := http.NewRequest("POST", urlObter, strings.NewReader(string(payloadObterBytes)))
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	httpReqObter.Header.Set("Content-Type", "application/json")

	respObter, err := (&http.Client{}).Do(httpReqObter)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	defer respObter.Body.Close()

	bodyObter, err := io.ReadAll(respObter.Body)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}

	if respObter.StatusCode < 200 || respObter.StatusCode >= 300 {
		return map[string]any{"erro": fmt.Sprintf("omie retornou status %d: %s", respObter.StatusCode, string(bodyObter))}, fmt.Errorf("omie retornou status %d: %s", respObter.StatusCode, string(bodyObter))
	}

	var obterResp map[string]any
	if err := json.Unmarshal(bodyObter, &obterResp); err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	fmt.Printf("RESPOSTA: %v\n", obterResp)

	cXmlNFSe, _ := obterResp["cXmlNFSe"].(string)
	cUrlNFSe, _ := obterResp["cUrlNFSe"].(string)
	cLinkPortal, _ := obterResp["cLinkPortal"].(string)
	cPdfNfse, _ := obterResp["cPdfNFSe"].(string)
	cNumNFSe, _ := obterResp["cNumNFSe"].(string)
	_, err = repositories.UpsertNFSEGerada(database.ConnectToDB(), nCodNF, codigoOs, cDataEmissao, cXmlNFSe, cUrlNFSe, cLinkPortal, cNumNFSe, cPdfNfse)
	if err != nil {
		return map[string]any{"erro": err.Error()}, err
	}
	return map[string]any{
		"nCodNF":       nCodNF,
		"codigoOs":     codigoOs,
		"cDataEmissao": cDataEmissao,
		"cXmlNFSe":     cXmlNFSe,
		"cUrlNFSe":     cUrlNFSe,
		"cLinkPortal":  cLinkPortal,
		"cNumNFSe":     cNumNFSe,
	}, nil
}
