/*
Copyright 2017 The KUAR Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package memqserver

import (
	"errors"
	"testing"
)

func TestQueueRoundTripPreservesOrder(t *testing.T) {
	b := NewBroker()
	if err := b.CreateQueue("work"); err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}

	for _, body := range []string{"first", "second", "third"} {
		if _, err := b.PutMessage("work", body); err != nil {
			t.Fatalf("PutMessage(%q): %v", body, err)
		}
	}

	for _, want := range []string{"first", "second", "third"} {
		m, err := b.GetMessage("work")
		if err != nil {
			t.Fatalf("GetMessage: %v", err)
		}
		if m.Body != want {
			t.Errorf("GetMessage body = %q, want %q", m.Body, want)
		}
	}
}

func TestDequeueFromEmptyQueueReportsEmpty(t *testing.T) {
	b := NewBroker()
	if err := b.CreateQueue("work"); err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}

	if _, err := b.GetMessage("work"); !errors.Is(err, ErrEmptyQueue) {
		t.Errorf("GetMessage on empty queue = %v, want %v", err, ErrEmptyQueue)
	}
}

func TestOperationsOnMissingQueueReportNotExist(t *testing.T) {
	b := NewBroker()

	if _, err := b.PutMessage("nope", "x"); !errors.Is(err, ErrNotExist) {
		t.Errorf("PutMessage on missing queue = %v, want %v", err, ErrNotExist)
	}
	if _, err := b.GetMessage("nope"); !errors.Is(err, ErrNotExist) {
		t.Errorf("GetMessage on missing queue = %v, want %v", err, ErrNotExist)
	}
	if err := b.DrainQueue("nope"); !errors.Is(err, ErrNotExist) {
		t.Errorf("DrainQueue on missing queue = %v, want %v", err, ErrNotExist)
	}
	if err := b.DeleteQueue("nope"); !errors.Is(err, ErrNotExist) {
		t.Errorf("DeleteQueue on missing queue = %v, want %v", err, ErrNotExist)
	}
}

// Empty names are rejected by the HTTP layer (see Server.CreateQueue), not here, so the broker
// only owns the duplicate rule.
func TestCreateQueueRejectsDuplicates(t *testing.T) {
	b := NewBroker()

	if err := b.CreateQueue("work"); err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	if err := b.CreateQueue("work"); !errors.Is(err, ErrAlreadyExist) {
		t.Errorf("CreateQueue twice = %v, want %v", err, ErrAlreadyExist)
	}
}

func TestDrainDiscardsMessagesButKeepsQueue(t *testing.T) {
	b := NewBroker()
	if err := b.CreateQueue("work"); err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	for range 3 {
		if _, err := b.PutMessage("work", "x"); err != nil {
			t.Fatalf("PutMessage: %v", err)
		}
	}

	if err := b.DrainQueue("work"); err != nil {
		t.Fatalf("DrainQueue: %v", err)
	}

	// Draining empties the queue; it does not remove it.
	if _, err := b.GetMessage("work"); !errors.Is(err, ErrEmptyQueue) {
		t.Errorf("GetMessage after drain = %v, want %v", err, ErrEmptyQueue)
	}
}

func TestDeleteRemovesTheQueueEntirely(t *testing.T) {
	b := NewBroker()
	if err := b.CreateQueue("work"); err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	if err := b.DeleteQueue("work"); err != nil {
		t.Fatalf("DeleteQueue: %v", err)
	}

	if _, err := b.GetMessage("work"); !errors.Is(err, ErrNotExist) {
		t.Errorf("GetMessage after delete = %v, want %v", err, ErrNotExist)
	}
}

func TestStatsTrackDepthAndCounters(t *testing.T) {
	b := NewBroker()
	if err := b.CreateQueue("work"); err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}
	for range 3 {
		if _, err := b.PutMessage("work", "x"); err != nil {
			t.Fatalf("PutMessage: %v", err)
		}
	}
	if _, err := b.GetMessage("work"); err != nil {
		t.Fatalf("GetMessage: %v", err)
	}

	stats := b.Stats()
	if len(stats.Queues) != 1 {
		t.Fatalf("stats reported %d queues, want 1", len(stats.Queues))
	}

	q := stats.Queues[0]
	if q.Name != "work" {
		t.Errorf("queue name = %q, want %q", q.Name, "work")
	}
	if q.Enqueued != 3 {
		t.Errorf("enqueued = %d, want 3", q.Enqueued)
	}
	if q.Dequeued != 1 {
		t.Errorf("dequeued = %d, want 1", q.Dequeued)
	}
	if q.Depth != 2 {
		t.Errorf("depth = %d, want 2 (3 enqueued - 1 dequeued)", q.Depth)
	}
}

func TestMessagesGetDistinctIDs(t *testing.T) {
	b := NewBroker()
	if err := b.CreateQueue("work"); err != nil {
		t.Fatalf("CreateQueue: %v", err)
	}

	seen := map[string]bool{}
	for range 20 {
		m, err := b.PutMessage("work", "x")
		if err != nil {
			t.Fatalf("PutMessage: %v", err)
		}
		if m.ID == "" {
			t.Fatal("message got an empty ID")
		}
		if seen[m.ID] {
			t.Fatalf("duplicate message ID %q", m.ID)
		}
		seen[m.ID] = true
	}
}
