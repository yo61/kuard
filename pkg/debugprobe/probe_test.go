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
package debugprobe

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func serve(t *testing.T, p *Probe) int {
	t.Helper()
	w := httptest.NewRecorder()
	p.Handle(w, httptest.NewRequest(http.MethodGet, "/healthy", nil), httprouter.Params{})
	return w.Code
}

// status reads the probe through the same JSON API the UI uses. Note the handler returns history
// newest-first, so index 0 is the most recent call.
func status(t *testing.T, p *Probe) ProbeStatus {
	t.Helper()
	w := httptest.NewRecorder()
	p.APIGet(w, httptest.NewRequest(http.MethodGet, "/healthy/api", nil), httprouter.Params{})

	var s ProbeStatus
	if err := json.Unmarshal(w.Body.Bytes(), &s); err != nil {
		t.Fatalf("decoding probe status: %v", err)
	}
	return s
}

// FailNext is the knob the book's probe demos turn: 0 succeeds forever, N fails N times then
// recovers, negative fails permanently.
func TestFailNextSemantics(t *testing.T) {
	tests := []struct {
		name     string
		failNext int
		want     []int
	}{
		{
			name:     "zero succeeds forever",
			failNext: 0,
			want:     []int{200, 200, 200, 200},
		},
		{
			name:     "positive fails that many times then recovers",
			failNext: 3,
			want:     []int{500, 500, 500, 200, 200},
		},
		{
			name:     "one fails exactly once",
			failNext: 1,
			want:     []int{500, 200},
		},
		{
			name:     "negative fails permanently",
			failNext: -1,
			want:     []int{500, 500, 500, 500},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.SetConfig(ProbeConfig{FailNext: tt.failNext})

			for i, want := range tt.want {
				if got := serve(t, p); got != want {
					t.Errorf("call %d = %d, want %d", i+1, got, want)
				}
			}
		})
	}
}

func TestHistoryRecordsEveryCall(t *testing.T) {
	p := New()
	p.SetConfig(ProbeConfig{FailNext: 2})

	for range 5 {
		serve(t, p)
	}

	s := status(t, p)
	if len(s.History) != 5 {
		t.Fatalf("history has %d entries, want 5", len(s.History))
	}

	// Newest first, so the two failures are at the end.
	wantCodes := []int{200, 200, 200, 500, 500}
	for i, want := range wantCodes {
		if got := s.History[i].Code; got != want {
			t.Errorf("history[%d].Code = %d, want %d", i, got, want)
		}
	}
}

// The history is a fixed-size window; without the trim a long-running probe would grow forever.
func TestHistoryIsCappedAtMaxHistory(t *testing.T) {
	p := New()

	for range maxHistory + 10 {
		serve(t, p)
	}

	if got := len(status(t, p).History); got != maxHistory {
		t.Errorf("history length = %d, want it capped at %d", got, maxHistory)
	}
}

func TestFailNextCountsDownInReportedStatus(t *testing.T) {
	p := New()
	p.SetConfig(ProbeConfig{FailNext: 2})

	if got := status(t, p).FailNext; got != 2 {
		t.Fatalf("FailNext before any call = %d, want 2", got)
	}

	serve(t, p)

	if got := status(t, p).FailNext; got != 1 {
		t.Errorf("FailNext after one call = %d, want 1", got)
	}
}
