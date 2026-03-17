package workers

import (
	"errors"
	"log"
	"sync"
	"time"

	"example/web-service-gin/internal/models"
)

var ErrFilaCheia = errors.New("fila de webhook cheia")

type Processor interface {
	ProcessarWebhookOsFaturada(data models.WebhookOsFaturadaResponse) (int, error)
	ProcessarWebhookContaReceber(data models.WebhookContaReceberResponseInclude) (int, error)
	ProcessarWebhookOsIncluida(data models.WebhookOsIncluidaResponse) (int, error)
}

type TipoJob string

const (
	JobOsFaturada   TipoJob = "os_faturada"
	JobContaReceber TipoJob = "conta_receber_include"
	JobOsIncluida   TipoJob = "os_incluida"
)

type WebhookJob struct {
	Tipo TipoJob

	OsFaturada   *models.WebhookOsFaturadaResponse
	OsIncluida   *models.WebhookOsIncluidaResponse
	ContaReceber *models.WebhookContaReceberResponseInclude
}

type WebhookWorkerPool struct {
	processor Processor
	jobs      chan WebhookJob
	wg        sync.WaitGroup
}

func NewWebhookWorkerPool(processor Processor, totalWorkers, tamanhoFila int) *WebhookWorkerPool {
	p := &WebhookWorkerPool{
		processor: processor,
		jobs:      make(chan WebhookJob, tamanhoFila),
	}

	for i := 0; i < totalWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i + 1)
	}

	return p
}

func (p *WebhookWorkerPool) Enqueue(job WebhookJob) error {
	select {
	case p.jobs <- job:
		return nil
	default:
		return ErrFilaCheia
	}
}

func (p *WebhookWorkerPool) EnqueueWithWait(job WebhookJob) error {
	p.jobs <- job
	return nil
}

func (p *WebhookWorkerPool) QueueSize() int {
	return len(p.jobs)
}

func (p *WebhookWorkerPool) QueueCapacity() int {
	return cap(p.jobs)
}

func (p *WebhookWorkerPool) Shutdown() {
	close(p.jobs)
	p.wg.Wait()
}

func (p *WebhookWorkerPool) worker(id int) {
	defer p.wg.Done()

	for job := range p.jobs {
		var err error

		for tentativa := 1; tentativa <= 3; tentativa++ {
			switch job.Tipo {
			case JobOsFaturada:
				if job.OsFaturada != nil {
					_, err = p.processor.ProcessarWebhookOsFaturada(*job.OsFaturada)
				}
			case JobContaReceber:
				if job.ContaReceber != nil {
					_, err = p.processor.ProcessarWebhookContaReceber(*job.ContaReceber)
				}
			case JobOsIncluida:
				if job.OsIncluida != nil {
					_, err = p.processor.ProcessarWebhookOsIncluida(*job.OsIncluida)
				}
			default:
				err = errors.New("tipo de job inválido")
			}

			if err == nil {
				break
			}

			log.Printf("[worker %d] erro no job %s, tentativa %d: %v", id, job.Tipo, tentativa, err)
			time.Sleep(time.Duration(tentativa) * 300 * time.Millisecond)
		}
	}
}
