/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExplainEchoesParamsAndReturnsExplanation(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/explain?p=hipaa&c=ctx1&d=doc-123",
		strings.NewReader("George Washington was a patient and his ssn was 123-45-6789."))
	rr := httptest.NewRecorder()

	explain(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	if id := rr.Header().Get("x-document-id"); id != "doc-123" {
		t.Errorf("x-document-id header = %q, want doc-123", id)
	}

	var e Explain
	if err := json.Unmarshal(rr.Body.Bytes(), &e); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if e.Context != "ctx1" {
		t.Errorf("context = %q, want ctx1", e.Context)
	}
	if e.DocumentId != "doc-123" {
		t.Errorf("documentId = %q, want doc-123", e.DocumentId)
	}
	if e.FilteredText == "" {
		t.Error("filteredText is empty")
	}
	if len(e.Explanation.AppliedSpans) == 0 {
		t.Fatal("appliedSpans is empty; the explanation must describe the redactions")
	}
	for i, span := range e.Explanation.AppliedSpans {
		if span.Id == "" {
			t.Errorf("appliedSpans[%d].id is empty", i)
		}
		if span.FilterType == "" {
			t.Errorf("appliedSpans[%d].filterType is empty", i)
		}
		if span.Replacement == "" {
			t.Errorf("appliedSpans[%d].replacement is empty", i)
		}
		// Context and documentId must be echoed into each span.
		if span.Context != "ctx1" {
			t.Errorf("appliedSpans[%d].context = %q, want ctx1", i, span.Context)
		}
		if span.DocumentId != "doc-123" {
			t.Errorf("appliedSpans[%d].documentId = %q, want doc-123", i, span.DocumentId)
		}
	}
}

func TestExplainDefaultsAndGeneratedDocumentId(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/explain", strings.NewReader("text"))
	rr := httptest.NewRecorder()

	explain(rr, req)

	raw := rr.Body.String()

	var e Explain
	if err := json.Unmarshal([]byte(raw), &e); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if e.Context != "none" {
		t.Errorf("default context = %q, want none", e.Context)
	}
	if e.DocumentId == "" {
		t.Error("expected a generated documentId when 'd' is omitted")
	}
	if hdr := rr.Header().Get("x-document-id"); hdr != e.DocumentId {
		t.Errorf("x-document-id header %q does not match body documentId %q", hdr, e.DocumentId)
	}
	if len(e.Explanation.AppliedSpans) > 0 && e.Explanation.AppliedSpans[0].DocumentId != e.DocumentId {
		t.Errorf("span documentId = %q, want generated %q",
			e.Explanation.AppliedSpans[0].DocumentId, e.DocumentId)
	}
	// Empty ignoredSpans must serialize as [] (not null) so clients can iterate it.
	if !strings.Contains(raw, `"ignoredSpans":[]`) {
		t.Errorf("ignoredSpans should serialize as [], got: %s", raw)
	}
}

func TestExplainGeneratesDistinctDocumentIds(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/explain", strings.NewReader("text"))
		rr := httptest.NewRecorder()
		explain(rr, req)
		var e Explain
		if err := json.Unmarshal(rr.Body.Bytes(), &e); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if ids[e.DocumentId] {
			t.Fatalf("generated documentId %q was not unique", e.DocumentId)
		}
		ids[e.DocumentId] = true
	}
}

func TestStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rr := httptest.NewRecorder()

	status(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	var s Status
	if err := json.Unmarshal(rr.Body.Bytes(), &s); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if s.Status == "" || s.Version == "" {
		t.Errorf("status response missing fields: %+v", s)
	}
}

func TestFilter(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/filter",
		strings.NewReader("George Washington was a patient."))
	rr := httptest.NewRecorder()

	filter(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if rr.Header().Get("x-document-id") == "" {
		t.Error("filter response is missing the x-document-id header")
	}
	if !strings.Contains(rr.Body.String(), "REDACTED") {
		t.Errorf("filter response should contain a redaction marker, got: %s", rr.Body.String())
	}
}
