package services

import (
	"example/web-service-gin/internal/models"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Tenho que ver para deixar isso OFUSCADO, appKey e AppSecret
type OmieService struct {
	AppKey    string
	AppSecret string
	BaseURL   string
}

func NewOmieService(appKey, appSecret, baseURL string) *OmieService {
	return &OmieService{
		AppKey:    appKey,
		AppSecret: appSecret,
		BaseURL:   baseURL,
	}
}
func (s *OmieService) ListarServicos() (string, error) {
	url := s.BaseURL + "/api/v1/servicos/servico/"
	payload := strings.NewReader(`{
				"call": "ListarCadastroServico",
				"param":[{"nPagina":1,"nRegPorPagina":20}],
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
func (s *OmieService) CadastrarServico(req models.ServicoRequest) (string, error) {
	url := s.BaseURL + "/api/v1/servicos/servico/"

	payload := strings.NewReader(`{
        "call": "IncluirCadastroServico",
        "param": [{
            "intIncluir": {"cCodIntServ": "` + req.CodInterno + `"},
            "descricao": {"cDescrCompleta": "` + req.Descricao + `"},
            "cabecalho": {
                "cDescricao": "` + req.Descricao + `",
                "cCodigo": "` + req.CCodigo + `",
                "nPrecoUnit": ` + req.PrecoUnitario + `
            }
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

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Resposta:", string(body))

	return string(body), nil
}
