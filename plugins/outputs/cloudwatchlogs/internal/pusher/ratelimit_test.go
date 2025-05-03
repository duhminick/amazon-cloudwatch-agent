package pusher

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

type ResponseCode int

const (
	ReplenishRate = 1
	MaxTPS = 10

	TotalOperations = 15

	SUCCESS	ResponseCode = 0
	THROTTLED ResponseCode = 1
)

type testScaffolding struct {
	// Throttling is generally by account ID so have a global limiter
	serviceLimiter *rate.Limiter

	attempted chan bool
	succeeded chan bool
	clientThrottled chan bool
	serviceThrottled chan bool
	timeSpent chan time.Duration
}

func TestServiceThrottling(t *testing.T) {
	s := testScaffolding {
		serviceLimiter: rate.NewLimiter(ReplenishRate, MaxTPS),
		attempted: make(chan bool, TotalOperations * numBackoffRetries),
		succeeded: make(chan bool, TotalOperations),
		clientThrottled: make(chan bool, TotalOperations * numBackoffRetries),
		serviceThrottled: make(chan bool, TotalOperations * numBackoffRetries),
		timeSpent: make(chan time.Duration, TotalOperations),
	}
	var wg sync.WaitGroup
	for i := 0; i < TotalOperations; i++ {
		wg.Add(1)
		go func() {
			startTime := time.Now()
			defer wg.Done()
			succeeded := false
			for attempt := 0; attempt < numBackoffRetries; attempt++ {
				// API calls are not instant so bake in some round trip time
				time.Sleep(randomRoundTripDuration())
				resp := s.runServiceOperation()
				s.attempted <- true

				if resp == THROTTLED {
					s.serviceThrottled <- true
					delay := calculateDelay(attempt)
					// fmt.Printf("delaying for %v on attempt %d\n", delay, attempt)
					time.Sleep(delay)
					continue
				}

				succeeded = true
				// fmt.Printf("succeeded operation\n")
				break
			}
			s.succeeded <- succeeded
			s.timeSpent <- time.Now().Sub(startTime)
		}()
	}
	wg.Wait()

	attempted := boolChannelSum(s.attempted)
	fmt.Printf("attempted calls: %d\n", attempted)
	fmt.Printf("service throttled calls: %d\n", boolChannelSum(s.serviceThrottled))
	fmt.Printf("client throttled calls: %d\n", boolChannelSum(s.clientThrottled))

	successful := boolChannelSum(s.succeeded)
	fmt.Printf("successful operations: %d\n", successful)
	fmt.Printf("failed operations: %d\n", TotalOperations - successful)

	fmt.Printf("total time spent: %v\n", durationChannelSum(s.timeSpent))
}

func (s *testScaffolding) runServiceOperation() ResponseCode {
	if ok := s.serviceLimiter.Allow(); !ok {
		return THROTTLED
	}
	return SUCCESS
}

func calculateDelay(retryCount int) time.Duration {
	delay := baseRetryDelay
	if retryCount < numBackoffRetries {
		delay = baseRetryDelay * time.Duration(1<<int64(retryCount))
	}
	if delay > maxRetryDelayTarget {
		delay = maxRetryDelayTarget
	}
	return withJitter(delay)
}

func randomRoundTripDuration() time.Duration {
	return time.Duration(rand.Int63n(500))
}

func boolChannelSum(c chan bool) int {
	close(c)
	sum := 0
	for b := range c {
		if b {
			sum += 1
		}
	}
	return sum
}

func durationChannelSum(c chan time.Duration) time.Duration {
	close(c)
	var sum time.Duration
	for d := range c {
		sum += d
	}
	return sum
}
