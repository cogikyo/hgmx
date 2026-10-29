module github.com/cogikyo/hgmx/telemetry/otlptest

go 1.26.0

require (
	github.com/cogikyo/hgmx v0.0.0
	go.opentelemetry.io/proto/otlp v1.11.1
	google.golang.org/protobuf v1.36.12
)

replace github.com/cogikyo/hgmx => ../..
