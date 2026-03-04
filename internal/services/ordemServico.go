package services

import (
	"example/web-service-gin/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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


func (s *OmieService) FaturarOrdemServico(req models.FaturaOrdemServicoRequest) (string, error) {
	url := s.BaseURL + "/api/v1/servicos/osp/"
	payload := strings.NewReader(`{
		"call":"FaturarOS",
		"param":[{
			"cCodIntOS":"` + req.CCodIntOS + `",
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