package middleware

import (
	"crypto/rand"
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDKey    = "request_id"
	RequestIDHeader = "X-Request-ID"
)

// RequestID genera un UUID v4 único por request y lo almacena en el contexto de Gin
// bajo RequestIDKey. Si el request entrante ya trae un header X-Request-ID
// (ej. inyectado por un CDN o gateway), ese valor se reutiliza para mantener
// consistencia a lo largo de toda la cadena de llamadas.
// El ID final siempre se devuelve en el header de respuesta X-Request-ID.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(RequestIDHeader)
		if rid == "" {
			rid = newUUID()
		}
		c.Set(RequestIDKey, rid)
		c.Header(RequestIDHeader, rid)
		c.Next()
	}
}

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
