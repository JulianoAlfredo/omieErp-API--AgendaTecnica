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
	"strconv"
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

func (s *OmieService) ImportarEmpresa(req models.ClienteImporta) (map[string]any, error) {
	url := s.BaseURL + "/api/v1/geral/clientes/"

	db := database.ConnectToDB()
	defer db.Close()

	empresa, err := repositories.SearchClientByField(db, req.Id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cliente no banco: %w", err)
	}
	if empresa == nil {
		return nil, fmt.Errorf("cliente com id %d não encontrado", req.Id)
	}

	pStr := func(key string) string {
		v := empresa[key]
		if v == nil {
			return ""
		}
		if p, ok := v.(*string); ok && p != nil {
			return *p
		}
		return ""
	}
	pInt64 := func(key string) (int64, bool) {
		v := empresa[key]
		switch n := v.(type) {
		case *int64:
			if n == nil {
				return 0, false
			}
			return *n, true
		case int64:
			return n, true
		case int:
			return int64(n), true
		case float64:
			return int64(n), true
		case string:
			if n == "" {
				return 0, false
			}
			parsed, err := strconv.ParseInt(n, 10, 64)
			if err != nil {
				return 0, false
			}
			return parsed, true
		default:
			return 0, false
		}
	}
	simNao := func(key string) string {
		v, _ := empresa[key].(string)
		if v == "Sim" {
			return "S"
		}
		return "N"
	}

	param := map[string]any{
		"codigo_cliente_integracao": fmt.Sprintf("%d", empresa["codigo_integracao"]),
		"optante_simples_nacional":  simNao("simples_nacional"),
		"produtor_rural":            simNao("produtor_rural"),
		"contribuinte":              simNao("contribuinte"),
		"bloquear_faturamento":      simNao("bloquear_faturamento"),
	}

	addPtr := func(src, dst string) {
		if v := strings.TrimSpace(pStr(src)); v != "" {
			param[dst] = v
		}
	}
	addRaw := func(src, dst string) {
		if v, ok := empresa[src].(string); ok {
			v = strings.TrimSpace(v)
			if v != "" {
				param[dst] = v
			}
		}
	}

	for src, dst := range map[string]string{
		"RAZAO_SOCIAL":      "razao_social",
		"CNPJ":              "cnpj_cpf",
		"nome_fantasia":     "nome_fantasia",
		"contato":           "contato",
		"numero":            "endereco_numero",
		"Bairro":            "bairro",
		"Estado":            "estado",
		"Cidade":            "cidade",
		"CEP":               "cep",
		"ddd_telefone2":     "telefone2_ddd",
		"telefone2":         "telefone2_numero",
		"ddd_fax":           "fax_ddd",
		"fax":               "fax_numero",
		"EMAILS":            "email",
		"siteCliente":       "homepage",
		"INS_ESTADUAL":      "inscricao_estadual",
		"INS_MUNICIPAL":     "inscricao_municipal",
		"inscricao_suframa": "inscricao_suframa",
		"tipo_atividade":    "tipo_atividade",
		"CNAE":              "cnae",
		"OBS":               "observacao",
		"txt_restricao":     "obs_detalhadas",
		"limite_credito":    "valor_limite_credito",
	} {
		addPtr(src, dst)
	}

	addRaw("ENDERECO", "endereco")
	addRaw("complemento", "complemento")

	if pais := strings.TrimSpace(pStr("Pais")); pais != "" && len(pais) <= 4 {
		param["codigo_pais"] = pais
	}

	endEnt := map[string]any{}
	addEndPtr := func(src, dst string) {
		if v := strings.TrimSpace(pStr(src)); v != "" {
			endEnt[dst] = v
		}
	}
	addEndRaw := func(src, dst string) {
		if v, ok := empresa[src].(string); ok {
			v = strings.TrimSpace(v)
			if v != "" {
				endEnt[dst] = v
			}
		}
	}

	for src, dst := range map[string]string{
		"cnpj_entrega":     "entCnpjCpf",
		"nome_entrega":     "entRazaoSocial",
		"numero_entrega":   "entNumero",
		"estado_entrega":   "entEstado",
		"cidade_entrega":   "entCidade",
		"cep_entrega":      "entCEP",
		"ie_entrega":       "entIE",
		"telefone_entrega": "entTelefone",
	} {
		addEndPtr(src, dst)
	}

	addEndRaw("endereco_entrega", "entEndereco")
	addEndRaw("complemento_entrega", "entComplemento")

	if len(endEnt) > 0 {
		param["enderecoEntrega"] = endEnt
	}

	recomendacoes := map[string]any{
		"gerar_boletos": simNao("gerar_boleto"),
	}
	if v := pStr("email_nf"); v != "" {
		recomendacoes["email_fatura"] = v
	}
	if vendedor, ok := pInt64("vendedor"); ok {
		recomendacoes["codigo_vendedor"] = vendedor
	}
	param["recomendacoes"] = recomendacoes

	bodyPayload := map[string]any{
		"call":       "IncluirCliente",
		"param":      []any{param},
		"app_key":    s.AppKey,
		"app_secret": s.AppSecret,
	}

	payloadBytes, err := json.Marshal(bodyPayload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[ImportarEmpresa] Status: %s | Resposta: %s", resp.Status, string(respBody))

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta da Omie: %w", err)
	}

	if faultMsg, ok := result["faultstring"]; ok {
		return result, fmt.Errorf("erro Omie: %v", faultMsg)
	}

	codigoOmieRaw, ok := result["codigo_cliente_omie"]
	if !ok {
		return result, fmt.Errorf("resposta Omie sem codigo_cliente_omie")
	}

	codigoOmie, ok := codigoOmieRaw.(float64)
	if !ok {
		return result, fmt.Errorf("tipo inválido de codigo_cliente_omie: %T", codigoOmieRaw)
	}

	if err := repositories.UpsertRelacaoCliente(db, fmt.Sprintf("%d", req.Id), int64(codigoOmie), pStr("CNPJ")); err != nil {
		return result, fmt.Errorf("cliente importado na Omie, mas falhou ao gravar relacao local: %w", err)
	}

	return result, nil
}
