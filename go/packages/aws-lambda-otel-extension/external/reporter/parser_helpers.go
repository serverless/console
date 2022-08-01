package reporter

import (
	"aws-lambda-otel-extension/external/lib"
	"fmt"
	"time"

	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

type AttributeHelper struct {
	key    string
	value  string
	source string
}

var ResourceAttributes []AttributeHelper = []AttributeHelper{
	{
		key:    "faas.id",
		source: "computeCustomArn",
	},
	{
		key:    "faas.name",
		source: "functionName",
	},
	{
		key:    "cloud.region",
		source: "computeRegion",
	},
	{
		key:    "sls.application_name",
		source: "sls_app_name",
	},
	{
		key:    "sls.app_uid",
		source: "sls_app_id",
	},
	{
		key:    "service.namespace",
		source: "sls_service_name",
	},
	{
		key:    "deployment.environment",
		source: "sls_stage",
	},
	{
		key:    "sls.org_uid",
		source: "sls_org_id",
	},
	{
		key:    "service.name",
		source: "functionName",
	},
	{
		key:    "telemetry.sdk.language",
		source: "computeRuntime",
	},
	{
		key:    "telemetry.sdk.name",
		value:  "opentelemetry",
		source: "opentelemetry",
	},
	{
		key:    "telemetry.sdk.version",
		value:  "1.0.1",
		source: "version",
	},
	{
		key:    "cloud.provider",
		value:  "aws",
		source: "provider",
	},
	{
		key:    "faas.version",
		source: "computeCustomFunctionVersion",
	},
	{
		key:    "sls.deployment_uid",
		source: "sls_deploy_id",
	},
	{
		key:    "cloud.account.id",
		source: "eventCustomAccountId",
	},
	{
		key:    "cloud.platform",
		value:  "lambda",
		source: "opentelemetry",
	},
	{
		key:    "faas.max_memory",
		source: "computeMemorySize",
	},
	{
		key:    "faas.log_group",
		source: "computeCustomLogGroupName",
	},
	{
		key:    "faas.log_stream_name",
		source: "computeCustomLogStreamName",
	},
	{
		key:    "faas.collector_version",
		value:  fmt.Sprintf("@serverless/aws-lambda-otel-extension@%s", lib.Version),
		source: fmt.Sprintf("@serverless/aws-lambda-otel-extension@%s", lib.Version),
	},
}

var MeasureAttributes []AttributeHelper = []AttributeHelper{
	{
		key:    "faas.coldstart",
		source: "computeIsColdStart",
	},
	{
		key:    "http.method",
		source: "eventCustomHttpMethod",
	},
	{
		key:    "http.raw_path",
		source: "rawHttpPath",
	},
	{
		key:    "http.domain",
		source: "eventCustomDomain",
	},
	{
		key:    "faas.error_stacktrace",
		source: "errorStacktrace",
	},
	{
		key:    "faas.error_message",
		source: "errorMessage",
	},
	{
		key:    "aws.xray.trace_id",
		source: "eventCustomXTraceId",
	},
	{
		key:    "faas.event_type",
		source: "eventType",
	},
	{
		key:    "faas.arch",
		source: "computeCustomEnvArch",
	},
	{
		key:    "faas.api_gateway_request_id",
		source: "eventCustomRequestId",
	},
	{
		key:    "faas.error_timeout",
		source: "timeout",
	},
	{
		key:    "faas.event_source",
		source: "eventSource",
	},
	{
		key:    "faas.api_gateway_app_id",
		source: "eventCustomApiId",
	},
	{
		key:    "faas.request_time_epoch",
		source: "eventCustomRequestTimeEpoch",
	},
	{
		key:    "faas.error_culprit",
		source: "errorCulprit",
	},
	{
		key:    "faas.error_type",
		source: "errorType",
	},
}

// Logs helpers
var LogsMetricAttributeNames map[string]bool = map[string]bool{
	"faas.arch":                   true,
	"faas.api_gateway_request_id": true,
	"faas.event_source":           true,
	"faas.api_gateway_app_id":     true,
}

var LogsResourceAttributeNames map[string]bool = map[string]bool{
	"faas.id":                true,
	"faas.name":              true,
	"cloud.region":           true,
	"sls.app_uid":            true,
	"service.namespace":      true,
	"deployment.environment": true,
	"service.name":           true,
	"telemetry.sdk.language": true,
	"telemetry.sdk.name":     true,
	"telemetry.sdk.version":  true,
	"cloud.provider":         true,
	"cloud.account.id":       true,
	"cloud.platform":         true,
	"faas.collector_version": true,
}

var LogsSeverityNumber map[string]int = map[string]int{
	"TRACE": 1,
	"DEBUG": 5,
	"INFO":  9,
	"WARN":  13,
	"ERROR": 17,
	"FATAL": 21,
}

// Protobuf helpers

func getAnyValue(value interface{}) *commonpb.AnyValue {
	switch value.(type) {
	case string:
		return &commonpb.AnyValue{
			Value: &commonpb.AnyValue_StringValue{
				StringValue: value.(string),
			},
		}
	case int64:
		return &commonpb.AnyValue{
			Value: &commonpb.AnyValue_IntValue{
				IntValue: value.(int64),
			},
		}
	case int32:
		return &commonpb.AnyValue{
			Value: &commonpb.AnyValue_IntValue{
				IntValue: int64(value.(int32)),
			},
		}
	case int:
		return &commonpb.AnyValue{
			Value: &commonpb.AnyValue_IntValue{
				IntValue: int64(value.(int)),
			},
		}
	case bool:
		return &commonpb.AnyValue{
			Value: &commonpb.AnyValue_BoolValue{
				BoolValue: value.(bool),
			},
		}
	}
	return nil
}

func getJsonValue(pv *commonpb.AnyValue) interface{} {
	if pv == nil {
		return nil
	}
	switch pv.Value.(type) {
	case *commonpb.AnyValue_StringValue:
		return pv.GetStringValue()
	case *commonpb.AnyValue_IntValue:
		return pv.GetIntValue()
	case *commonpb.AnyValue_BoolValue:
		return pv.GetBoolValue()
	}
	return nil
}

func getTimeUnixNano(value interface{}) uint64 {
	switch value.(type) {
	case string:
		t, err := time.Parse(time.RFC3339Nano, value.(string))
		if err != nil {
			fmt.Printf("Unknown time string %s\n", err)
			return 0
		}
		return uint64(t.UnixNano())
	case int64:
		return uint64(value.(int64)) * 1000000
	case int:
		return uint64(value.(int)) * 1000000
	case float32:
		return uint64(value.(float32)) * 1000000
	case float64:
		return uint64(value.(float64)) * 1000000
	default:
		fmt.Printf("Unknown time type %T\n", value)
	}
	return uint64(0)
}
