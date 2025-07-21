package notify

import (
	"container/heap"
	"context"
	"fmt"
	"sync"
	"time"
)

// Priority levels for queue items
const (
	PriorityLow      = 0
	PriorityNormal   = 1
	PriorityHigh     = 2
	PriorityCritical = 3
)

// QueueItem represents a notification in the queue
type QueueItem struct {
	Notification interface{} // Can be Notification or OutputNotification
	Priority     int
	Timestamp    time.Time
	RetryCount   int
	MaxRetries   int
	index        int // Used by heap interface
}

// NotificationQueue manages asynchronous notification delivery
type NotificationQueue struct {
	mu      sync.Mutex
	items   priorityQueue
	ch      chan *QueueItem
	manager *Manager
	ctx     context.Context //nolint:containedctx // Required for queue lifecycle management
	cancel  context.CancelFunc
	wg      sync.WaitGroup

	// Configuration
	maxQueueSize int
	workers      int
	rateLimit    time.Duration
	retryBackoff time.Duration
	maxRetries   int

	// Metrics
	processed      uint64
	failed         uint64
	dropped        uint64
	lastNotifyTime time.Time
}

// QueueOptions configures the notification queue
type QueueOptions struct {
	MaxQueueSize int           // Maximum items in queue (default: 1000)
	Workers      int           // Number of worker goroutines (default: 2)
	RateLimit    time.Duration // Minimum time between notifications (default: 100ms)
	RetryBackoff time.Duration // Base backoff for retries (default: 1s)
	MaxRetries   int           // Maximum retry attempts (default: 3)
}

// DefaultQueueOptions returns default queue options
func DefaultQueueOptions() QueueOptions {
	return QueueOptions{
		MaxQueueSize: 1000,
		Workers:      2,
		RateLimit:    100 * time.Millisecond,
		RetryBackoff: time.Second,
		MaxRetries:   3,
	}
}

// NewNotificationQueue creates a new notification queue
func NewNotificationQueue(manager *Manager, opts QueueOptions) *NotificationQueue {
	if opts.MaxQueueSize <= 0 {
		opts.MaxQueueSize = 1000
	}
	if opts.Workers <= 0 {
		opts.Workers = 2
	}
	if opts.RateLimit <= 0 {
		opts.RateLimit = 100 * time.Millisecond
	}
	if opts.RetryBackoff <= 0 {
		opts.RetryBackoff = time.Second
	}
	if opts.MaxRetries < 0 {
		opts.MaxRetries = 3
	}

	ctx, cancel := context.WithCancel(context.Background())

	q := &NotificationQueue{
		items:        make(priorityQueue, 0),
		ch:           make(chan *QueueItem, opts.MaxQueueSize),
		manager:      manager,
		ctx:          ctx,
		cancel:       cancel,
		maxQueueSize: opts.MaxQueueSize,
		workers:      opts.Workers,
		rateLimit:    opts.RateLimit,
		retryBackoff: opts.RetryBackoff,
		maxRetries:   opts.MaxRetries,
	}

	heap.Init(&q.items)
	return q
}

// Start begins processing the queue
func (q *NotificationQueue) Start() {
	// Start dispatcher
	q.wg.Add(1)
	go q.dispatcher()

	// Start workers
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
}

// Stop gracefully shuts down the queue
func (q *NotificationQueue) Stop(timeout time.Duration) error {
	// Signal shutdown
	q.cancel()

	// Wait for workers with timeout
	done := make(chan struct{})
	go func() {
		q.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("queue shutdown timed out after %v", timeout)
	}
}

// Enqueue adds a notification to the queue
func (q *NotificationQueue) Enqueue(notification interface{}, priority int) error {
	return q.EnqueueWithRetries(notification, priority, q.maxRetries)
}

// EnqueueWithRetries adds a notification with custom retry count
func (q *NotificationQueue) EnqueueWithRetries(notification interface{}, priority int, maxRetries int) error {
	item := &QueueItem{
		Notification: notification,
		Priority:     priority,
		Timestamp:    time.Now(),
		MaxRetries:   maxRetries,
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	// Check total queue size (heap + channel)
	totalSize := len(q.items) + len(q.ch)
	if totalSize >= q.maxQueueSize {
		q.dropped++
		return fmt.Errorf("queue is full (%d items)", q.maxQueueSize)
	}

	// Try to send directly to channel first
	select {
	case q.ch <- item:
		// Successfully sent to channel
		return nil
	default:
		// Channel is full, add to heap
		heap.Push(&q.items, item)
		return nil
	}
}

// dispatcher moves items from heap to channel
func (q *NotificationQueue) dispatcher() {
	defer q.wg.Done()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-q.ctx.Done():
			// Drain remaining items
			q.drainQueue()
			return

		case <-ticker.C:
			q.mu.Lock()
			// Move items from heap to channel
			for len(q.items) > 0 {
				poppedItem := heap.Pop(&q.items)
				item, ok := poppedItem.(*QueueItem)
				if !ok {
					// This should never happen, log and continue
					continue
				}
				select {
				case q.ch <- item:
					// Successfully sent
				default:
					// Channel is full, push item back
					heap.Push(&q.items, item)
					break
				}
			}
			q.mu.Unlock()
		}
	}
}

// worker processes notifications from the queue
func (q *NotificationQueue) worker(_ int) {
	defer q.wg.Done()

	for {
		select {
		case <-q.ctx.Done():
			return

		case item := <-q.ch:
			q.processItem(item)
		}
	}
}

// processItem handles a single notification
func (q *NotificationQueue) processItem(item *QueueItem) {
	// Apply rate limiting
	q.applyRateLimit()

	ctx, cancel := context.WithTimeout(q.ctx, 30*time.Second)
	defer cancel()

	var err error

	// Process based on notification type
	switch n := item.Notification.(type) {
	case Notification:
		err = q.manager.Notify(ctx, n)
	case OutputNotification:
		err = q.manager.NotifyWithOutput(ctx, n)
	default:
		err = fmt.Errorf("unknown notification type: %T", n)
	}

	if err != nil {
		q.handleError(item, err)
	} else {
		q.mu.Lock()
		q.processed++
		q.mu.Unlock()
	}
}

// handleError handles notification delivery errors
func (q *NotificationQueue) handleError(item *QueueItem, _ error) {
	item.RetryCount++

	if item.RetryCount <= item.MaxRetries {
		// Calculate backoff
		backoff := q.retryBackoff * time.Duration(item.RetryCount)

		// Schedule retry
		time.AfterFunc(backoff, func() {
			if err := q.EnqueueWithRetries(item.Notification, item.Priority, item.MaxRetries-item.RetryCount); err != nil {
				q.mu.Lock()
				q.failed++
				q.mu.Unlock()
			}
		})
	} else {
		q.mu.Lock()
		q.failed++
		q.mu.Unlock()
	}
}

// applyRateLimit ensures minimum time between notifications
func (q *NotificationQueue) applyRateLimit() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.lastNotifyTime.IsZero() {
		elapsed := time.Since(q.lastNotifyTime)
		if elapsed < q.rateLimit {
			time.Sleep(q.rateLimit - elapsed)
		}
	}

	q.lastNotifyTime = time.Now()
}

// drainQueue processes all remaining items during shutdown
func (q *NotificationQueue) drainQueue() {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Create a context with timeout for draining
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Process all items in heap
	for len(q.items) > 0 {
		poppedItem := heap.Pop(&q.items)
		item, ok := poppedItem.(*QueueItem)
		if !ok {
			continue
		}

		switch n := item.Notification.(type) {
		case Notification:
			_ = q.manager.Notify(ctx, n)
		case OutputNotification:
			_ = q.manager.NotifyWithOutput(ctx, n)
		}

		q.processed++
	}

	// Process all items in channel
	for {
		select {
		case item := <-q.ch:
			switch n := item.Notification.(type) {
			case Notification:
				_ = q.manager.Notify(ctx, n)
			case OutputNotification:
				_ = q.manager.NotifyWithOutput(ctx, n)
			}
			q.processed++
		default:
			return
		}
	}
}

// GetMetrics returns queue metrics
func (q *NotificationQueue) GetMetrics() (processed, failed, dropped uint64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.processed, q.failed, q.dropped
}

// GetQueueSize returns the current queue size
func (q *NotificationQueue) GetQueueSize() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items) + len(q.ch)
}

// Priority queue implementation
type priorityQueue []*QueueItem

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// Higher priority first
	if pq[i].Priority != pq[j].Priority {
		return pq[i].Priority > pq[j].Priority
	}
	// Earlier timestamp first for same priority
	return pq[i].Timestamp.Before(pq[j].Timestamp)
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*QueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
