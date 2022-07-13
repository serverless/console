from os import getenv

from serverless.aws_lambda_otel_extension.shared.constants import (
    _HANDLER_ENV_VAR,
    AWS_DEFAULT_REGION_ENV_VAR,
    AWS_LAMBDA_FUNCTION_MEMORY_SIZE_ENV_VAR,
    AWS_LAMBDA_FUNCTION_NAME_ENV_VAR,
    AWS_LAMBDA_FUNCTION_VERSION_ENV_VAR,
    AWS_LAMBDA_LOG_GROUP_NAME_ENV_VAR,
    AWS_LAMBDA_LOG_STREAM_NAME_ENV_VAR,
    AWS_REGION_ENV_VAR,
    ORIG_HANDLER_ENV_VAR,
    OTEL_PYTHON_LOG_CORRELATION_ENV_VAR,
    SLS_EXTENSION_ENV_VAR,
    SLS_EXTENSION_INTERNAL_DISABLED_INSTRUMENTATIONS_ENV_VAR,
    SLS_EXTENSION_INTERNAL_ENABLED_INSTRUMENTATIONS_ENV_VAR,
    SLS_TEST_EXTENSION_INTERNAL_DEBUG_ENV_VAR,
    SLS_TEST_EXTENSION_INTERNAL_EXPORT_FLUSH_TIMEOUT_ENV_VAR,
    SLS_TEST_EXTENSION_INTERNAL_EXPORT_URL_ENV_VAR,
    SLS_TEST_EXTENSION_INTERNAL_LOG_ENV_VAR,
    SLS_TEST_EXTENSION_INTERNAL_LOG_PRETTY_ENV_VAR,
)

ENV__HANDLER = getenv(_HANDLER_ENV_VAR)
ENV_ORIG_HANDLER = getenv(ORIG_HANDLER_ENV_VAR)

ENV_AWS_LAMBDA_FUNCTION_MEMORY_SIZE = getenv(AWS_LAMBDA_FUNCTION_MEMORY_SIZE_ENV_VAR)
ENV_AWS_LAMBDA_FUNCTION_NAME = getenv(AWS_LAMBDA_FUNCTION_NAME_ENV_VAR)
ENV_AWS_LAMBDA_FUNCTION_VERSION = getenv(AWS_LAMBDA_FUNCTION_VERSION_ENV_VAR)
ENV_AWS_LAMBDA_LOG_GROUP_NAME = getenv(AWS_LAMBDA_LOG_GROUP_NAME_ENV_VAR)
ENV_AWS_LAMBDA_LOG_STREAM_NAME = getenv(AWS_LAMBDA_LOG_STREAM_NAME_ENV_VAR)

ENV_AWS_DEFAULT_REGION = getenv(AWS_DEFAULT_REGION_ENV_VAR)
ENV_AWS_REGION = getenv(AWS_REGION_ENV_VAR)

ENV_OTEL_PYTHON_LOG_CORRELATION = getenv(OTEL_PYTHON_LOG_CORRELATION_ENV_VAR)

ENV_SLS_EXTENSION = getenv(SLS_EXTENSION_ENV_VAR)
ENV_SLS_EXTENSION_INTERNAL_DISABLED_INSTRUMENTATIONS = getenv(SLS_EXTENSION_INTERNAL_DISABLED_INSTRUMENTATIONS_ENV_VAR)
ENV_SLS_EXTENSION_INTERNAL_ENABLED_INSTRUMENTATIONS = getenv(SLS_EXTENSION_INTERNAL_ENABLED_INSTRUMENTATIONS_ENV_VAR)

ENV_SLS_TEST_EXTENSION_INTERNAL_DEBUG = getenv(SLS_TEST_EXTENSION_INTERNAL_DEBUG_ENV_VAR)
ENV_SLS_TEST_EXTENSION_INTERNAL_EXPORT_FLUSH_TIMEOUT = getenv(SLS_TEST_EXTENSION_INTERNAL_EXPORT_FLUSH_TIMEOUT_ENV_VAR)
ENV_SLS_TEST_EXTENSION_INTERNAL_EXPORT_URL = getenv(SLS_TEST_EXTENSION_INTERNAL_EXPORT_URL_ENV_VAR)
ENV_SLS_TEST_EXTENSION_INTERNAL_LOG = getenv(SLS_TEST_EXTENSION_INTERNAL_LOG_ENV_VAR)
ENV_SLS_TEST_EXTENSION_INTERNAL_LOG_PRETTY = getenv(SLS_TEST_EXTENSION_INTERNAL_LOG_PRETTY_ENV_VAR)