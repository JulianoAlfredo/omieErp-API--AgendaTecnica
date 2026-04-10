package services

import (
	"encoding/json"
	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repositories"
	"fmt"
	"io"
	"log"
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

func (s *OmieService) ConsultarCliente(req models.ClienteConsulta) (map[string]any, error) {
	url := s.BaseURL + "/api/v1/geral/clientes/"
	payload := strings.NewReader(`{
		"call":"ConsultarCliente",
		"param":[{"codigo_cliente_omie":0,"codigo_cliente_integracao":"` + req.CodCliIntegra + `"}],
		"app_key":"` + s.AppKey + `",
		"app_secret":"` + s.AppSecret + `"
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

func (s *OmieService) SincronizarClientes() (string, error) {
	url := s.BaseURL + "/api/v1/geral/clientes/"

	db := database.ConnectToDB()
	defer db.Close()

	if err := repositories.CriarTabelaRelacaoClientes(db); err != nil {
		return "", err
	}

	pagina := 1
	totalPaginas := 1
	totalSincronizados := 0

	for pagina <= totalPaginas {
		payload := strings.NewReader(fmt.Sprintf(`{
			"call": "ListarClientes",
			"param":[{"pagina":%d,"registros_por_pagina":50,"apenas_importado_api":"N"}],
			"app_key": "%s",
			"app_secret": "%s"
		}`, pagina, s.AppKey, s.AppSecret))

		httpReq, err := http.NewRequest("POST", url, payload)
		if err != nil {
			return "", err
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := (&http.Client{}).Do(httpReq)
		if err != nil {
			return "", err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", err
		}

		var result models.OmieListarClientesResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return "", err
		}

		totalPaginas = result.TotalDePaginas

		for _, cliente := range result.ClientesCadastro {
			if err := repositories.UpsertRelacaoCliente(db, cliente.CodigoClienteIntegracao, cliente.CodigoClienteOmie, cliente.CnpjCpf); err != nil {
				log.Printf("Erro ao sincronizar cliente %s: %v", cliente.CodigoClienteIntegracao, err)
			} else {
				totalSincronizados++
			}
		}

		log.Printf("Página %d/%d processada", pagina, totalPaginas)
		pagina++
	}

	return fmt.Sprintf("Sincronização concluída: %d clientes sincronizados", totalSincronizados), nil
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
