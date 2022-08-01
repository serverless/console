package reporter

import (
	"aws-lambda-otel-extension/external/types"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricspb "go.opentelemetry.io/proto/otlp/metrics/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

const MaxSpansPerBatch = 100

func ParseInternalPayload(data []byte) (*types.RecordPayload, error) {
	var payload *types.RecordPayload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func ParseEventDataPayload(data json.RawMessage) (*types.EventDataPayload, error) {
	var payload *types.EventDataPayload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func ParseTelemetryDataPayload(data json.RawMessage) (*types.TelemetryDataPayload, error) {
	var payload *types.TelemetryDataPayload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func createMetricAttributes(fun map[string]interface{}, record *types.PlatformObjectRecord) []*commonpb.KeyValue {
	var attributes []*commonpb.KeyValue

	// for loop to iterate MeasureAttributes
	for _, ra := range MeasureAttributes {
		if ra.value != "" {
			attributes = append(attributes, &commonpb.KeyValue{
				Key:   ra.key,
				Value: getAnyValue(ra.value),
			})
		} else {
			if ra.source != "" {
				if val, ok := fun[ra.source]; ok {
					attributes = append(attributes, &commonpb.KeyValue{
						Key:   ra.key,
						Value: getAnyValue(val),
					})
				}
			}
		}
	}

	attributes = append(attributes, &commonpb.KeyValue{
		Key:   "faas.error",
		Value: getAnyValue(fun["error"]),
	})

	if record != nil {
		attributes = append(attributes, &commonpb.KeyValue{
			Key:   "faas.duration",
			Value: getAnyValue(record.ReportLogItem.DurationMs),
		})
		attributes = append(attributes, &commonpb.KeyValue{
			Key:   "faas.billed_duration",
			Value: getAnyValue(record.ReportLogItem.BilledDurationMs),
		})
		attributes = append(attributes, &commonpb.KeyValue{
			Key:   "faas.max_memory_used_mb",
			Value: getAnyValue(record.ReportLogItem.MaxMemoryUsedMB),
		})
	}

	// Add timeout info
	hasErrorTimeout := false
	et, ok := fun["errorTimeout"].(string)
	if ok {
		hasErrorTimeout = et == "timeout"
	}

	attributes = append(attributes, &commonpb.KeyValue{
		Key:   "faas.error_timeout",
		Value: getAnyValue(hasErrorTimeout),
	})

	return attributes
}

func createResourceAttributes(fun map[string]interface{}) []*commonpb.KeyValue {
	var attributes []*commonpb.KeyValue

	// for loop to iterate ResourceAttributes
	for _, ra := range ResourceAttributes {
		if ra.value != "" {
			attributes = append(attributes, &commonpb.KeyValue{
				Key:   ra.key,
				Value: getAnyValue(ra.value),
			})
		} else {
			if ra.source != "" {
				if val, ok := fun[ra.source]; ok {
					attributes = append(attributes, &commonpb.KeyValue{
						Key:   ra.key,
						Value: getAnyValue(val),
					})
				}
			}
		}
	}

	return attributes
}

func getTimeUnixNanoInterval(record map[string]interface{}) (uint64, uint64) {
	startTimeUnixNano := getTimeUnixNano(record["startTime"])
	if startTimeUnixNano == 0 {
		fmt.Printf(">> startTime is not valid (%v)", record["startTime"])
		return 0, 0
	}
	endTimeUnixNano := getTimeUnixNano(record["endTime"])
	return startTimeUnixNano, endTimeUnixNano
}

func createHistogramMetric(count uint64, sum float64, record map[string]interface{}, attributes []*commonpb.KeyValue) *metricspb.Metric_Histogram {
	startTime, endTime := getTimeUnixNanoInterval(record)
	return &metricspb.Metric_Histogram{
		Histogram: &metricspb.Histogram{
			DataPoints: []*metricspb.HistogramDataPoint{
				{
					StartTimeUnixNano: startTime,
					TimeUnixNano:      endTime,
					Count:             count,
					Sum:               sum,
					BucketCounts:      []uint64{1, 0},
					// ExplicitBounds:    []float64{math.Inf(1)},
					Attributes: attributes,
				},
			},
		},
	}
}

func createCountMetric(count uint64, asInt int64, record map[string]interface{}, attributes []*commonpb.KeyValue) *metricspb.Metric_Sum {
	startTime, endTime := getTimeUnixNanoInterval(record)
	return &metricspb.Metric_Sum{
		Sum: &metricspb.Sum{
			DataPoints: []*metricspb.NumberDataPoint{
				{
					StartTimeUnixNano: startTime,
					TimeUnixNano:      endTime,
					Attributes:        attributes,
					Value:             &metricspb.NumberDataPoint_AsInt{AsInt: asInt},
				},
			},
		},
	}
}

func CreateMetricsPayload(requestId string, fun map[string]interface{}, record *types.PlatformObjectRecord) *metricspb.MetricsData {

	metricAttributes := createMetricAttributes(fun, record)

	metricAttributes = append(metricAttributes, &commonpb.KeyValue{
		Key:   "faas.execution",
		Value: getAnyValue(requestId),
	})

	if v, ok := fun["httpStatusCode"]; ok {
		metricAttributes = append(metricAttributes, &commonpb.KeyValue{
			Key:   "http.status_code",
			Value: getAnyValue(v),
		})
	}

	if v, ok := fun["httpPath"]; ok {
		metricAttributes = append(metricAttributes, &commonpb.KeyValue{
			Key:   "http.path",
			Value: getAnyValue(v),
		})
	}

	metrics := []*metricspb.Metric{}

	if record == nil {
		metrics = append(metrics, &metricspb.Metric{
			Name: "faas.invoke",
			Unit: "1",
			Data: createCountMetric(1, 1, fun, metricAttributes),
		})

	} else {
		metrics = append(metrics, &metricspb.Metric{
			Name: "faas.duration",
			Unit: "1",
			Data: createHistogramMetric(1, record.ReportLogItem.DurationMs, fun, metricAttributes),
		})

		memory := (float64(record.ReportLogItem.MaxMemoryUsedMB) / float64(record.ReportLogItem.MemorySizeMB)) * 100

		metrics = append(metrics, &metricspb.Metric{
			Name: "faas.memory",
			Unit: "1",
			Data: createHistogramMetric(1, memory, fun, metricAttributes),
		})

		if record.ReportLogItem.InitDurationMs > 0 {
			metrics = append(metrics, &metricspb.Metric{
				Name: "faas.coldstart_duration",
				Unit: "1",
				Data: createHistogramMetric(1, record.ReportLogItem.InitDurationMs, fun, metricAttributes),
			})
		}
	}

	return &metricspb.MetricsData{
		ResourceMetrics: []*metricspb.ResourceMetrics{
			{
				Resource: &resourcepb.Resource{
					Attributes: createResourceAttributes(fun),
				},
				InstrumentationLibraryMetrics: []*metricspb.InstrumentationLibraryMetrics{
					{
						InstrumentationLibrary: &commonpb.InstrumentationLibrary{
							Name:    "serverless-meter",
							Version: "1.0.0",
						},
						Metrics: metrics,
					},
				},
			},
		},
	}
}

// Traces logic

func batchOverflowSpans(traces *tracepb.TracesData) *tracepb.TracesData {
	var batches []*tracepb.InstrumentationLibrarySpans

	for _, libSpans := range traces.ResourceSpans[0].InstrumentationLibrarySpans {
		var cutSpans []*tracepb.Span
		lenSpans := len(libSpans.Spans)
		iSpan := 0
		for (lenSpans - iSpan) > 0 {
			if (lenSpans - iSpan) > MaxSpansPerBatch {
				cutSpans = libSpans.Spans[iSpan : iSpan+MaxSpansPerBatch]
				iSpan += MaxSpansPerBatch
			} else {
				cutSpans = libSpans.Spans[iSpan:]
				iSpan += len(cutSpans)
			}
			batches = append(batches, &tracepb.InstrumentationLibrarySpans{
				InstrumentationLibrary: libSpans.InstrumentationLibrary,
				Spans:                  cutSpans,
			})

		}
	}

	return &tracepb.TracesData{
		ResourceSpans: []*tracepb.ResourceSpans{
			{
				Resource:                    traces.ResourceSpans[0].Resource,
				InstrumentationLibrarySpans: batches,
			},
		},
	}
}

func CreateTracePayload(requestId string, fun map[string]interface{}, traces *types.Traces) (*tracepb.TracesData, error) {

	if traces == nil {
		return nil, nil
	}
	metrics := createMetricAttributes(fun, nil)

	if len(traces.ResourceSpans) == 0 {
		return nil, errors.New("no resource spans found")
	}

	instLibSpans := make([]*tracepb.InstrumentationLibrarySpans, len(traces.ResourceSpans[0].InstrumentationLibrarySpans))

	for libIndex, librarySpans := range traces.ResourceSpans[0].InstrumentationLibrarySpans {
		if librarySpans.Spans == nil {
			return nil, errors.New("no spans found")
		}
		cSpans := make([]*tracepb.Span, len(librarySpans.Spans))

		for spanIndex, span := range librarySpans.Spans {
			// finally convert spans to key-value format
			var attribs []*commonpb.KeyValue
			existingAttribs := map[string]bool{}
			for k, v := range span.Attributes {
				attribs = append(attribs, &commonpb.KeyValue{
					Key:   k,
					Value: getAnyValue(v),
				})
				existingAttribs[k] = true
			}
			for _, kv := range metrics {
				if _, ok := existingAttribs[kv.Key]; !ok {
					attribs = append(attribs, kv)
				}
			}

			traceId, _ := b64.StdEncoding.DecodeString(span.TraceID)
			spanId, _ := b64.StdEncoding.DecodeString(span.SpanID)
			parentSpanId, _ := b64.StdEncoding.DecodeString(span.ParentSpanID)
			startTime, _ := strconv.ParseInt(span.StartTimeUnixNano, 10, 64)
			endTime, _ := strconv.ParseInt(span.EndTimeUnixNano, 10, 64)

			cSpans[spanIndex] = &tracepb.Span{
				Attributes:        attribs,
				TraceId:           traceId,
				SpanId:            spanId,
				ParentSpanId:      parentSpanId,
				Name:              span.Name,
				StartTimeUnixNano: uint64(startTime),
				EndTimeUnixNano:   uint64(endTime),
			}
		}

		instLibSpans[libIndex] = &tracepb.InstrumentationLibrarySpans{

			InstrumentationLibrary: &commonpb.InstrumentationLibrary{
				Name:    librarySpans.InstrumentationLibrary.Name,
				Version: librarySpans.InstrumentationLibrary.Version,
			},
			Spans: cSpans,
		}
	}

	resourceSpans := []*tracepb.ResourceSpans{
		{
			Resource: &resourcepb.Resource{
				Attributes: createResourceAttributes(fun),
			},
			InstrumentationLibrarySpans: instLibSpans,
		},
	}
	return batchOverflowSpans(&tracepb.TracesData{
		ResourceSpans: resourceSpans,
	}), nil

}
