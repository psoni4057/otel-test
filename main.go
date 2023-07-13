package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	"go.opentelemetry.io/otel/attribute"
	otelMetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

type RequestError struct {
	StatusCode string
	Status     string
}

func New() *OpenTelemetry {
	return &OpenTelemetry{}
}

type OpenTelemetry struct {
	provider *metric.MeterProvider

	Meter otelMetric.Meter
}

func (openTelemetry *OpenTelemetry) StartOT() {
	fmt.Println("START: StartOT")
	exporter, err := mexporter.New(mexporter.WithProjectID("techyon-dev-main"))
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	//sdk/metric
	openTelemetry.provider = metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(3*time.Minute))),
	)
	//sdk/metric
	openTelemetry.Meter = openTelemetry.provider.Meter("test")
	fmt.Println("END: StartOT")
}

func (openTelemetry *OpenTelemetry) StopOT() {
	fmt.Println("START: StopOT")
	if openTelemetry.provider != nil {
		openTelemetry.provider.Shutdown(context.Background())
	}
	fmt.Println("END: StopOT")
}

// sync counter;inserting code into applications that will update a value each time it is executed
func myFunc(opentelemetry *OpenTelemetry) {
	fmt.Println("START: myFunc")
	meterCounter1, error := opentelemetry.Meter.Int64Counter("my-metric/test1")
	if error != nil {
		log.Fatalf("Failed to create counter: %v", error)
	}
	meterCounter2, error2 := opentelemetry.Meter.Int64Counter("my-metric/test2")
	if error2 != nil {
		log.Fatalf("Failed to create counter: %v", error2)
	}
	ctx := context.Background()

	//code
	a := reflect.ValueOf(RequestError{"404", "error 4xx"})
	property1 := "StatusCode"
	property2 := "Status"
	a1 := a.FieldByName(property1)
	a2 := a.FieldByName(property2)

	fmt.Println("adding labels key-value")
	labels := []attribute.KeyValue{
		attribute.String("techyon_biodome", "MY_BIODOME"),
		attribute.String("techyon_tenant", "MY_TENANT"),
		attribute.String("techyon_habitat", "MY_HABITAT"),
		attribute.String("techyon_workload", "my-dispatcher"),
		attribute.String("techyon_team", "MY_TEAM"),
		attribute.String("pod", "MY_HOST_NAME"),
		attribute.String("subscription_id", "MY_SUBSCRIPTION"),
		attribute.String("content_type", "MY_CONTENT_TYPE"),
		attribute.String("status", a2.String()),
		attribute.String("status_code", a1.String()),
	}
	meterCounter1.Add(ctx, 1, otelMetric.WithAttributes(labels...))
	labels2 := []attribute.KeyValue{
		attribute.String("techyon_biodome", "MY_BIODOME2"),
		attribute.String("techyon_tenant", "MY_TENANT2"),
		attribute.String("techyon_habitat", "MY_HABITAT2"),
		attribute.String("techyon_workload", "my-dispatcher2"),
		attribute.String("techyon_team", "MY_TEAM2"),
		attribute.String("pod", "MY_HOST_NAME2"),
		attribute.String("subscription_id", "MY_SUBSCRIPTION2"),
		attribute.String("content_type", "MY_CONTENT_TYPE2"),
		attribute.String("status", a2.String()),
		attribute.String("status_code", a1.String()),
	}
	meterCounter2.Add(ctx, 1, otelMetric.WithAttributes(labels2...))
	fmt.Println("added labels")
	fmt.Println("END: myFunc")
}

func main() {
	opentelemetry := New()
	opentelemetry.StartOT()
	myFunc(opentelemetry)
	opentelemetry.StopOT()
}
