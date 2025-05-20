package pusher

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/amazon-cloudwatch-agent/internal/logscommon"
)

const (
	ttlTime = 5 * time.Minute
)

type retentionPolicyTTL struct {
	statePath     string
	// oldTimestamps come from the TTL file on agent start
	oldTimestamps map[string]time.Time
	// newTimestamps are the new TTLs that will be saved periodically and when the agent is done
	newTimestamps map[string]time.Time
	mu            sync.RWMutex
	ch            chan string
	done          chan struct{}
}

func NewRetentionPolicyTTL(fileStatePath string) *retentionPolicyTTL {
	r := &retentionPolicyTTL{
		statePath:     filepath.Join(fileStatePath, logscommon.RetentionPolicyTTLFileName),
		oldTimestamps: make(map[string]time.Time),
		newTimestamps: make(map[string]time.Time),
		ch:            make(chan string, retentionChannelSize),
		done:          make(chan struct{}),
	}

	r.loadTTLState()
	go r.process()
	return r
}

func (r *retentionPolicyTTL) Update(group string) {
	r.ch <- group
}

func (r *retentionPolicyTTL) Done() {
	<-r.done
}

func (r *retentionPolicyTTL) IsExpired(group string) bool {
	if ts, ok := r.oldTimestamps[escapeLogGroup(group)]; ok {
		return ts.Add(ttlTime).After(time.Now())
	}
	return false
}

func (r *retentionPolicyTTL) loadTTLState() {
}

func (r *retentionPolicyTTL) process() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	for {
		select {
		case group := <-r.ch:
			r.updateTimestamp(group)
		case <-t.C:
			r.saveTTLState()
		case <-r.done:
			r.saveTTLState()
			return
		}
	}
}

func (r *retentionPolicyTTL) updateTimestamp(group string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.newTimestamps[escapeLogGroup(group)] = time.Now()
}

func (r *retentionPolicyTTL) saveTTLState() {
	r.mu.RLock()
	defer r.mu.Unlock()
}

func escapeLogGroup(group string) (escapedLogGroup string) {
	escapedLogGroup = filepath.ToSlash(group)
	escapedLogGroup = strings.Replace(escapedLogGroup, "/", "_", -1)
	escapedLogGroup = strings.Replace(escapedLogGroup, " ", "_", -1)
	escapedLogGroup = strings.Replace(escapedLogGroup, ":", "_", -1)
	return
}
