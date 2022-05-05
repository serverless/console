'use strict';

module.exports = function eventType(event) {
  const type = 'aws.firehose';
  const { records = [] } = event;
  return event.deliveryStreamArn && records[0] && records[0].kinesisRecordMetadata ? type : false;
};
