# Implementation Plan

- [ ] 1. Modify the logEventBatch struct to separate state-updating callbacks
  - Modify the batch.go file to separate state-updating callbacks from other success callbacks
  - Add a new field to track state-updating callbacks separately from other callbacks
  - _Requirements: 3.1, 3.2_

- [ ] 2. Add a new method to update state only
  - [x] 2.1 Implement updateStateOnly method in batch.go
    - Create a new method that only executes the state-updating callbacks
    - Ensure it properly handles the range queue batchers
    - _Requirements: 3.2_
  
  - [x] 2.2 Add unit tests for updateStateOnly method
    - Create test cases to verify the method correctly updates state without executing other callbacks
    - _Requirements: 3.2_

- [x] 3. Modify the sender to update state for failed batches
  - [x] 3.1 Update the Send method in sender.go
    - Modify the retry exhaustion logic to call batch.updateStateOnly() instead of just returning
    - Add debug logging to indicate state is being updated for a failed batch
    - _Requirements: 1.1, 1.2, 1.3, 1.4_
  
  - [x] 3.2 Add unit tests for the modified Send method
    - Create test cases to verify state is updated when retry attempts are exhausted
    - _Requirements: 1.1, 1.2_

- [x] 4. Update the queue implementation to use the new approach
  - [x] 4.1 Modify the queue.go file to register state callbacks separately
    - Update the converter.convert method to identify state-updating callbacks
    - Ensure the queue.send method registers callbacks appropriately
    - _Requirements: 3.2, 3.3_
  
  - [x] 4.2 Add unit tests for the modified queue implementation
    - Create test cases to verify callbacks are registered correctly
    - _Requirements: 3.2, 3.3_

