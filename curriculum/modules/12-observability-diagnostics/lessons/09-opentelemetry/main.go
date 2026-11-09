package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func newExporter() (*stdouttrace.Exporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		"https://opentelemetry.io/schemas/1.26.0",
		attribute.String("service.name", "order-service"),
		attribute.String("service.version", "1.0.0"),
		attribute.String("environment", "development"),
	)
}

func newTraceProvider(exp *stdouttrace.Exporter) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(newResource()),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exp)),
	)
}

func main() {
	exp, err := newExporter()
	if err != nil {
		panic(err)
	}

	tp := newTraceProvider(exp)
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	tracer := otel.Tracer("order-service")
	ctx := context.Background()

	ctx, span := tracer.Start(ctx, "create_order")
	span.SetAttributes(
		attribute.String("order.id", "ORD-12345"),
		attribute.Float64("order.amount", 99.50),
		attribute.String("order.currency", "USD"),
	)
	span.AddEvent("order.validation_passed")

	if err := processPayment(ctx, tracer); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "payment failed: "+err.Error())
		fmt.Println("Payment error:", err)
	} else {
		span.SetStatus(codes.Ok, "order created successfully")
	}

	span.End()
	time.Sleep(100 * time.Millisecond)
}

func processPayment(ctx context.Context, tracer trace.Tracer) error {
	ctx, span := tracer.Start(ctx, "process_payment")
	defer span.End()

	span.SetAttributes(attribute.String("payment.method", "credit_card"))
	span.AddEvent("charging_card")

	time.Sleep(50 * time.Millisecond)

	if time.Now().Unix()%2 == 0 {
		return errors.New("insufficient funds")
	}
	return nil
}
