package services

import (
	"encoding/json"
	"example/web-service-gin/internal/models"
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
