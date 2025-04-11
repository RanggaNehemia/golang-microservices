package tracing

import (
	"context"

	"github.com/RanggaNehemia/golang-microservices/auth-service/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func InitTracer() func() {
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		utils.Logger.Fatal("failed to initialize stdouttrace exporter", zap.Error(err))
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
	)

	otel.SetTracerProvider(tp)
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			utils.Logger.Fatal("Error shutting down tracer provider", zap.Error(err))
		}
	}
}
