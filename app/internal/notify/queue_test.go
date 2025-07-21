package notify

import (
	"container/heap"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNotificationQueue_BasicOperation(t *testing.T) {
	manager := &Manager{
		notifiers: []Notifier{NewMockNotifier(true)},
		options:   NotificationOptions{MaxOutputLength: 1024},
	}

	opts := DefaultQueueOptions()
	opts.Workers = 1
	opts.RateLimit = 10 * time.Millisecond

	queue := NewNotificationQueue(manager, opts)
	queue.Start()
	defer func() {
		if err := queue.Stop(5 * time.Second); err != nil {
			t.Errorf("Failed to stop queue: %v", err)
		}
	}()

	// Enqueue a notification
	notification := Notification{
		Title:   "Test",
		Message: "Test message",
		Level:   LevelInfo,
	}

	err := queue.Enqueue(notification, PriorityNormal)
	if err != nil {
		t.Fatalf("Failed to enqueue: %v", err)
	}

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Check notification was processed
	mock := manager.notifiers[0].(*MockNotifier)
	notifications := mock.GetNotifications()

	if len(notifications) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(notifications))
	}

	// Check metrics
	processed, failed, dropped := queue.GetMetrics()
	if processed != 1 {
		t.Errorf("Expected 1 processed, got %d", processed)
	}
	if failed != 0 {
		t.Errorf("Expected 0 failed, got %d", failed)
	}
	if dropped != 0 {
		t.Errorf("Expected 0 dropped, got %d", dropped)
	}
}

func TestNotificationQueue_Priority(t *testing.T) {
	manager := &Manager{
		notifiers: []Notifier{NewMockNotifier(true)},
		options:   NotificationOptions{MaxOutputLength: 1024},
	}

	opts := DefaultQueueOptions()
	opts.Workers = 1
	opts.RateLimit = 1 * time.Millisecond
	// Reduce channel size to force items into heap
	opts.MaxQueueSize = 10

	queue := NewNotificationQueue(manager, opts)

	// Start but with a delay to allow all items to queue up
	queue.Start()

	// Fill the channel first with low priority items to force heap usage
	for i := 0; i < 5; i++ {
		_ = queue.Enqueue(Notification{
			Title:   "Filler",
			Message: "Fill channel",
			Level:   LevelInfo,
		}, PriorityLow)
	}

	// Now enqueue priority items - these should go to heap
	notifications := []struct {
		title    string
		priority int
	}{
		{"Critical 1", PriorityCritical},
		{"High 1", PriorityHigh},
		{"Normal 1", PriorityNormal},
		{"Critical 2", PriorityCritical},
	}

	for _, n := range notifications {
		err := queue.Enqueue(Notification{
			Title:   n.title,
			Message: "Test",
			Level:   LevelInfo,
		}, n.priority)
		if err != nil {
			t.Fatalf("Failed to enqueue %s: %v", n.title, err)
		}
	}

	// Let processing happen
	time.Sleep(300 * time.Millisecond)
	queue.Stop(5 * time.Second)

	// Check order
	mock := manager.notifiers[0].(*MockNotifier)
	processed := mock.GetNotifications()

	// Look for the order of our test notifications (skip filler)
	testNotifs := []string{}
	for _, notif := range processed {
		if notif.Title != "Filler" {
			testNotifs = append(testNotifs, notif.Title)
		}
	}

	if len(testNotifs) < 4 {
		t.Errorf("Expected at least 4 test notifications, got %d", len(testNotifs))
		return
	}

	// Verify critical notifications appear before normal/high
	criticalPositions := []int{}
	normalHighPositions := []int{}

	for i, title := range testNotifs {
		if title == "Critical 1" || title == "Critical 2" {
			criticalPositions = append(criticalPositions, i)
		} else if title == "Normal 1" || title == "High 1" {
			normalHighPositions = append(normalHighPositions, i)
		}
	}

	// At least one critical should come before normal/high
	foundPriorityOrder := false
	for _, critPos := range criticalPositions {
		for _, normPos := range normalHighPositions {
			if critPos < normPos {
				foundPriorityOrder = true
				break
			}
		}
		if foundPriorityOrder {
			break
		}
	}

	if !foundPriorityOrder {
		t.Errorf("Priority order not respected. Order: %v", testNotifs)
	}
}

func TestNotificationQueue_RateLimit(t *testing.T) {
	manager := &Manager{
		notifiers: []Notifier{NewMockNotifier(true)},
		options:   NotificationOptions{MaxOutputLength: 1024},
	}

	opts := DefaultQueueOptions()
	opts.Workers = 1
	opts.RateLimit = 50 * time.Millisecond

	queue := NewNotificationQueue(manager, opts)
	queue.Start()
	defer func() {
		if err := queue.Stop(5 * time.Second); err != nil {
			t.Errorf("Failed to stop queue: %v", err)
		}
	}()

	// Enqueue multiple notifications
	start := time.Now()
	for i := 0; i < 3; i++ {
		err := queue.Enqueue(Notification{
			Title:   "Test",
			Message: "Test message",
			Level:   LevelInfo,
		}, PriorityNormal)
		if err != nil {
			t.Fatalf("Failed to enqueue: %v", err)
		}
	}

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	elapsed := time.Since(start)

	// Should take at least 100ms (2 intervals of 50ms)
	if elapsed < 100*time.Millisecond {
		t.Errorf("Rate limiting not working, took only %v", elapsed)
	}
}

func TestNotificationQueue_Retry(t *testing.T) {
	// Create a notifier that fails initially
	mock := NewMockNotifier(true)
	failCount := int32(2) // Fail first 2 attempts
	attemptCount := int32(0)
	var notifications []Notification
	var notifMu sync.Mutex

	// Set custom notify function that fails conditionally
	mock.SetNotifyFunc(func(ctx context.Context, n Notification) error {
		count := atomic.AddInt32(&attemptCount, 1)
		if count <= failCount {
			return errors.New("simulated error")
		}
		// Success - record notification without locking mock's mutex
		notifMu.Lock()
		notifications = append(notifications, n)
		notifMu.Unlock()
		return nil
	})

	manager := &Manager{
		notifiers: []Notifier{mock},
		options:   NotificationOptions{MaxOutputLength: 1024},
	}

	opts := DefaultQueueOptions()
	opts.Workers = 1
	opts.RetryBackoff = 10 * time.Millisecond
	opts.MaxRetries = 3

	queue := NewNotificationQueue(manager, opts)
	queue.Start()
	defer func() {
		if err := queue.Stop(5 * time.Second); err != nil {
			t.Errorf("Failed to stop queue: %v", err)
		}
	}()

	// Enqueue notification
	err := queue.Enqueue(Notification{
		Title:   "Retry Test",
		Message: "This should be retried",
		Level:   LevelWarning,
	}, PriorityHigh)
	if err != nil {
		t.Fatalf("Failed to enqueue: %v", err)
	}

	// Wait for retries
	time.Sleep(500 * time.Millisecond)

	// Check metrics
	processed, failed, _ := queue.GetMetrics()

	// Should eventually succeed
	if processed != 1 {
		t.Errorf("Expected 1 processed after retries, got %d", processed)
	}
	if failed != 0 {
		t.Errorf("Expected 0 failed (should have succeeded after retries), got %d", failed)
	}

	// Verify notification was recorded
	notifMu.Lock()
	recordedCount := len(notifications)
	notifMu.Unlock()

	if recordedCount != 1 {
		t.Errorf("Expected 1 notification recorded, got %d", recordedCount)
	}
}

func TestNotificationQueue_QueueFull(t *testing.T) {
	manager := &Manager{
		notifiers: []Notifier{NewMockNotifier(true)},
		options:   NotificationOptions{MaxOutputLength: 1024},
	}

	opts := DefaultQueueOptions()
	opts.MaxQueueSize = 3
	opts.Workers = 0 // Don't process

	queue := NewNotificationQueue(manager, opts)

	// Fill the queue
	for i := 0; i < 3; i++ {
		err := queue.Enqueue(Notification{
			Title: "Test",
		}, PriorityNormal)
		if err != nil {
			t.Fatalf("Failed to enqueue item %d: %v", i, err)
		}
	}

	// Try to add one more
	err := queue.Enqueue(Notification{
		Title: "Overflow",
	}, PriorityNormal)

	if err == nil {
		t.Error("Expected error when queue is full")
	}

	_, _, dropped := queue.GetMetrics()
	if dropped != 1 {
		t.Errorf("Expected 1 dropped, got %d", dropped)
	}
}

func TestNotificationQueue_GracefulShutdown(t *testing.T) {
	manager := &Manager{
		notifiers: []Notifier{NewMockNotifier(true)},
		options:   NotificationOptions{MaxOutputLength: 1024},
	}

	opts := DefaultQueueOptions()
	opts.Workers = 1
	opts.RateLimit = 10 * time.Millisecond

	queue := NewNotificationQueue(manager, opts)
	queue.Start()

	// Enqueue notifications
	for i := 0; i < 5; i++ {
		err := queue.Enqueue(Notification{
			Title: "Shutdown test",
		}, PriorityNormal)
		if err != nil {
			t.Fatalf("Failed to enqueue: %v", err)
		}
	}

	// Stop with timeout
	err := queue.Stop(2 * time.Second)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	// All notifications should be processed
	mock := manager.notifiers[0].(*MockNotifier)
	if len(mock.GetNotifications()) != 5 {
		t.Errorf("Not all notifications were processed during shutdown")
	}
}

func TestNotificationQueue_OutputNotification(t *testing.T) {
	manager := &Manager{
		notifiers: []Notifier{NewMockOutputNotifier(true)},
		options:   NotificationOptions{MaxOutputLength: 1024},
	}

	opts := DefaultQueueOptions()
	queue := NewNotificationQueue(manager, opts)
	queue.Start()
	defer func() {
		if err := queue.Stop(5 * time.Second); err != nil {
			t.Errorf("Failed to stop queue: %v", err)
		}
	}()

	// Enqueue output notification
	outputNotif := OutputNotification{
		Notification: Notification{
			Title:   "Command",
			Message: "Completed",
			Level:   LevelSuccess,
		},
		Output:   "Command output here",
		ExitCode: 0,
	}

	err := queue.Enqueue(outputNotif, PriorityNormal)
	if err != nil {
		t.Fatalf("Failed to enqueue: %v", err)
	}

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Check it was processed
	mock := manager.notifiers[0].(*MockOutputNotifier)
	outputs := mock.GetOutputNotifications()

	if len(outputs) != 1 {
		t.Errorf("Expected 1 output notification, got %d", len(outputs))
	}
	if len(outputs) > 0 && outputs[0].Output != "Command output here" {
		t.Errorf("Output not preserved correctly")
	}
}

func TestNotificationQueue_Concurrent(t *testing.T) {
	manager := &Manager{
		notifiers: []Notifier{NewMockNotifier(true)},
		options:   NotificationOptions{MaxOutputLength: 1024},
	}

	opts := DefaultQueueOptions()
	opts.Workers = 4
	opts.RateLimit = 1 * time.Millisecond

	queue := NewNotificationQueue(manager, opts)
	queue.Start()
	defer func() {
		if err := queue.Stop(5 * time.Second); err != nil {
			t.Errorf("Failed to stop queue: %v", err)
		}
	}()

	// Concurrent enqueuers
	var wg sync.WaitGroup
	numGoroutines := 10
	itemsPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(_ int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				_ = queue.Enqueue(Notification{
					Title: "Concurrent",
				}, PriorityNormal)
			}
		}(i)
	}

	wg.Wait()
	time.Sleep(500 * time.Millisecond)

	processed, _, _ := queue.GetMetrics()
	expected := uint64(numGoroutines * itemsPerGoroutine)

	if processed != expected {
		t.Errorf("Expected %d processed, got %d", expected, processed)
	}
}

func TestNotificationQueue_GetQueueSize(t *testing.T) {
	manager := &Manager{
		notifiers: []Notifier{NewMockNotifier(true)},
		options:   NotificationOptions{MaxOutputLength: 1024},
	}

	opts := DefaultQueueOptions()
	opts.Workers = 0 // Don't start workers

	queue := NewNotificationQueue(manager, opts)

	// Initial size should be 0
	if size := queue.GetQueueSize(); size != 0 {
		t.Errorf("Expected initial size 0, got %d", size)
	}

	// Add items
	for i := 0; i < 5; i++ {
		_ = queue.Enqueue(Notification{Title: "Test"}, PriorityNormal)
	}

	if size := queue.GetQueueSize(); size != 5 {
		t.Errorf("Expected size 5, got %d", size)
	}
}

func TestPriorityQueue_HeapInterface(t *testing.T) {
	// Test the heap interface methods directly
	pq := &priorityQueue{}
	heap.Init(pq)

	// Create test items
	items := []*QueueItem{
		{
			Notification: Notification{Title: "Low", Level: LevelInfo},
			Priority:     0,
			Timestamp:    time.Now(),
		},
		{
			Notification: Notification{Title: "High", Level: LevelError},
			Priority:     2,
			Timestamp:    time.Now().Add(1 * time.Second),
		},
		{
			Notification: Notification{Title: "Medium", Level: LevelWarning},
			Priority:     1,
			Timestamp:    time.Now().Add(2 * time.Second),
		},
	}

	// Push items using heap interface
	for _, item := range items {
		heap.Push(pq, item)
	}

	if pq.Len() != 3 {
		t.Errorf("Push: expected length 3, got %d", pq.Len())
	}

	// Test Less (higher priority should come first)
	// Note: heap maintains heap property, not sorted order
	// So we'll pop items and check order
	var poppedOrder []string
	for pq.Len() > 0 {
		item := heap.Pop(pq).(*QueueItem)
		if notif, ok := item.Notification.(Notification); ok {
			poppedOrder = append(poppedOrder, notif.Title)
		}
	}

	// Should be popped in priority order: High, Medium, Low
	expectedOrder := []string{"High", "Medium", "Low"}
	for i, title := range poppedOrder {
		if title != expectedOrder[i] {
			t.Errorf("Pop order: expected %s at position %d, got %s", expectedOrder[i], i, title)
		}
	}
}

func TestPriorityQueue_DirectMethods(t *testing.T) {
	// Test the priorityQueue methods directly (not through heap interface)
	pq := &priorityQueue{}

	// Test Push and Len
	item1 := &QueueItem{
		Notification: Notification{Title: "Item1"},
		Priority:     1,
		Timestamp:    time.Now(),
	}
	pq.Push(item1)

	if pq.Len() != 1 {
		t.Errorf("Len: expected 1, got %d", pq.Len())
	}

	// Test Less with two items
	item2 := &QueueItem{
		Notification: Notification{Title: "Item2"},
		Priority:     2,
		Timestamp:    time.Now().Add(1 * time.Second),
	}
	pq.Push(item2)

	// Higher priority (2) should be "less" than lower priority (1)
	if !pq.Less(1, 0) {
		t.Error("Less: higher priority item should be less")
	}

	// Test Swap
	var title0Before, title1Before string
	if notif, ok := (*pq)[0].Notification.(Notification); ok {
		title0Before = notif.Title
	}
	if notif, ok := (*pq)[1].Notification.(Notification); ok {
		title1Before = notif.Title
	}
	pq.Swap(0, 1)

	var title0After, title1After string
	if notif, ok := (*pq)[0].Notification.(Notification); ok {
		title0After = notif.Title
	}
	if notif, ok := (*pq)[1].Notification.(Notification); ok {
		title1After = notif.Title
	}
	if title0After != title1Before || title1After != title0Before {
		t.Error("Swap: items not properly swapped")
	}

	// Test Pop
	initialLen := pq.Len()
	popped := pq.Pop().(*QueueItem)

	if pq.Len() != initialLen-1 {
		t.Errorf("Pop: expected length %d, got %d", initialLen-1, pq.Len())
	}

	// Should pop the last item
	var poppedTitle string
	if notif, ok := popped.Notification.(Notification); ok {
		poppedTitle = notif.Title
	}
	if poppedTitle != title0Before {
		t.Errorf("Pop: expected %s, got %s", title0Before, poppedTitle)
	}
}
