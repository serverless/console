'use strict';

const { expect } = require('chai');
const { EventEmitter } = require('events');
const evilDns = require('evil-dns');
const log = require('log').get('test');
const requireUncached = require('ncjsm/require-uncached');
const overwriteStdoutWrite = require('process-utils/override-stdout-write');
const getExtensionServerMock = require('../../utils/get-extension-server-mock');
const normalizeOtelAttributes = require('../../utils/normalize-otel-attributes');
const { default: fetch } = require('node-fetch');
const ensureNpmDependencies = require('../../../scripts/lib/ensure-npm-dependencies');

const OTEL_SERVER_PORT = 2772;
const port = 9001;

describe('external', () => {
  before(async () => {
    ensureNpmDependencies('external/otel-extension-external');
    evilDns.add('sandbox', '127.0.0.1');
    process.env.AWS_LAMBDA_RUNTIME_API = `127.0.0.1:${port}`;
    process.env.SLS_OTEL_USER_SETTINGS = JSON.stringify({
      metrics: { outputType: 'json' },
      traces: { outputType: 'json' },
    });
    process.env.SLS_TEST_RUN = '1';
  });

  it('should handle plain success invocation', async () => {
    const requestId = 'bf8bcf52-ff05-4f30-85cc-8a8bb1a27ae0';
    const emitter = new EventEmitter();
    const { server, listenerEmitter } = getExtensionServerMock(emitter, { requestId });

    server.listen(port);
    let stdoutData = '';
    const extensionProcess = overwriteStdoutWrite(
      (data) => (stdoutData += data),
      async () => requireUncached(() => require('../../../external/otel-extension-external'))
    );
    await Promise.all([
      new Promise((resolve) => listenerEmitter.once('listener', resolve)),
      new Promise((resolve) => listenerEmitter.once('logsSubscription', resolve)),
    ]);
    emitter.emit('event', { eventType: 'INVOKE', requestId });

    emitter.emit('logs', [
      {
        time: '2022-02-14T15:31:24.674Z',
        type: 'platform.start',
        record: {
          requestId,
          version: '$LATEST',
        },
      },
      {
        time: '2022-02-14T15:31:24.676Z',
        type: 'platform.extension',
        record: {
          name: 'otel-extension',
          state: 'Ready',
          events: ['INVOKE', 'SHUTDOWN'],
        },
      },
    ]);
    // Emit init logs
    await fetch(`http://localhost:${OTEL_SERVER_PORT}`, {
      method: 'post',
      body: JSON.stringify({
        recordType: 'eventData',
        record: {
          eventData: {
            [requestId]: {
              functionName: 'testFunction',
              computeCustomEnvArch: 'x86',
              computeRegion: 'us-east-1',
              eventCustomApiId: requestId,
            },
          },
          span: {
            traceId: 'trace-123',
            spanId: 'span-123',
          },
        },
      }),
      headers: {
        'Content-Type': 'application/json',
      },
    });

    await fetch(`http://localhost:${OTEL_SERVER_PORT}`, {
      method: 'post',
      body: JSON.stringify({
        recordType: 'telemetryData',
        requestId,
        record: {
          function: {
            'service.name': 'test-otel-extension-success',
            'telemetry.sdk.language': 'nodejs',
            'telemetry.sdk.name': 'opentelemetry',
            'telemetry.sdk.version': '1.0.1',
            'cloud.provider': 'aws',
            'cloud.platform': 'aws_lambda',
            'cloud.region': 'us-east-1',
            'faas.name': 'test-otel-extension-success',
            'faas.version': '$LATEST',
            'process.pid': 18,
            'process.executable.name': '/var/lang/bin/node',
            'process.command': '/var/runtime/index.js',
            'process.command_line': '/var/lang/bin/node /var/runtime/index.js',
            'computeCustomArn':
              'arn:aws:lambda:us-east-1:992311060759:function:test-otel-extension-success',
            'functionName': 'test-otel-extension-success',
            'computeRegion': 'us-east-1',
            'computeRuntime': 'aws.lambda.nodejs.14.18.1',
            'computeCustomFunctionVersion': '$LATEST',
            'computeMemorySize': '128',
            'eventCustomXTraceId':
              'Root=1-620a52a7-17a756211c086b284224da3d;Parent=2b73def916cea6f4;Sampled=0',
            'computeCustomLogStreamName': '2022/02/14/[$LATEST]a98c4430fe4f4f59a0eb0d360d96213d',
            'computeCustomEnvArch': 'x64',
            'eventType': null,
            'eventCustomRequestId': 'bf8bcf52-ff05-4f30-85cc-8a8bb1a27ae0',
            'computeIsColdStart': true,
            'eventCustomDomain': null,
            'eventCustomRequestTimeEpoch': null,
            'startTime': 1644843688147,
            'endTime': 1644843690251,
            'error': false,
            'httpStatusCode': 200,
          },
          traces: {
            resourceSpans: [
              {
                resource: {
                  'service.name': 'test-otel-extension-success',
                  'telemetry.sdk.language': 'nodejs',
                  'telemetry.sdk.name': 'opentelemetry',
                  'telemetry.sdk.version': '1.0.1',
                  'cloud.provider': 'aws',
                  'cloud.platform': 'aws_lambda',
                  'cloud.region': 'us-east-1',
                  'faas.name': 'test-otel-extension-success',
                  'faas.version': '$LATEST',
                  'process.pid': 18,
                  'process.executable.name': '/var/lang/bin/node',
                  'process.command': '/var/runtime/index.js',
                  'process.command_line': '/var/lang/bin/node /var/runtime/index.js',
                },
                instrumentationLibrarySpans: [
                  {
                    instrumentationLibrary: {
                      name: '@opentelemetry/instrumentation-aws-lambda',
                      version: '0.28.1',
                    },
                    spans: [
                      {
                        traceId: '571ce462a1147e9e5c586d547259451e',
                        spanId: 'df3acdfbb928da48',
                        name: 'test-otel-extension-success',
                        kind: 'SPAN_KIND_SERVER',
                        startTimeUnixNano: '1644843688147651600',
                        endTimeUnixNano: '1644843690251790300',
                        attributes: {
                          'faas.execution': 'bf8bcf52-ff05-4f30-85cc-8a8bb1a27ae0',
                          'faas.id':
                            'arn:aws:lambda:us-east-1:992311060759:function:test-otel-extension-success',
                          'cloud.account.id': '992311060759',
                          'http.status_code': 200,
                        },
                        status: {},
                      },
                    ],
                  },
                ],
              },
            ],
          },
        },
      }),
      headers: {
        'Content-Type': 'application/json',
      },
    });
    emitter.emit('logs', [
      {
        time: '2022-02-14T15:31:26.713Z',
        type: 'function',
        record: `2022-02-14T13:01:30.307Z\t${requestId}\tINFO\tHi mom`,
      },
      {
        time: '2022-02-14T15:31:26.713Z',
        type: 'platform.runtimeDone',
        record: {
          requestId,
          status: 'success',
        },
      },
      {
        time: '2022-02-14T15:31:26.742Z',
        type: 'platform.end',
        record: {
          requestId,
        },
      },
    ]);
    await new Promise((resolve) => listenerEmitter.once('listener', resolve));
    emitter.emit('event', { eventType: 'SHUTDOWN', requestId });
    emitter.emit('logs', [
      {
        time: '2022-02-14T15:31:26.742Z',
        type: 'platform.report',
        record: {
          requestId,
          metrics: {
            durationMs: 2064.05,
            billedDurationMs: 2065,
            memorySizeMB: 128,
            maxMemoryUsedMB: 67,
            initDurationMs: 238.12,
          },
        },
      },
    ]);

    await extensionProcess;
    server.close();

    log.debug('report string %s', stdoutData);
    const reports = {};
    for (const reportString of stdoutData.split('\n').filter(Boolean)) {
      const reportType = reportString.slice(2, reportString.indexOf(':'));
      if (!reports[reportType]) reports[reportType] = [];
      const jsonString = reportString.slice(reportString.indexOf(':') + 2);
      log.debug('report %s JSON string %s', reportType, jsonString);
      reports[reportType].push(JSON.parse(jsonString));
    }
    log.debug('reports %o', reports);
    const metricsReport = reports.metrics[0];
    const tracesReport = reports.traces[0];
    const logReport = reports.logs[0][0];

    const resourceMetrics = normalizeOtelAttributes(
      metricsReport.resourceMetrics[0].resource.attributes
    );
    expect(resourceMetrics['faas.name']).to.equal('test-otel-extension-success');
    const resourceSpans = normalizeOtelAttributes(
      tracesReport.resourceSpans[0].resource.attributes
    );
    expect(resourceSpans['faas.name']).to.equal('test-otel-extension-success');

    expect(logReport.Body).to.equal(`2022-02-14T13:01:30.307Z\t${requestId}\tINFO\tHi mom`);
    expect(logReport.Attributes['faas.name']).to.equal('testFunction');
    expect(logReport.Resource['faas.arch']).to.equal('x86');
  });
});
