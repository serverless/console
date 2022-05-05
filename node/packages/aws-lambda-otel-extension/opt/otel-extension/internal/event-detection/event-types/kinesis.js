'use strict';

module.exports = function eventType(event) {
  const type = 'aws.kinesis';
  const { Records = [] } = event;
  const [firstEvent = {}] = Records;
  const { eventSource } = firstEvent;
  // test is for firstEvent.eventVersion === '1.0'
  return eventSource === 'aws:kinesis' ? type : false;
};
