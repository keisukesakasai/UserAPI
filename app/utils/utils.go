package utils

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("UserAPI-utils")

func LoggerAndCreateSpan(c *gin.Context, msg string) trace.Span {
	_, span := tracer.Start(c.Request.Context(), msg)
	span.SetAttributes(
		attribute.Int("status", c.Writer.Status()),
		attribute.String("method", c.Request.Method),
		attribute.String("client_ip", c.ClientIP()),
		attribute.String("message", msg),
		attribute.String("span_id", span.SpanContext().SpanID().String()),
		attribute.String("trace_id", span.SpanContext().TraceID().String()),
	)

	start := time.Now()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	defer logger.Sync()
	logger.Info("Logger",
		zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()),
		zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		zap.Duration("elapsed", time.Since(start)),
		zap.String("message", msg),
		zap.String("span_id", span.SpanContext().SpanID().String()),
		zap.String("trace_id", span.SpanContext().TraceID().String()),
	)

	return span
}
