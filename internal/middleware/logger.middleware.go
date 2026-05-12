package middleware

import (
	"bytes"
	"database/sql"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogger(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		blw := &responseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		c.Next()

		metodo := c.Request.Method
		rota := c.FullPath()
		if rota == "" {
			rota = c.Request.URL.Path
		}
		statusCode := c.Writer.Status()
		ipOrigem := c.ClientIP()
		corpoRequisicao := string(reqBody)
		corpoResposta := blw.body.String()
		duracaoMs := time.Since(start).Milliseconds()

		go func() {
			_, err := db.Exec(
				`INSERT INTO amm_omie_logs
					(metodo, rota, status_code, ip_origem, corpo_requisicao, corpo_resposta, duracao_ms)
				VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`,
				metodo,
				rota,
				statusCode,
				ipOrigem,
				corpoRequisicao,
				corpoResposta,
				duracaoMs,
			)
			if err != nil {
				log.Printf("[middleware] erro ao salvar log: %v", err)
			}
		}()
	}
}
