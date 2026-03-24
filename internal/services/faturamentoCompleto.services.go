package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"example/web-service-gin/internal/models"
)

type FluxoFaturamento struct {
	OsIncluida   chan models.WebhookOsIncluidaResponse
	ContaReceber chan models.WebhookContaReceberResponseInclude
	BoletoGerado chan models.WebhookBoletoGeradoResponse
	numeroOs     string // preenchido após OrdemServico.Incluida
}

type FaturamentoOrquestrador struct {
	mu     sync.Mutex
	fluxos map[string]*FluxoFaturamento
}

var orquestrador = &FaturamentoOrquestrador{
	fluxos: make(map[string]*FluxoFaturamento),
}

func GetOrquestrador() *FaturamentoOrquestrador {
	return orquestrador
}

func (o *FaturamentoOrquestrador) registrar(codIntOS string) *FluxoFaturamento {
	fluxo := &FluxoFaturamento{
		OsIncluida:   make(chan models.WebhookOsIncluidaResponse, 1),
		ContaReceber: make(chan models.WebhookContaReceberResponseInclude, 1),
		BoletoGerado: make(chan models.WebhookBoletoGeradoResponse, 1),
	}
	o.mu.Lock()
	o.fluxos[codIntOS] = fluxo
	o.mu.Unlock()
	return fluxo
}

func (o *FaturamentoOrquestrador) remover(codIntOS string) {
	o.mu.Lock()
	delete(o.fluxos, codIntOS)
	o.mu.Unlock()
}

func (o *FaturamentoOrquestrador) NotificarOsIncluida(data models.WebhookOsIncluidaResponse) bool {
	o.mu.Lock()
	fluxo, ok := o.fluxos[data.CodigoIntegra]
	if ok {
		fluxo.numeroOs = data.NumeroOs
	}
	o.mu.Unlock()

	if !ok {
		return false
	}
	select {
	case fluxo.OsIncluida <- data:
	default:
	}
	return true
}

func (o *FaturamentoOrquestrador) NotificarContaReceber(data models.WebhookContaReceberResponseInclude) bool {
	o.mu.Lock()
	var fluxo *FluxoFaturamento
	for _, f := range o.fluxos {
		if f.numeroOs == data.NumeroPedido {
			fluxo = f
			break
		}
	}
	o.mu.Unlock()

	if fluxo == nil {
		return false
	}
	select {
	case fluxo.ContaReceber <- data:
	default:
	}
	return true
}

func (o *FaturamentoOrquestrador) NotificarBoletoGerado(data models.WebhookBoletoGeradoResponse) bool {
	o.mu.Lock()
	var fluxo *FluxoFaturamento
	for _, f := range o.fluxos {
		if f.numeroOs == data.NumeroPedido {
			fluxo = f
			break
		}
	}
	o.mu.Unlock()

	if fluxo == nil {
		return false
	}
	select {
	case fluxo.BoletoGerado <- data:
	default:
	}
	return true
}

func (s *OmieService) CriarFaturamentoCompleto(req models.OrdemServicoRequest) (map[string]any, error) {
	codIntOS := req.Cabecalho.CCodIntOS

	orq := GetOrquestrador()
	fluxo := orq.registrar(codIntOS)
	defer orq.remover(codIntOS)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fmt.Printf("[FaturamentoCompleto] Criando OS: %s\n", codIntOS)
	if _, err := s.CriarOrdemServico(req); err != nil {
		return nil, fmt.Errorf("erro ao criar OS: %w", err)
	}

	fmt.Println("[FaturamentoCompleto] Aguardando webhook OrdemServico.Incluida...")
	var osIncluida models.WebhookOsIncluidaResponse
	select {
	case osIncluida = <-fluxo.OsIncluida:
		fmt.Printf("[FaturamentoCompleto] OS incluída — nCodOS: %d, numeroOs: %s\n", osIncluida.IdOs, osIncluida.NumeroOs)
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout aguardando webhook OrdemServico.Incluida")
	}

	fmt.Printf("[FaturamentoCompleto] Faturando OS: %d\n", osIncluida.IdOs)
	if _, err := s.FaturarOrdemServico(models.FaturaOrdemServicoRequest{
		NCodOS: int(osIncluida.IdOs),
	}); err != nil {
		return nil, fmt.Errorf("erro ao faturar OS: %w", err)
	}

	fmt.Println("[FaturamentoCompleto] Aguardando webhook Financas.ContaReceber.Incluido...")
	var contaReceber models.WebhookContaReceberResponseInclude
	select {
	case contaReceber = <-fluxo.ContaReceber:
		fmt.Printf("[FaturamentoCompleto] Conta a receber incluída — CodigoConta: %d\n", contaReceber.CodigoConta)
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout aguardando webhook Financas.ContaReceber.Incluido")
	}

	fmt.Printf("[FaturamentoCompleto] Gerando boleto para conta: %d\n", contaReceber.CodigoConta)
	resultadoBoleto, err := s.GerarBoletoConta(models.GerarBoletoConta{
		NCodTitulo: contaReceber.CodigoConta,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar boleto: %w", err)
	}

	fmt.Println("[FaturamentoCompleto] Aguardando webhook Financas.ContaReceber.BoletoGerado...")
	var boletoGerado models.WebhookBoletoGeradoResponse
	select {
	case boletoGerado = <-fluxo.BoletoGerado:
		fmt.Printf("[FaturamentoCompleto] Boleto gerado — CodigoBarras: %s\n", boletoGerado.CodigoBarras)
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout aguardando webhook Financas.ContaReceber.BoletoGerado")
	}

	return map[string]any{
		"os_incluida":    osIncluida,
		"conta_receber":  contaReceber,
		"boleto_gerado":  boletoGerado,
		"resultado_omie": resultadoBoleto,
	}, nil
}
