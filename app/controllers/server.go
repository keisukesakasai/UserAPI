package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"
	"userapi/config"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("UserAPI"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var tracerProvider *sdktrace.TracerProvider

	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		// stdouttrace.WithWriter(os.Stderr),
		stdouttrace.WithWriter(io.Discard),
	)

	/*
		conn, err := grpc.DialContext(ctx, "otel-collector-collector.tracing.svc.cluster.local:4318", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
		}

		// Set up a trace exporter
		traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}
	*/

	// idg := xray.NewIDGenerator()

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
		// sdktrace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

//-- otelcollecotr
var tracer = otel.Tracer("UserAPI-controllers")

func StartMainServer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Otel Collecotor への接続設定
	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	// router 設定
	r := gin.New()

	r.Use(otelgin.Middleware("UserAPI-server"))

	//--- handler 設定
	r.POST("/createUser", createUser)
	r.POST("/getUserByEmail", getUserByEmail)

	r.POST("/encrypt", Encrypt)

	r.Run(":" + config.Config.Port)
}
