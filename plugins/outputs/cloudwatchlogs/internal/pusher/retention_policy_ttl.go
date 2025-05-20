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
	statePath  string
	timestamps map[string]time.Time
	mu         sync.RWMutex
	ch         chan string
	done       chan struct{}
}

func NewRetentionPolicyTTL(fileStatePath string) *retentionPolicyTTL {
	r := &retentionPolicyTTL{
		statePath:  filepath.Join(fileStatePath, logscommon.RetentionPolicyTTLFileName),
		timestamps: make(map[string]time.Time),
		ch:         make(chan string, retentionChannelSize),
		done:       make(chan struct{}),
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

// DOMINIC: this should read from the state file when loaded
// that way when we save, it'll just have the new log groups that we've seen
// this is for the case when a log group was removed by the user. shouldn't need
// a mutex then
func (r *retentionPolicyTTL) IsExpired(group string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if ts, ok := r.timestamps[escapeFilePath(group)]; ok {
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
			r.saveState()
		case <-r.done:
			r.saveState()
			return
		}
	}
}

func (r *retentionPolicyTTL) updateTimestamp(group string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.timestamps[group] = time.Now()
}

func (r *retentionPolicyTTL) saveState() {
	r.mu.RLock()
	defer r.mu.Unlock()
}

func escapeFilePath(filePath string) (escapedFilePath string) {
	escapedFilePath = filepath.ToSlash(filePath)
	escapedFilePath = strings.Replace(escapedFilePath, "/", "_", -1)
	escapedFilePath = strings.Replace(escapedFilePath, " ", "_", -1)
	escapedFilePath = strings.Replace(escapedFilePath, ":", "_", -1)
	return
}
