'use strict';

const { expect } = require('chai');

const path = require('path');
const log = require('log').get('test');
const normalizeOtelAttributes = require('../utils/normalize-otel-attributes');
const cleanup = require('../lib/cleanup');
const createCoreResources = require('../lib/create-core-resources');
const processFunction = require('../lib/process-function');
const resolveTestVariantsConfig = require('../lib/resolve-test-variants-config');

for (const name of [
  'TEST_LAYER_FILENAME',
  'TEST_EXTERNAL_LAYER_FILENAME',
  'TEST_INTERNAL_LAYER_FILENAME',
]) {
  // In tests, current working directory is mocked,
  // so if relative path is provided in env var it won't be resolved properly
  // with this patch we resolve it before cwd mocking
  if (process.env[name]) process.env[name] = path.resolve(process.env[name]);
}

describe('integration', function () {
  this.timeout(120000);
  const coreConfig = {};

  const useCasesConfig = new Map([
    [
      'callback',
      {
        variants: new Map([
          ['v12', { configuration: { Runtime: 'nodejs12.x' } }],
          ['v14', { configuration: { Runtime: 'nodejs14.x' } }],
          ['v16', { configuration: { Runtime: 'nodejs16.x' } }],
          [
            'invalidSettings',
            {
              isInstrumentationDisabled: true,
              configuration: {
                Environment: {
                  Variables: {
                    AWS_LAMBDA_EXEC_WRAPPER: '/opt/otel-extension-internal-node/exec-wrapper.sh',
                    SLS_DEBUG_EXTENSION: '1',
                    SLS_EXTENSION: '{}',
                  },
                },
              },
            },
          ],
        ]),
      },
    ],
    ['esbuild-from-esm-callback', {}],
    [
      'esm-callback/index',
      {
        variants: new Map([
          ['v14', { configuration: { Runtime: 'nodejs14.x' } }],
          ['v16', { configuration: { Runtime: 'nodejs16.x' } }],
        ]),
      },
    ],
    [
      'express',
      {
        invokePayload: {
          version: '2.0',
          routeKey: '$default',
          rawPath: '/foo',
          rawQueryString: '',
          headers: {
            'accept':
              'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
            'accept-encoding': 'gzip, deflate, br',
            'accept-language': 'en-US,pl;q=0.7,en;q=0.3',
            'content-length': '0',
            'host': '1hqnqp4a70.execute-api.us-east-1.amazonaws.com',
            'sec-fetch-dest': 'document',
            'sec-fetch-mode': 'navigate',
            'sec-fetch-site': 'none',
            'sec-fetch-user': '?1',
            'sec-gpc': '1',
            'upgrade-insecure-requests': '1',
            'user-agent':
              'Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:98.0) Gecko/20100101 Firefox/98.0',
            'x-amzn-trace-id': 'Root=1-624605c4-7fcc8fe9188a3cb762dcd189',
            'x-forwarded-for': '80.55.87.22',
            'x-forwarded-port': '443',
            'x-forwarded-proto': 'https',
          },
          requestContext: {
            accountId: '992311060759',
            apiId: '1hqnqp4a70',
            domainName: '1hqnqp4a70.execute-api.us-east-1.amazonaws.com',
            domainPrefix: '1hqnqp4a70',
            http: {
              method: 'GET',
              path: '/foo',
              protocol: 'HTTP/1.1',
              sourceIp: '80.55.87.22',
              userAgent:
                'Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:98.0) Gecko/20100101 Firefox/98.0',
            },
            requestId: 'P3XWwjfgIAMEVFw=',
            routeKey: '$default',
            stage: '$default',
            time: '31/Mar/2022:19:49:24 +0000',
            timeEpoch: 1648756164620,
          },
          isBase64Encoded: false,
        },
        test: ({ instrumentationSpans, invocationsData }) => {
          expect(
            instrumentationSpans['@opentelemetry/instrumentation-express'].length
          ).to.be.at.least(4);
          expect(
            invocationsData.map(({ responsePayload }) => responsePayload.bodyJson)
          ).to.deep.equal([{ message: 'Hello from /foo!' }, { message: 'Hello from /foo!' }]);
        },
      },
    ],
    [
      'error-timeout',
      {
        // On timeout re-initialization of external extension gets slow, and we observe that second
        // invocation times out before actually lambda is initialized.
        // This is either because currently our external extension is Node.js based,
        // so has slow startup time, or it can be performance issue on AWS side.
        // To ensure reliable result increase timeout, so we get second invocation correct
        configuration: { Timeout: 7 },
        expectedOutcome: 'error:unhandled',
        test: ({ instrumentationSpans }) => {
          const { attributes } =
            instrumentationSpans['@opentelemetry/instrumentation-aws-lambda'][0];
          expect(attributes['faas.error_type']).to.equal('timeout');
        },
      },
    ],
    ['callback-error', { expectedOutcome: 'error:handled' }],
  ]);

  const testVariantsConfig = resolveTestVariantsConfig(useCasesConfig);

  before(async () => {
    await createCoreResources(coreConfig, {
      layerTypes:
        process.env.TEST_EXTERNAL_LAYER_FILENAME || process.env.TEST_INTERNAL_LAYER_FILENAME
          ? ['external', 'nodeInternal']
          : ['nodeAll'],
    });
    for (const testConfig of testVariantsConfig) {
      testConfig.deferredResult = processFunction(testConfig, coreConfig).catch((error) => ({
        // As we process result promises sequentially step by step in next turn, allowing them to
        // reject will generate unhandled rejection.
        // Therefore this scenario is converted to successuful { error } resolution
        error,
      }));
    }
  });

  for (const testConfig of testVariantsConfig) {
    it(testConfig.name, async () => {
      const testResult = await testConfig.deferredResult;
      if (testResult.error) throw testResult.error;
      const { expectedOutcome, isInstrumentationDisabled } = testConfig;
      const { invocationsData } = testResult;
      if (expectedOutcome !== 'error:unhandled' && !isInstrumentationDisabled) {
        // Current timeout handling is unreliable, therefore do not attempt to confirm
        // on all reports

        // While reports should come in order as specified below,
        // there were observed cases when it wasn't the case,
        // e.g. telemetryData (response) was received before eventData (request)
        const reportTypes = { request: 0, response: 0, metrics: 0, traces: 0 };
        invocationsData.map(({ reports }) => reports.map(([type]) => ++reportTypes[type]));
        expect(reportTypes.request).to.equal(2);
        expect(reportTypes.response).to.equal(2);
        expect(reportTypes.metrics).to.be.at.least(3);
        expect(reportTypes.traces).to.equal(2);
      }

      const allInvocationReports = invocationsData.map(({ reports }) => reports).flat();
      if (isInstrumentationDisabled) {
        expect(allInvocationReports.length).to.equal(0);
        if (testConfig.test) testConfig.test({ invocationsData });
      } else {
        const metricsReport = allInvocationReports.find(
          ([reportType]) => reportType === 'metrics'
        )[1];

        const tracesReport = allInvocationReports.find(
          ([reportType]) => reportType === 'traces'
        )[1];
        const resourceMetrics = normalizeOtelAttributes(
          metricsReport.resourceMetrics[0].resource.attributes
        );
        expect(resourceMetrics['faas.name']).to.equal(testConfig.configuration.FunctionName);
        const resourceSpans = normalizeOtelAttributes(
          tracesReport.resourceSpans[0].resource.attributes
        );
        expect(resourceSpans['faas.name']).to.equal(testConfig.configuration.FunctionName);

        const instrumentationSpans = {};
        for (const {
          instrumentationLibrary: { name: instrumentationLibraryName },
          spans,
        } of tracesReport.resourceSpans[0].instrumentationLibrarySpans) {
          instrumentationSpans[instrumentationLibraryName] = spans.map((span) => ({
            ...span,
            attributes: normalizeOtelAttributes(span.attributes),
          }));
        }
        log.debug('instrumentationSpans %o', instrumentationSpans);
        if (testConfig.test) {
          testConfig.test({
            invocationsData,
            metricsReport,
            tracesReport,
            resourceMetrics,
            resourceSpans,
            instrumentationSpans,
          });
        }
      }
    });
  }

  after(async () => cleanup({ skipFunctionsCleanup: true }));
});
