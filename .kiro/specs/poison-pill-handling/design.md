# Design Document: Poison Pill Handling for CloudWatch Agent

## Overview

The Amazon CloudWatch Agent currently updates the state file only when a batch of logs is successfully sent to CloudWatch Logs. This is done by calling the `Done()` method on log events, which ultimately enqueues the range to a range manager that updates the state file. However, when a batch fails to be sent after exhausting all retry attempts (a "poison pill" batch), the state file is not updated. This means that after the agent restarts, it will attempt to process the same problematic batch again, potentially causing the same failure to repeat.

This design document outlines the approach to improve the agent's resilience by updating the state file even for poison pill batches that have exhausted their retry attempts, allowing the agent to move past problematic log entries.

## Architecture

The CloudWatch Agent's log processing pipeline consists of several key components:

1. **Log Sources** (e.g., logfile, windows_event_log): These components read logs from various sources and generate `LogEvent` objects.

2. **LogAgent**: Connects log sources to log destinations.

3. **CloudWatch Logs Output Plugin**: Processes log events and sends them to CloudWatch Logs.

4. **Pusher**: Handles the actual sending of log batches to CloudWatch Logs, including retry logic.

5. **State Management**: Tracks the processing state of logs to avoid duplicate processing after restarts.

The current issue occurs in the `sender.Send()` method in the pusher component. When a batch fails to be sent after exhausting all retry attempts, the method simply returns without calling `batch.done()`, which means the state file is not updated.

## Components and Interfaces

### Key Components to Modify

1. **Sender** (`plugins/outputs/cloudwatchlogs/internal/pusher/sender.go`):
   - The `Send()` method needs to be modified to call `batch.done()` even when a batch fails after exhausting all retry attempts.

2. **LogEventBatch** (`plugins/outputs/cloudwatchlogs/internal/pusher/batch.go`):
   - The `done()` method is responsible for executing all registered callbacks, which ultimately update the state file.

### Interfaces

No changes to interfaces are required. The existing interfaces provide all the necessary functionality:

- `LogEvent` interface with the `Done()` method
- `StatefulLogEvent` interface that extends `LogEvent` with state management capabilities
- `FileRangeManager` interface for managing the state file

## Data Models

No changes to data models are required. The existing data models are sufficient:

- `Range` represents a range of log entries that have been processed
- `RangeTracker` manages collections of ranges and handles serialization to the state file
- `logEventBatch` represents a batch of log events to be sent to CloudWatch Logs

## Implementation Details

The implementation will focus on modifying the `Send()` method in the `sender` struct to ensure that `batch.done()` is called even when a batch fails after exhausting all retry attempts.

### Current Implementation

Currently, in `sender.Send()`, the `batch.done()` method is only called when the batch is successfully sent:

```go
output, err := s.service.PutLogEvents(input)
if err == nil {
    // ... handle successful response ...
    batch.done()
    s.logger.Debugf("Pusher published %v log events...", ...)
    return
}
// ... handle errors and retry logic ...
```

If all retry attempts fail, the method simply returns without calling `batch.done()`:

```go
if time.Since(startTime)+wait > s.RetryDuration() {
    s.logger.Errorf("All %v retries to %v/%v failed for PutLogEvents, request dropped.", ...)
    return
}
```

### Proposed Implementation

The proposed implementation needs to be more nuanced than simply calling `batch.done()` when all retry attempts are exhausted. The `batch.done()` method executes all registered callbacks, which include not only updating the state file but also other operations like updating the last sent time and adding stats, which we don't want to execute for failed batches.

We need to specifically call only the callbacks related to updating the state file (enqueueing the range to the range manager) without executing other success-related callbacks. This requires a more targeted approach:

1. Modify the `logEventBatch` struct to separate state-updating callbacks from other success callbacks.

2. Add a new method `updateStateOnly()` to the `logEventBatch` struct that only executes the state-updating callbacks.

3. Call this new method when all retry attempts are exhausted:

```go
if time.Since(startTime)+wait > s.RetryDuration() {
    s.logger.Errorf("All %v retries to %v/%v failed for PutLogEvents, request dropped.", ...)
    s.logger.Debugf("Updating state file for failed batch to prevent reprocessing after restart")
    batch.updateStateOnly() // Only update state file, not other success metrics
    return
}
```

This approach ensures that we only update the state file to prevent reprocessing of the failed batch after restart, without executing other success-related callbacks that would incorrectly indicate the batch was successfully sent.

## Error Handling

The existing error handling in the `sender.Send()` method is robust and will be maintained. The only change is to ensure that `batch.done()` is called even when all retry attempts are exhausted.

The method already handles various AWS errors appropriately:
- `ResourceNotFoundException`: Attempts to create the log stream
- `InvalidParameterException` and `DataAlreadyAcceptedException`: Logs an error and returns without retrying
- Other AWS errors: Logs an error and retries with appropriate backoff

## Testing Strategy

### Unit Tests

1. Modify existing unit tests for the `sender.Send()` method to verify that `batch.done()` is called when all retry attempts are exhausted.

2. Add a new test case that simulates a batch failing after exhausting all retry attempts and verifies that the state file is updated.

### Integration Tests

1. Create an integration test that simulates a poison pill batch (e.g., by injecting a batch that will always be rejected by CloudWatch Logs) and verifies that the state file is updated correctly.

2. Verify that after restarting the agent, the poison pill batch is not reprocessed.

## Conclusion

This design proposes a simple but effective solution to the poison pill batch problem. By ensuring that the state file is updated even for failed batches, the agent will be able to move past problematic log entries after a restart, improving its resilience and reliability.

The implementation requires minimal changes to the existing codebase, focusing only on the `sender.Send()` method in the pusher component. This approach maintains the existing architecture and interfaces while addressing the specific issue of poison pill batches.