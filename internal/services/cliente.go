package services

import (
	"example/web-service-gin/internal/models"
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

func (s *OmieService) ListarClientes() (string, error) {
	url := s.BaseURL + "/api/v1/geral/clientes/"
	payload := strings.NewReader(`{
				"call": "ListarClientes",
				"param":[{"pagina":1,"registros_por_pagina":20,"apenas_importado_api":"N"}],
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

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Resposta:", string(body))

	return string(body), nil
}
