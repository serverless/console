# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [0.5.2](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.5.1...@serverless/aws-lambda-otel-extension@0.5.2) (2022-07-14)

### Bug Fixes

- Ensure broken settings do not break extension ([9d09dfd](https://github.com/serverless/runtime/commit/9d09dfd31a7b361e8f348b1b3bdafb99b1162c6d))

### [0.5.1](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.5.0...@serverless/aws-lambda-otel-extension@0.5.1) (2022-07-14)

### Performance Improvements

- Patch ineffective `require-in-the-middle` setup ([62672b9](https://github.com/serverless/runtime/commit/62672b95dc98a61488facf63a68347da1b17d626))

## [0.5.0](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.4.6...@serverless/aws-lambda-otel-extension@0.5.0) (2022-07-12)

### ⚠ BREAKING CHANGES

- Configuration variable `SLS_OTEL_USER_SETTINGS` is renamed to `SLS_EXTENSION`
- Introduce new mandatory `orgId`, `namespace` and `environment` settings (required by the Serverless Console)
- Extension is no longer observability tool agnostic. It's reconfigured to report specifically to the Serverless Console.
  At the same time its configuration settings are simplified:
  - Support for `destination` in all configuration groups is dropped.
  - New required `ingestToken` property is introduced (for injection of kinesis server authentication token)
  - Support for `metrics.outputType` and `traces.outputType` was removed (reports are hardcoded to be send in protobuf.

### Features

- New mandatory `orgId`, `namespace` and `environment` settings ([cc009de](https://github.com/serverless/runtime/commit/cc009de84cfe54f1934644c4bc79beb890869262))
- On invalid user settings make extension ineffective ([8f7128b](https://github.com/serverless/runtime/commit/8f7128b2c747f5665273fd71af2a753f4a75bc89))
- Remove support for `[metrics|traces].outputType` ([91e7e03](https://github.com/serverless/runtime/commit/91e7e03bda720dbf124f3c26465e0621942f6c24))
- Rename `SLS_OTEL_USER_SETTINGS` into `SLS_EXTENSION` ([aa2b00b](https://github.com/serverless/runtime/commit/aa2b00beebee2e93e56f81e73cde13467db4c1ba))
- Unconditionally send reports to the Serverless Console ([80ef8dc](https://github.com/serverless/runtime/commit/80ef8dc2982c1ad3ce64cabbe7d246ee2fb9ce60))
- Support bundling of custom settings into layers ([cd699bc](https://github.com/serverless/runtime/commit/cd699bceda46df76d1d7d5b95923edcbb206938f))
- Allow environment changes in bundled settings ([6e7bf67](https://github.com/serverless/runtime/commit/6e7bf67bbba786fa5dfdb1ca3cb4313f8e5d2323))

### Maintenance Improvements

- Add debug log on event/next response ([33950f2](https://github.com/serverless/runtime/commit/33950f21f1dcf55d29f55150a0f1a2b436e8c76d))
- Centralize user settings handling ([b2bb89d](https://github.com/serverless/runtime/commit/b2bb89dd919eee73e9d72c795828db16d16c71f5))
- Improve clarity of build script ([3db6c0d](https://github.com/serverless/runtime/commit/3db6c0dc40404ff101732428095bd4a623fdb307))

### [0.4.6](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.4.5...@serverless/aws-lambda-otel-extension@0.4.6) (2022-06-29)

### Bug Fixes

- Fix handling of promise resolution ([867b3bc](https://github.com/serverless/runtime/commit/867b3bc59e73bcda99f8e30a008fadd98b9c0f52))

### Maintenance Improvements

- Improve resolution handling ([938ce98](https://github.com/serverless/runtime/commit/938ce98fe03900739f9a1c9b1765e3f242d5c657))
- Simplify stale callback handling ([fb2885a](https://github.com/serverless/runtime/commit/fb2885a9004c06e31db3eec2d2bab83f2fa39617))

### [0.4.5](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.4.4...@serverless/aws-lambda-otel-extension@0.4.5) (2022-06-29)

### Bug Fixes

- Fix handling of a case where `telemetryData` arrives before `eventData` ([26991ee](https://github.com/serverless/runtime/commit/26991ee530110901ec1de1a945351a7d42ec4166))

### Maintenance Improvements

- Add debug logs on internal extension requests timing ([2d2c2db](https://github.com/serverless/runtime/commit/2d2c2dbfe2b5fe99cd13a262031ae1d0b81be688))
- Ensure to reuse HTTP(S) connection sockets ([82d295a](https://github.com/serverless/runtime/commit/82d295ac6f1ab7255ba6c67a890c4153b220d25c))
- Remove `node-fetch` dependency ([688e24d](https://github.com/serverless/runtime/commit/688e24d08e94219df23038fc90b8325616d50de7))

### [0.4.4](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.4.3...@serverless/aws-lambda-otel-extension@0.4.4) (2022-06-27)

### Bug Fixes

- Fix handling of doubled resolution ([8eacf17](https://github.com/serverless/runtime/commit/8eacf17644a82d1e164a2752aa35a5289becce8b))

### Performance Improvements

- Optimize protobuf encoding by using pre-compiled interface ([70f26e6](https://github.com/serverless/runtime/commit/70f26e6faf1c5315f11e6b3f7dee1b02956b8186))
- Improve telemetry server response time ([60ae043](https://github.com/serverless/runtime/commit/60ae043d31e5f13594a2743b2365343d47898a72))

### Maintenance Improvements

- Clarify wrapping logic ([2b203f0](https://github.com/serverless/runtime/commit/2b203f0445d1e322b518b91c555446f3fc017e6d))
- Cleanup invalid _async_ mark ([7423971](https://github.com/serverless/runtime/commit/74239715d008f26639856294dc611f3cd429eb07))

### [0.4.3](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.4.2...@serverless/aws-lambda-otel-extension@0.4.3) (2022-06-23)

### Features

- Add attributes and resource to req/res payload ([18d69e9](https://github.com/serverless/runtime/commit/18d69e9a785a7019dceb2f3d657b7b31e2cf7164))

### Maintenance Improvements

- Changed getCurrentRequestData to getRequestContext ([a80a961](https://github.com/serverless/runtime/commit/a80a96163119615065c377c9bb3b18c4e2b8629f))
- Added invoke time to internal context ([c70d60a](https://github.com/serverless/runtime/commit/c70d60aefb80099cb1422b1f2d745860f7258135))

### [0.4.2](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.4.1...@serverless/aws-lambda-otel-extension@0.4.2) (2022-06-15)

### Maintenance Improvements

- Adjusted errorType for handled/unhandled ([4fe0437](https://github.com/serverless/runtime/commit/4fe0437a2d4407dfba8f0aff191fca07c89afc48))
- Updated timeout errorType ([d9b60cd](https://github.com/serverless/runtime/commit/d9b60cd6e2df9613ae4c765d9587d7f0cb6ce82d))

### [0.4.1](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.4.0...@serverless/aws-lambda-otel-extension@0.4.1) (2022-06-13)

### Bug Fixes

- Fix ESM support in Node.js v16 runtime ([2547826](https://github.com/serverless/runtime/commit/25478264044f453ea64acd71940f271077095537))

### Maintenance Improvements

- Cleanup timeout setup ([ec7c67a](https://github.com/serverless/runtime/commit/ec7c67a88d219f43cabcd9c16ac004b5eb7fc64b))
- Debug log telemetry data cache handling ([b69e62f](https://github.com/serverless/runtime/commit/b69e62f5c47a707082f3d8f13523a089c4e37e70))
- Ensure to flush report in one log line ([64ca1dd](https://github.com/serverless/runtime/commit/64ca1dd40886dfa2c4c5b5fe4ec965ae8811cb14))
- Generalize debugLogs handling ([bd6340b](https://github.com/serverless/runtime/commit/bd6340b629c87787ca7bb4e8f4dbd8b94275740d))
- Improve debug logging ([6aef91c](https://github.com/serverless/runtime/commit/6aef91c2eb2ad0f98eeacff50922efbd77cbd494))
- Unify debug log handling ([57aceba](https://github.com/serverless/runtime/commit/57aceba04e05df1f664237401d7e7811ef5f26be))

## [0.4.0](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.3.2...@serverless/aws-lambda-otel-extension@0.4.0) (2022-06-02)

### ⚠ BREAKING CHANGES

- All configuration settings are now set in JSON object in context of `SLS_OTEL_USER_SETTINGS` environment variable. Handling of following environment variables was replaced:
  - `SLS_OTEL_REPORT_METRICS_URL` (replaced by `metrics.destination`)
  - `SLS_OTEL_REPORT_TRACES_URL` (replaced by `traces.destination`)
  - `SLS_OTEL_REPORT_LOGS_URL` (replaced by `logs.destination`)
  - `SLS_OTEL_REPORT_REQUEST_RESPONSE_URL` (replaced by `request.destination` and `response.destination`)
  - `SLS_OTEL_REPORT_REQUEST_HEADERS` (replaced by `common.destination.requestHeaders`)
- `disableLogsMonitoring` is replaced by `logs.disabled`
- `disableRequestResponseMonitoring` is replaced by `request.disabled` and `response.disabled`

### Features

- Generalize request handling and integrate settings propagated via environment variables into user settings ([e73acca](https://github.com/serverless/runtime/commit/e73acca9db5fc1ceb0681f7ddead978b6b6651b9))
- Support new user settings format ([c6904b6](https://github.com/serverless/runtime/commit/c6904b6b6ae94a60fe8d496a1f11ac654cb8173d))

### Bug Fixes

- Wait for request/response hooks before settling invocation ([6742606](https://github.com/serverless/runtime/commit/67426068b78cce5db7a59888aee7b553baff82a7))
- Ensure to handle property thenable/callback race ([77bf67f](https://github.com/serverless/runtime/commit/77bf67fdaa01aa58e620e04e80c68a943480d599))

### Maintenance Improvements

- Do not accept multiple paylaods at request functions ([db33381](https://github.com/serverless/runtime/commit/db3338183925fb8b55dc10a7f57553636e07dbf3))

### [0.3.2](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.3.1...@serverless/aws-lambda-otel-extension@0.3.2) (2022-05-31)

### Features

- Filter out Serverless Dashboard data logs ([b35c41c](https://github.com/serverless/runtime/commit/b35c41c61dd162606722944bce46f7a71d3475ce))

### [0.3.1](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.3.0...@serverless/aws-lambda-otel-extension@0.3.1) (2022-05-26)

### Bug Fixes

- Fix dependencies configuration ([9b41fa6](https://github.com/serverless/runtime/commit/9b41fa6198738669a5aebd58b7477f724888793d))

### Maintenance Improvements

- Bundle internal and external extensions into a single file ([9195c1e](https://github.com/serverless/runtime/commit/9195c1e8f760e27323ea65155fc3cef5e55d9ffb))
- Distinguish JSON from request data ([aa0141b](https://github.com/serverless/runtime/commit/aa0141be19c7b18cbb0f3abcab027d958f2ab052))
- Improve debug log for ingestion server requests ([d2b5233](https://github.com/serverless/runtime/commit/d2b5233020357c2384c894df5ef87e4ad2d03a0e))
- Improve debug log for telemetry payload ([b133ea6](https://github.com/serverless/runtime/commit/b133ea665b2d0308cbd41f7605e5dec15d3acba5))
- Log extension overhead processing durations ([7ddb265](https://github.com/serverless/runtime/commit/7ddb26576148f6c5466db00370fc5de770a6507a))
- Remove obsolete debug log ([59ae9a8](https://github.com/serverless/runtime/commit/59ae9a8b1b48356c5974f1dd758ceb404c098bd5))
- Rename `logMessage` to `debugLog` ([b7c2edc](https://github.com/serverless/runtime/commit/b7c2edcd3332d9383369e5ab1d8a68262b3be928))
- Simplify protobuf processing ([e339066](https://github.com/serverless/runtime/commit/e339066d4b1f7046e06fdbd911b7f838b8945596))

## [0.3.0](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.2.21...@serverless/aws-lambda-otel-extension@0.3.0) (2022-05-23)

### ⚠ BREAKING CHANGES

- Path to exec wrapper has changed from `/opt/otel-extension/internal/exec-wrapper.sh` to `/opt/otel-extension-internal-node/exec-wrapper.sh`. Change should be reflected when configuring `AWS_LAMBDA_EXEC_WRAPPER` environment variable for lambda

### Features

- Build script for runtime agnostic version of external extension ([e9fb48c](https://github.com/serverless/runtime/commit/e9fb48c0a75a058f6edac85033f8a8406fb1054a))
- Split internal and external extension into independent extensions ([13a7a72](https://github.com/serverless/runtime/commit/13a7a72a01c43d79176ef0f34b54fca64aba9651))

### Bug Fixes

- Fix handling of JSON type payloads ([5ba0175](https://github.com/serverless/runtime/commit/5ba0175b04a01311e9283c8d0b5a470a5b0d6347))

### Maintenance Improvements

- Rewrite external extension ([73cfc5c](https://github.com/serverless/runtime/commit/73cfc5c9c979f2d7b089f5fb3deca39ab1f52086))
- Drop `node-fetch` dependency in external extension ([cca15ed](https://github.com/serverless/runtime/commit/cca15ed850123ee2d7130bb7eb03d0de97ce8a52))
- Provide separate `userSettings` source for both extensions ([5a6909d](https://github.com/serverless/runtime/commit/5a6909dceabb7f73f15fb5055476e5b9d6f141dd))
- Remove `lodash.isobject` dependency ([4773a7c](https://github.com/serverless/runtime/commit/4773a7cf52b0684fc8db100b0b1321444539c612))
- Remove `OTEL_SERVER_PORT` definition from common scope ([fe443a7](https://github.com/serverless/runtime/commit/fe443a7fd5060b6182de5797aece98ba4d45aaac))
- Remove common `logMessage` ([a620a06](https://github.com/serverless/runtime/commit/a620a065e2c20d7a996ce0c64093686fea4c8c05))

### [0.2.21](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.2.20...@serverless/aws-lambda-otel-extension@0.2.21) (2022-04-29)

### Features

- Avoid sending binary response data ([dbb3711](https://github.com/serverless/runtime/commit/dbb3711b696a77290f93a36f1f2f2e190067c36a))

### [0.2.20](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.2.19...@serverless/aws-lambda-otel-extension@0.2.20) (2022-04-26)

### Bug Fixes

- Updated timeout error message ([2f49bcd](https://github.com/serverless/runtime/commit/2f49bcd291d4574384515a028060897fa36043c0))

### [0.2.19](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.2.18...@serverless/aws-lambda-otel-extension@0.2.19) (2022-04-21)

### Features

- Suport `userSettings.disableLogsMonitoring` ([54d1fa9](https://github.com/serverless/runtime/commit/54d1fa996e82566fd505d9e1975d9c641651d147))
- Support `userSettings.disableRequestResponseMonitoring` ([b7998b4](https://github.com/serverless/runtime/commit/b7998b42843e02e1fb76faf60a73e127c86bdb7d))

### Maintenance Improvements

- Centralize log listener server handling ([39fb512](https://github.com/serverless/runtime/commit/39fb51297105f8c3c491b41f151e592169601fb1))
- Centralize log listener server host and port configuration ([1c38e2d](https://github.com/serverless/runtime/commit/1c38e2d6a1a8b5dd3bc5ddef1f5bc24f8cc0e230))
- Centralize subscription configuration ([c6ffbc8](https://github.com/serverless/runtime/commit/c6ffbc8f61b2ebebceb35f4a576900ee27b3274a))
- Improve handling of dead paths ([9336587](https://github.com/serverless/runtime/commit/933658710e2dc814f0b2f2c1221c8044e924b132))
- Remove obsolete error handling ([7310032](https://github.com/serverless/runtime/commit/7310032176d1e9ddd17bde9378ea260844f6b044))
- Remove obsolete handling from telemetry server ([247dec6](https://github.com/serverless/runtime/commit/247dec6665597b7f6d56fa8454c18a32352d5458))
- Reorder initial setup ([7f9e364](https://github.com/serverless/runtime/commit/7f9e364d23ab7cf0341aa2e9e7e255f398cb2fa7))

### [0.2.18](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.2.12...@serverless/aws-lambda-otel-extension@0.2.18) (2022-04-13)

### Bug Fixes

- Adjusted code for integration tests ([5709511](https://github.com/serverless/runtime/commit/5709511ed5afd4b06c541383d39e0880feacd2fa))
- fixed integration test ([81aa93f](https://github.com/serverless/runtime/commit/81aa93fb295b980cad87b43e9b166005f138a020))
- Fixed waitFor bug ([8cf697a](https://github.com/serverless/runtime/commit/8cf697ad29d75ed04b8a3272c5e4c409c5fe37e4))

### [0.2.17](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.2.12...@serverless/aws-lambda-otel-extension@0.2.17) (2022-04-13)

### Bug Fixes

- fixed integration test ([81aa93f](https://github.com/serverless/runtime/commit/81aa93fb295b980cad87b43e9b166005f138a020))
- Fixed waitFor bug ([8cf697a](https://github.com/serverless/runtime/commit/8cf697ad29d75ed04b8a3272c5e4c409c5fe37e4))

### [0.2.13](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.2.12...@serverless/aws-lambda-otel-extension@0.2.13) (2022-04-13)

### Bug Fixes

- Fixed waitFor bug ([aa0658c](https://github.com/serverless/runtime/commit/aa0658cc8c3ab87708ec33c00f21dc93c1dc76e7))

### [0.2.12](https://github.com/serverless/runtime/compare/@serverless/aws-lambda-otel-extension@0.2.10...@serverless/aws-lambda-otel-extension@0.2.12) (2022-04-13)

### Bug Fixes

- Include spanId in req/res data ([1801503](https://github.com/serverless/runtime/commit/1801503557b09d97edbda16f94f990b6914c5bad))

### Maintenance Improvements

- Automate extension version resolution ([d1d8c12](https://github.com/serverless/runtime/commit/d1d8c124563b0481383c21b69e7f13576dafbab9))
- Improve `faas.collector_version` format ([9ccbfa1](https://github.com/serverless/runtime/commit/9ccbfa10c301d1234d595ae5b242bd43a42b1ade))
- Optimize instrumentations setup ([99b47fd](https://github.com/serverless/runtime/commit/99b47fd339f1a64312ef7c37baca54f6e9967ac1))
- Prevent misleading "No modules.." warning ([85fbc39](https://github.com/serverless/runtime/commit/85fbc399ad61914104d7f86e60699042d7fe6f79))
- Remove obsolete config configuration ([5b448ba](https://github.com/serverless/runtime/commit/5b448ba4c0bac1d9f6653151cad9d075916def3f))
- Remove obsolete instrumentation registration ([014b6b3](https://github.com/serverless/runtime/commit/014b6b39647df54643aa18086183479810aaf79a))
- Remove obsolete payload compression ([f5ceb3d](https://github.com/serverless/runtime/commit/f5ceb3d9151ad2183d3ebe5e1b079f4a7e788fea))
