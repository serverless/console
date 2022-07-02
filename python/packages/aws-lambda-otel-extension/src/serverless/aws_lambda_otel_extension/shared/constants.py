import logging
import typing

from pkg_resources import get_distribution

_package_distribution = get_distribution("serverless-aws-lambda-otel-extension")

PACKAGE_NAME = _package_distribution.project_name
PACKAGE_NAMESPACE = "serverless.aws_lambda_otel_extension"
PACKAGE_VERSION = _package_distribution.version

HANDLER_INSTRUMENTATION_NAME = "serverless.aws_lambda_otel_extension"
WRAPPER_INSTRUMENTATION_NAME = "serverless.aws_lambda_otel_extension.aws_lambda.instrumentor"

_X_AMZN_TRACE_ID_ENV_VAR = "_X_AMZN_TRACE_ID"

AWS_DEFAULT_REGION_ENV_VAR = "AWS_DEFAULT_REGION"
AWS_LAMBDA_FUNCTION_MEMORY_SIZE_ENV_VAR = "AWS_LAMBDA_FUNCTION_MEMORY_SIZE"
AWS_LAMBDA_FUNCTION_NAME_ENV_VAR = "AWS_LAMBDA_FUNCTION_NAME"
AWS_LAMBDA_FUNCTION_VERSION_ENV_VAR = "AWS_LAMBDA_FUNCTION_VERSION"
AWS_LAMBDA_LOG_GROUP_NAME_ENV_VAR = "AWS_LAMBDA_LOG_GROUP_NAME"
AWS_LAMBDA_LOG_STREAM_NAME_ENV_VAR = "AWS_LAMBDA_LOG_STREAM_NAME"
AWS_LAMBDA_RUNTIME_API_ENV_VAR = "AWS_LAMBDA_RUNTIME_API"
AWS_REGION_ENV_VAR = "AWS_REGION"

OTEL_PYTHON_DISABLED_INSTRUMENTATIONS_ENV_VAR = "OTEL_PYTHON_DISABLED_INSTRUMENTATIONS"
OTEL_PYTHON_ENABLED_INSTRUMENTATIONS_ENV_VAR = "OTEL_PYTHON_ENABLED_INSTRUMENTATIONS"
OTEL_PYTHON_LOG_CORRELATION_ENV_VAR = "OTEL_PYTHON_LOG_CORRELATION"

SLS_OTEL_EXTENSION_FLUSH_TIMEOUT_ENV_VAR = "SLS_OTEL_EXTENSION_FLUSH_TIMEOUT"
SLS_OTEL_EXTENSION_LOG_LEVEL_ENV_VAR = "SLS_OTEL_EXTENSION_LOG_LEVEL"
SLS_OTEL_USER_SETTINGS_ENV_VAR = "SLS_OTEL_USER_SETTINGS"

TEST_DRY_LOG_ENV_VAR = "TEST_DRY_LOG"
TEST_DRY_LOG_PRETTY_ENV_VAR = "TEST_DRY_LOG_PRETTY"

HTTP_CONTENT_TYPE_APPLICATION_JSON = "application/json"
HTTP_CONTENT_TYPE_HEADER = "Content-Type"

HTTP_METHOD_GET = "GET"
HTTP_METHOD_POST = "POST"
HTTP_METHOD_PUT = "PUT"

# Falsy would be if this doesn't match in cases where don't want to specifically test for a valid value.
TRUTHY = ["true", "1", "1.0", "yes", "on", "yup", "mmhmm"]

LOG_LEVEL_MAP: typing.Dict[str, int] = {
    "critical": logging.CRITICAL,
    "debug": logging.DEBUG,
    "error": logging.ERROR,
    "info": logging.INFO,
    "warning": logging.WARNING,
}

INSTRUMENTATION_TILDE_MAP = {
    "~common": [
        "botocore",
        "logging",
        "requests",
        "urllib",
    ],
}

_HANDLER_ENV_VAR = "_HANDLER"
ORIG_HANDLER_ENV_VAR = "ORIG_HANDLER"
