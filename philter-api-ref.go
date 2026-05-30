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
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Status struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type Span struct {
	Id             string  `json:"id"`
	CharacterStart int     `json:"characterStart"`
	CharacterEnd   int     `json:"characterEnd"`
	FilterType     string  `json:"filterType"`
	Context        string  `json:"context"`
	DocumentId     string  `json:"documentId"`
	Confidence     float64 `json:"confidence"`
	Text           string  `json:"text"`
	Replacement    string  `json:"replacement"`
	Ignored        bool    `json:"ignored"`
}

type Explanation struct {
	AppliedSpans []Span `json:"appliedSpans"`
	IgnoredSpans []Span `json:"ignoredSpans"`
}

type Explain struct {
	FilteredText string      `json:"filteredText"`
	Context      string      `json:"context"`
	DocumentId   string      `json:"documentId"`
	Explanation  Explanation `json:"explanation"`
}

func main() {

	fmt.Println("Starting service...")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/filter", filter).Methods("POST")
	router.HandleFunc("/api/explain", explain).Methods("POST")
	router.HandleFunc("/api/status", status).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))

}

func filter(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("x-document-id", "asdfghjkl12345678")

	fmt.Fprintln(w, "{{{REDACTED-entity}}} was a patient.")

}

func explain(w http.ResponseWriter, r *http.Request) {

	// The Philter API passes policy, context, and document id as query
	// parameters: p (policy), c (context), d (document id).
	query := r.URL.Query()

	context := query.Get("c")
	if context == "" {
		context = "none"
	}

	documentId := query.Get("d")
	if documentId == "" {
		documentId = generateDocumentId()
	}

	// The policy is accepted for parity with the real API but does not change
	// this reference response.
	_ = query.Get("p")

	filteredText := "{{{REDACTED-entity}}} was a patient and his ssn was {{{REDACTED-ssn}}}."

	appliedSpans := []Span{
		{
			Id:             generateDocumentId(),
			CharacterStart: 0,
			CharacterEnd:   17,
			FilterType:     "NER_ENTITY",
			Context:        context,
			DocumentId:     documentId,
			Confidence:     0.918,
			Text:           "George Washington",
			Replacement:    "{{{REDACTED-entity}}}",
			Ignored:        false,
		},
		{
			Id:             generateDocumentId(),
			CharacterStart: 48,
			CharacterEnd:   59,
			FilterType:     "SSN",
			Context:        context,
			DocumentId:     documentId,
			Confidence:     1,
			Text:           "123-45-6789",
			Replacement:    "{{{REDACTED-ssn}}}",
			Ignored:        false,
		},
	}

	explain := Explain{
		FilteredText: filteredText,
		Context:      context,
		DocumentId:   documentId,
		Explanation: Explanation{
			AppliedSpans: appliedSpans,
			IgnoredSpans: []Span{},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-document-id", documentId)

	if err := json.NewEncoder(w).Encode(explain); err != nil {
		panic(err)
	}

}

// generateDocumentId returns a random hex identifier, mirroring the unique
// document id Philter assigns to each request.
func generateDocumentId() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(b)
}

func status(w http.ResponseWriter, r *http.Request) {

	status := Status{Status: "Healthy", Version: "1.0.0"}

	if err := json.NewEncoder(w).Encode(status); err != nil {
		panic(err)
	}

}
