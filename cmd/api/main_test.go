package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleNetworkConfig_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/network", nil)
	rr := httptest.NewRecorder()

	handleNetworkConfig(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("expected json content-type, got %q", ct)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["status"] != "not implemented yet" {
		t.Fatalf("unexpected status: %q", body["status"])
	}
}

func TestHandleNetworkConfig_POST(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/network", nil)
	rr := httptest.NewRecorder()

	handleNetworkConfig(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"status": "received"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestHandleNetworkConfig_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/v1/network", nil)
	rr := httptest.NewRecorder()

	handleNetworkConfig(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}
