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

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	fmt.Println("Resposta:", result)
	return result, nil
}

func (s *OmieService) GerarBoletoConta(req models.GerarBoletoConta) (string, error) {
	url := s.BaseURL + "/api/v1/financas/contareceberboleto/"
	payload := strings.NewReader(`{
				"call": "GerarBoletoConta",
	"param":[{
		"call":"GerarBoleto",
		"param":[{"nCodTitulo":` + fmt.Sprint(req.NCodTitulo) + `,
		"cCodIntTitulo":"` + req.CCodIntTitulo + `"}],
				"app_key": "` + s.AppKey + `",
				"app_secret": "` + s.AppSecret + `"
		}`)
	httpReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
