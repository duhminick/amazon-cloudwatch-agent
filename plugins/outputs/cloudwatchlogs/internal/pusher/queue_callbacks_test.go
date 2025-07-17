// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package pusher

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/aws/amazon-cloudwatch-agent/internal/state"
	"github.com/aws/amazon-cloudwatch-agent/logs"
	"github.com/aws/amazon-cloudwatch-agent/sdk/service/cloudwatchlogs"
	"github.com/aws/amazon-cloudwatch-agent/tool/testutil"
	"github.com/aws/amazon-cloudwatch-agent/tool/util"
)

// TestQueueCallbackRegistration tests that callbacks are registered correctly
func TestQueueCallbackRegistration(t *testing.T) {
	t.Run("RegistersStateCallbacks", func(t *testing.T) {
		var wg sync.WaitGroup
		var s stubLogsService
		var called bool

		// Mock the PutLogEvents method to verify the batch has state callbacks registered
		s.ple = func(in *cloudwatchlogs.PutLogEventsInput) (*cloudwatchlogs.PutLogEventsOutput, error) {
			called = true
			return &cloudwatchlogs.PutLogEventsOutput{}, nil
		}

		// Create a mock sender that captures the batch for inspection
		mockSender := &mockSender{}
		mockSender.On("Send", mock.AnythingOfType("*pusher.logEventBatch")).Run(func(args mock.Arguments) {
			batch := args.Get(0).(*logEventBatch)

			// Verify that both regular callbacks and state callbacks are registered
			assert.NotEmpty(t, batch.doneCallbacks, "Regular callbacks should be registered")
			assert.NotEmpty(t, batch.stateCallbacks, "State callbacks should be registered")

			// Call the sender's original implementation
			s.PutLogEvents(batch.build())
		}).Return()

		// Create a queue with our mock sender
		logger := testutil.NewNopLogger()
		stop := make(chan struct{})
		q := &queue{
			target:          Target{"G", "S", util.StandardLogGroupClass, -1},
			logger:          logger,
			converter:       newConverter(logger, Target{"G", "S", util.StandardLogGroupClass, -1}),
			batch:           newLogEventBatch(Target{"G", "S", util.StandardLogGroupClass, -1}, nil),
			sender:          mockSender,
			eventsCh:        make(chan logs.LogEvent, 100),
			flushCh:         make(chan struct{}),
			resetTimerCh:    make(chan struct{}),
			flushTimer:      time.NewTimer(10 * time.Millisecond),
			stop:            stop,
			startNonBlockCh: make(chan struct{}),
			wg:              &wg,
		}
		q.flushTimeout.Store(10 * time.Millisecond)

		// Add an event and trigger send
		q.batch.append(newLogEvent(time.Now(), "test message", nil))
		q.send()

		// Verify expectations
		mockSender.AssertExpectations(t)
		assert.True(t, called, "PutLogEvents should have been called")
	})

	t.Run("RegistersStateCallbacksForStatefulEvents", func(t *testing.T) {
		var wg sync.WaitGroup

		// Create a mock range queue
		mrq := &mockRangeQueue{}
		mrq.On("ID").Return("test-queue")
		mrq.On("Enqueue", mock.Anything).Return()

		// Create a mock sender that captures the batch for inspection
		mockSender := &mockSender{}
		mockSender.On("Send", mock.AnythingOfType("*pusher.logEventBatch")).Run(func(args mock.Arguments) {
			batch := args.Get(0).(*logEventBatch)

			// Verify that state callbacks are registered
			assert.NotEmpty(t, batch.stateCallbacks, "State callbacks should be registered")

			// Verify that the batch has a batcher for our queue
			batcher, ok := batch.batchers["test-queue"]
			assert.True(t, ok, "Batch should have a batcher for our queue")
			assert.NotNil(t, batcher, "Batcher should not be nil")
		}).Return()

		// Create a queue with our mock sender
		logger := testutil.NewNopLogger()
		stop := make(chan struct{})
		q := &queue{
			target:          Target{"G", "S", util.StandardLogGroupClass, -1},
			logger:          logger,
			converter:       newConverter(logger, Target{"G", "S", util.StandardLogGroupClass, -1}),
			batch:           newLogEventBatch(Target{"G", "S", util.StandardLogGroupClass, -1}, nil),
			sender:          mockSender,
			eventsCh:        make(chan logs.LogEvent, 100),
			flushCh:         make(chan struct{}),
			resetTimerCh:    make(chan struct{}),
			flushTimer:      time.NewTimer(10 * time.Millisecond),
			stop:            stop,
			startNonBlockCh: make(chan struct{}),
			wg:              &wg,
		}
		q.flushTimeout.Store(10 * time.Millisecond)

		// Create a stateful log event
		event := newStubStatefulLogEvent("test message", time.Now(), state.NewRange(10, 20), mrq)

		// Convert and append the event
		convertedEvent := q.converter.convert(event)
		q.batch.append(convertedEvent)

		// Send the batch
		q.send()

		// Verify expectations
		mockSender.AssertExpectations(t)
	})

	t.Run("ConverterHandlesStatefulEvents", func(t *testing.T) {
		// Create a mock range queue
		mrq := &mockRangeQueue{}
		mrq.On("ID").Return("test-queue")

		// Create a converter
		logger := testutil.NewNopLogger()
		converter := newConverter(logger, Target{"G", "S", util.StandardLogGroupClass, -1})

		// Create a stateful log event
		event := newStubStatefulLogEvent("test message", time.Now(), state.NewRange(10, 20), mrq)

		// Convert the event
		convertedEvent := converter.convert(event)

		// Verify the converted event has the state information
		require.NotNil(t, convertedEvent.state, "Converted event should have state information")
		assert.Equal(t, state.NewRange(10, 20), convertedEvent.state.r, "Range should be preserved")
		assert.Equal(t, mrq, convertedEvent.state.queue, "Queue should be preserved")

		// Verify the Done callback is not set for stateful events
		// This is important because we want to use the state callbacks instead
		assert.Nil(t, convertedEvent.doneCallback, "Done callback should be nil for stateful events")
	})
}

// mockSender is a mock implementation of the Sender interface
type mockSender struct {
	mock.Mock
}

func (m *mockSender) Send(batch *logEventBatch) {
	m.Called(batch)
}

func (m *mockSender) SetRetryDuration(d time.Duration) {
	m.Called(d)
}

func (m *mockSender) RetryDuration() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}
