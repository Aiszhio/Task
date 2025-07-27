package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func NewMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)

		logger.Info(
			"method", ctx.Request.Method,
			"url", ctx.Request.URL.String(),
			slog.Int("status", ctx.Writer.Status()),
			slog.Int("latency_ms", int(duration.Milliseconds())),
			"client_ip", ctx.ClientIP(),
		)

		if len(ctx.Errors) > 0 {
			for _, err := range ctx.Errors {
				logger.Error("handle error",
					"method", ctx.Request.Method,
					"url", ctx.Request.URL.String(),
					"err", err,
				)
			}
		}
	}
}
