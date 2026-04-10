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
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"example/web-service-gin/internal/models"
)

func (s *OmieService) ListarContasCorrente() ([]models.ContaCorrente, error) {
	url := s.BaseURL + "/api/v1/geral/contacorrente/"

	payload := strings.NewReader(`{
				"call": "ListarContasCorrentes",
				"param":[{"pagina":1,"registros_por_pagina":50,"apenas_importado_api":"N"}],
				"app_key": "` + s.AppKey + `",
				"app_secret": "` + s.AppSecret + `"
		}`)
	httpReq, err := http.NewRequest("GET", url, payload)
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
