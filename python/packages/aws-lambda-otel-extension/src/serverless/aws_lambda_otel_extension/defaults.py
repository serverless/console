import logging
import typing

from serverless.aws_lambda_otel_extension import constants

SLS_OPENTELEMETRY_SERVER_URL = "http://sandbox:2772"
# TODO: Decide on this value with rest of team.
SLS_OPENTELEMETRY_HTTP_REQUEST_TIMEOUT = 5

SLS_OTEL_METRICS_ENABLED = False

AWS_REGION = "us-east-1"

# See: https://github.com/serverless/serverless/blob/main/lib/plugins/aws/invoke-local/runtime-wrappers/invoke.py#L17
# See:
AWS_LAMBDA_FUNCTION_NAME = "Fake"
AWS_LAMBDA_FUNCTION_VERSION = "LATEST"
AWS_LAMBDA_FUNCTION_TIMEOUT = 6
AWS_LAMBDA_FUNCTION_ARN = "arn:aws:lambda:serverless:Fake"
AWS_LAMBDA_FUNCTION_MEMORY_SIZE = 1024
AWS_LAMBDA_FUNCTION_REQUEST_ID = "1234567890"
AWS_LAMBDA_FUNCTION_LOG_GROUP_NAME = "/aws/lambda/Fake"

OTEL_PYTHON_ENABLED_INSTRUMENTATIONS: typing.List[str] = constants.INSTRUMENTATION_TILDE_MAP["~common"]
OTEL_PYTHON_DISABLED_INSTRUMENTATIONS: typing.List[str] = []

SLS_WRAPPER_LOG_LEVEL = logging.CRITICAL

OTEL_INSTRUMENTATION_AWS_LAMBDA_FLUSH_TIMEOUT = 1000
