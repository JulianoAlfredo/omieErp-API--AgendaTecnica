package services

import (
	"encoding/json"
	"example/web-service-gin/internal/models"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Quero inserir essas 3 funcs em algum utils ou helpers, mas por enquanto vou deixar aqui para não perder o código
func formatEmail(email interface{}) string {
	if email == nil {
		return "null"
	}
	data, _ := json.Marshal(email)
	return string(data)
}
func formatInfoAdicionais(info interface{}) string {
	if info == nil {
		return "null"
	}
	data, _ := json.Marshal(info)
	return string(data)
}
func formatServicos(servicos interface{}) string {
	if servicos == nil {
		return "[]"
	}
	data, _ := json.Marshal(servicos)
	return string(data)
}

// JA CRIA OS EM EXECUCAO... AS OS AQUI VAO SERVIR APENAS PARA FATURAMENTO, SEM SER ACOMPANHAMENTO DA OS MESMO
func (s *OmieService) CriarOrdemServico(req models.OrdemServicoRequest) (string, error) {
	url := s.BaseURL + "/api/v1/servicos/os/"
	fmt.Println("REQ:", req.Cabecalho.CCodIntOS)
	payload := strings.NewReader(`{
		"call": "IncluirOS",
		"param": [{
			"Cabecalho": {
				"cCodIntOS": "` + req.Cabecalho.CCodIntOS + `",
				"cCodParc": "` + req.Cabecalho.CCodParc + `",
				"cEtapa": "` + req.Cabecalho.CEtapa + `",
				"dDtPrevisao": "` + req.Cabecalho.DDtPrevisao + `",
				"nCodCli": ` + fmt.Sprint(req.Cabecalho.NCodCli) + `,
				"nQtdeParc": ` + fmt.Sprint(req.Cabecalho.NQtdeParc) + `
			},
			"Departamentos": [],
			"Email": ` + formatEmail(req.Email) + `,
			"InformacoesAdicionais": ` + formatInfoAdicionais(req.InformacoesAdicionais) + `,
			"ServicosPrestados": ` + formatServicos(req.ServicosPrestados) + `
		}],
		"app_key": "` + s.AppKey + `",
		"app_secret": "` + s.AppSecret + `"
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
	fmt.Println(err)
	fmt.Println(resp)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Resposta:", string(body))

	return "Cliente cadastrado com sucesso", nil
}
func (s *OmieService) ListarOrdemServico() (string, error) {
	url := s.BaseURL + "/api/v1/servicos/os/"
	payload := strings.NewReader(`{
	"call":"ListarOS",
	"param":[{
		"pagina":1,
		"registros_por_pagina":50,
		"apenas_importado_api":"N"
		}],
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
func (s *OmieService) FaturarOrdemServico(req models.FaturaOrdemServicoRequest) (string, error) {
	url := s.BaseURL + "/api/v1/servicos/osp/"
	payload := strings.NewReader(`{
		"call":"FaturarOS",
		"param":[{
			"nCodOS":` + fmt.Sprint(req.NCodOS) + `
		}],
		"app_key": "` + s.AppKey + `",
		"app_secret": "` + s.AppSecret + `"
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
	fmt.Println(err)
	fmt.Println(resp)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Resposta:", string(body))

	return "OS FATURADA com sucesso", nil

}

func (s *OmieService) ConsultarOsFase(req models.ListarOSResponse) (map[string]any, error) {
	url := s.BaseURL + "/api/v1/servicos/os/"
	payload := strings.NewReader(`{
	"call":"ConsultarOS",
	"param":[{
		"cCodIntOS":"` + req.CCodIntOS + `",
		"nCodOS":` + fmt.Sprint(req.NCodOS) + `,
		"cNumOS":"` + req.CNumOS + `"
		}],
	"app_key":"` + s.AppKey + `",
	"app_secret":"` + s.AppSecret + `"
	}`)

	httpReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var resultMap map[string]any
	err = json.Unmarshal(body, &resultMap)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", resultMap)

	return resultMap, nil

}

func (s *OmieService) VerificaOsFaturada(req models.ListarOSResponse) (map[string]any, error) {
	url := s.BaseURL + "/api/v1/servicos/os/"
	payload := strings.NewReader(`{
	"call":"ConsultarOS",
	"param":[{
		"cCodIntOS":"` + req.CCodIntOS + `",
		"nCodOS":` + fmt.Sprint(req.NCodOS) + `,
		"cNumOS":"` + req.CNumOS + `"
		}],
	"app_key":"` + s.AppKey + `",
	"app_secret":"` + s.AppSecret + `"
	}`)

	httpReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Erro API: Status %d - Resposta: %s\n", resp.StatusCode, string(bodyBytes))
	}

	var resultMap map[string]any
	err = json.Unmarshal(body, &resultMap)
	if err != nil {
		return nil, err
	}
	cabecalho, ok := resultMap["Cabecalho"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("campo Cabecalho não encontrado ou inválido")
	}

	fmt.Printf("cEtapa: %v\n", cabecalho["cEtapa"])
	codOs, ok := resultMap["Cabecalho"].(map[string]any)["nCodOS"]
	if !ok {
		return nil, fmt.Errorf("campo nCodOS não encontrado ou inválido")
	}
	fmt.Println(codOs)

	infoCadastro, ok := resultMap["InfoCadastro"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("campo InfoCadastro não encontrado ou inválido")
	}
	cFaturada, ok := infoCadastro["cFaturada"].(string)
	if !ok {
		return nil, fmt.Errorf("campo cFaturada não encontrado ou inválido")
	}
	fmt.Printf("cFaturada: %v\n", cFaturada)

	return map[string]any{"cFaturada": cFaturada, "codigoOs": codOs}, nil

}
