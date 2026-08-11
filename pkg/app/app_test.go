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
package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// newTestHandler wires an App the way cmd/kuard does -- BindConfig, parse, LoadConfig -- but with
// a private viper and FlagSet so tests do not fight over the package-level globals.
func newTestHandler(t *testing.T, args ...string) http.Handler {
	t.Helper()

	k := NewApp()
	v := viper.New()
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	k.BindConfig(v, fs)
	if err := fs.Parse(args); err != nil {
		t.Fatalf("parsing %v: %v", args, err)
	}
	k.LoadConfig(v)

	return k.Handler()
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
}

func TestProbesReportHealthy(t *testing.T) {
	h := newTestHandler(t)

	for _, path := range []string{"/healthy", "/ready"} {
		w := get(t, h, path)
		if w.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want %d", path, w.Code, http.StatusOK)
		}
		if got := w.Body.String(); got != "ok" {
			t.Errorf("GET %s body = %q, want %q", path, got, "ok")
		}
	}
}

// The metrics endpoint moved from prometheus.Handler() to promhttp.Handler() when client_golang
// reached 1.0. Asserting on the exported series, not just the status code, is what makes this
// catch a regression rather than a 200 from an empty handler.
func TestMetricsExportsRequestDuration(t *testing.T) {
	h := newTestHandler(t)

	// The histogram is a Vec, so it exports nothing until a request has been observed.
	get(t, h, "/healthy")

	w := get(t, h, "/metrics")
	if w.Code != http.StatusOK {
		t.Fatalf("GET /metrics = %d, want %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "request_duration_seconds") {
		t.Error("GET /metrics did not export request_duration_seconds")
	}
}

func TestEnvAPIReturnsCommandLineAndEnv(t *testing.T) {
	h := newTestHandler(t)

	w := get(t, h, "/env/api")
	if w.Code != http.StatusOK {
		t.Fatalf("GET /env/api = %d, want %d", w.Code, http.StatusOK)
	}

	var body struct {
		CommandLine []string          `json:"commandLine"`
		Env         map[string]string `json:"env"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding /env/api: %v", err)
	}
	if len(body.CommandLine) == 0 {
		t.Error("/env/api returned an empty commandLine")
	}
	if body.Env == nil {
		t.Error("/env/api returned no env map")
	}
}

func TestMemAPIReturnsMemStats(t *testing.T) {
	h := newTestHandler(t)

	w := get(t, h, "/mem/api")
	if w.Code != http.StatusOK {
		t.Fatalf("GET /mem/api = %d, want %d", w.Code, http.StatusOK)
	}

	var body struct {
		MemStats struct {
			HeapAlloc uint64
		} `json:"memStats"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding /mem/api: %v", err)
	}
	if body.MemStats.HeapAlloc == 0 {
		t.Error("/mem/api reported HeapAlloc of 0, which cannot be right for a running process")
	}
}

// The React bundle is mounted by the server-rendered shell, so this covers the template and the
// go-bindata assets as well as the route.
func TestRootServesReactShell(t *testing.T) {
	h := newTestHandler(t)

	w := get(t, h, "/")
	if w.Code != http.StatusOK {
		t.Fatalf("GET / = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	for _, want := range []string{`id="root"`, "pageContext", "built/bundle.js"} {
		if !strings.Contains(body, want) {
			t.Errorf("GET / did not contain %q", want)
		}
	}
}

// kuard serves the whole app under "", /a, /b and /c. The client passes urlBase to React Router as
// its basename, so losing a prefix breaks routing rather than just a URL.
func TestAllURLBasePrefixesServeTheApp(t *testing.T) {
	h := newTestHandler(t)

	for _, prefix := range []string{"", "/a", "/b", "/c"} {
		w := get(t, h, prefix+"/")
		if w.Code != http.StatusOK {
			t.Errorf("GET %s/ = %d, want %d", prefix, w.Code, http.StatusOK)
			continue
		}
		if want := `"urlBase":"` + prefix + `"`; !strings.Contains(w.Body.String(), want) {
			t.Errorf("GET %s/ did not report %s", prefix, want)
		}
	}
}

func TestUnknownPathIsNotFound(t *testing.T) {
	h := newTestHandler(t)

	if w := get(t, h, "/no/such/route"); w.Code != http.StatusNotFound {
		t.Errorf("GET /no/such/route = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// Proves the whole config path: pflag -> viper BindPFlag -> UnmarshalExact -> component SetConfig.
// A viper upgrade that broke any link would show up here, which a compile check cannot catch --
// UnmarshalExact panics at runtime rather than failing to build.
func TestFlagsReachComponentConfig(t *testing.T) {
	h := newTestHandler(t, "--liveness-fail-next=2")

	want := []int{http.StatusInternalServerError, http.StatusInternalServerError, http.StatusOK}
	for i, code := range want {
		if got := get(t, h, "/healthy").Code; got != code {
			t.Errorf("GET /healthy call %d = %d, want %d", i+1, got, code)
		}
	}
}

func TestDefaultConfigLeavesProbesPassing(t *testing.T) {
	h := newTestHandler(t)

	for i := range 3 {
		if got := get(t, h, "/healthy").Code; got != http.StatusOK {
			t.Errorf("GET /healthy call %d = %d, want %d", i+1, got, http.StatusOK)
		}
	}
}
