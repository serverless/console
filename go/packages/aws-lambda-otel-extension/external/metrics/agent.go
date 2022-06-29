package metrics

import (
	"aws-lambda-otel-extension/external/lib"
	"aws-lambda-otel-extension/external/reporter"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const OTEL_SERVER_PORT = "2772"

type TelemetryPayload struct {
	RecordType string `json:"recordType"`
}

// InternalHttpClient is a simple client for the OTEL server
type InternalHttpListener struct {
	httpServer *http.Server
	// metricsQueue is a synchronous queue and is used to put the received metrics to be consumed later (see main)
	reportAgent        *reporter.HttpClient
	currentRequestData *reporter.CurrentRequestData
	logger             *lib.Logger
}

func NewInternalHttpListener(reportAgent *reporter.HttpClient, currentRequestData *reporter.CurrentRequestData) *InternalHttpListener {
	return &InternalHttpListener{
		httpServer:         nil,
		reportAgent:        reportAgent,
		logger:             lib.NewLogger(),
		currentRequestData: currentRequestData,
	}
}

// Start initiates the server in a goroutine
func (l *InternalHttpListener) Start() bool {
	address := fmt.Sprintf("localhost:%s", OTEL_SERVER_PORT)
	mux := http.NewServeMux()
	mux.HandleFunc("/", l.http_handler)
	l.httpServer = &http.Server{Addr: address, Handler: mux}
	go func() {
		l.logger.Info("Serving internal agent on address: " + address)
		err := l.httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			l.logger.Error("Unexpected stop on Http Server", zap.Error(err))
			l.Shutdown()
		} else {
			l.logger.Debug("Http Server closed", zap.Error(err))
		}
	}()
	return true
}

// http_handler handles the requests coming from the Internal API.
// Everytime Internal API sends metrics, this function will read the metrics from the response body
// and put them into a synchronous queue to be read by the main goroutine.
// Logging or printing besides the error cases below is not recommended if you have subscribed to receive extension metrics.
// Otherwise, logging here will cause Logs API to send new logs for the printed lines which will create an infinite loop.
func (l *InternalHttpListener) http_handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.logger.Error("Error reading body", zap.Error(err))
		return
	}

	// fmt.Println("Internal API event received:", string(body))

	payload, err := reporter.ParseInternalPayload(body)
	if err != nil {
		l.logger.Error("Error parsing payload", zap.Error(err))
		return
	}

	switch payload.RecordType {
	case "eventData":
		eventData, err := reporter.ParseEventDataPayload(payload.Record)
		if err != nil {
			l.logger.Error("Error parsing payload", zap.Error(err))
			return
		}
		l.currentRequestData.SetUniqueName("request")
		l.currentRequestData.SetEventData(eventData)
		if eventData.RequestEventPayload != nil {
			l.reportAgent.PostRequest(*eventData.RequestEventPayload)
		}
		return

	case "telemetryData":
		telemetryData, err := reporter.ParseTelemetryDataPayload(payload.Record)
		if err != nil {
			l.logger.Error("Error parsing payload", zap.Error(err))
			return
		}
		if telemetryData.ResponseEventPayload != nil {
			l.reportAgent.PostResponse(*telemetryData.ResponseEventPayload)
		}
		metrics := reporter.CreateMetricsPayload(telemetryData.RequestID, telemetryData.Function, nil)
		data, err := proto.Marshal(metrics)
		if err != nil {
			l.logger.Error("Error marshalling metrics", zap.Error(err))
		}
		l.reportAgent.PostMetric(data)

		traces, err := reporter.CreateTracePayload(telemetryData.RequestID, telemetryData.Function, telemetryData.Traces)
		if err != nil {
			l.logger.Error("Error creating traces", zap.Error(err))
			return
		}

		for _, trace := range traces {
			data, err = proto.Marshal(trace)
			if err != nil {
				l.logger.Error("Error marshalling traces", zap.Error(err))
				return
			}

			l.reportAgent.PostTrace(data)
		}

	default:
		l.logger.Error("Unknown payload type", zap.String("payload", string(body)))
	}
}

// Shutdown terminates the HTTP server
func (l *InternalHttpListener) Shutdown() {
	if l.httpServer != nil {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		err := l.httpServer.Shutdown(ctx)
		if err != nil {
			l.logger.Error("Failed to shutdown http server gracefully ", zap.Error(err))
		} else {
			l.httpServer = nil
		}
	}
}