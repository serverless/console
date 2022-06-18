# This is an outer wrapper script meant to address issues with the `imp` module used by the Python3.6/3.7/3.8 runtime.
# The `imp` module attempts to load a module from a file path that is not a valid path since it is a dot notated
# namespace where the `__init__.py` is not present in the top level directory.
#
#   src/serverless (doesn't have __init__.py)
#   src/serverless/aws_lambda_otel_extension (has __init__.py)
#   src/serverless_aws_lambda_otel_extension_internal/wrapper.py (has __init__.py)
#
# We want to retain the use of the `serverless.*` namespace for practicality reasons and wedges like this are required
# when dealing with code that operates similar to the AWS bootstrap code via older import discovery methods.

import serverless.aws_lambda_otel_extension.internal.wrapper

# "Export" the auto_instrumenting_handler function from the internal wrapper module.
auto_instrumenting_handler = serverless.aws_lambda_otel_extension.internal.wrapper.auto_instrumenting_handler
