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

func (s *OmieService) CriarCliente(req models.ClienteRequest) (string, error) {

	url := s.BaseURL + "/api/v1/geral/clientes/"

	payload := strings.NewReader(`{
		"call":"IncluirCliente",
		"param":[{"codigo_cliente_integracao": "` + req.CodCliIntegra + `",
		"email":"` + req.Email + `",
		"razao_social":"` + req.RazaoSocial + `",
		"nome_fantasia":"` + req.NomeFantasia + `",
		"cnpj_cpf":"` + req.CnpjCpf + `"}],
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
func (s *OmieService) ListarClientes() (map[string]any, error) {
	url := s.BaseURL + "/api/v1/geral/clientes/"
	payload := strings.NewReader(`{
        "call": "ListarClientes",
        "param":[{"pagina":1,"registros_por_pagina":20,"apenas_importado_api":"N"}],
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

func (s *OmieService) ImportarEmpresa(req models.ClienteImporta) (string, error) {

	url := s.BaseURL + "/api/v1/geral/clientes/"

	db := database.ConnectToDB()
	empresa := repositories.SearchClients(db, req.Id)

	payload := strings.NewReader(`{
	"call":"IncluirCliente",
	"param":[{"codigo_cliente_integracao": "` + req.Id + `",
	"email":"` + empresa[0]["emails"].(string) + `",
	"razao_social":"` + empresa[0]["razao_social"].(string) + `",
	"nome_fantasia":"` + empresa[0]["nome_fantasia"].(string) + `",
	"cnpj_cpf":"` + empresa[0]["cnpj"].(string) + `"}],
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

	return "Cliente importado com sucesso", nil
}
