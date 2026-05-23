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
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Status struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type Explain struct {
	FilteredText string `json:"filteredText"`
	Context      string `json:"context"`
	DocumentId   string `json:"documentId"`
}

const (
	medicalRequest = `Patient Margaret Collins, born on 04/12/1978, with SSN 523-88-4021 was admitted to the ER at St. Luke's Medical Center. Her primary care physician, Dr. Howard Banks, can be reached at hbanks@stlukesmed.org or (555) 342-9187.`

	medicalResponse = `Patient {{{REDACTED-person-name}}}, born on 04/12/1978, with SSN {{{REDACTED-ssn}}} was admitted to the ER at St. Luke's Medical Center. Her primary care physician, Dr. {{{REDACTED-person-name}}}, can be reached at {{{REDACTED-email-address}}} or {{{REDACTED-phone-number}}}.`

	legalRequest = `This agreement is entered into between Robert T. Harmon (SSN: 412-67-9034) of 1842 Birchwood Drive, Austin, TX 78701, and Meridian Law Group. Mr. Harmon can be contacted at robert.harmon@legalmail.com or by phone at (512) 778-4490.`

	legalResponse = `This agreement is entered into between {{{REDACTED-person-name}}} (SSN: {{{REDACTED-ssn}}}) of {{{REDACTED-street-address}}}, and Meridian Law Group. Mr. {{{REDACTED-person-name}}} can be contacted at {{{REDACTED-email-address}}} or by phone at {{{REDACTED-phone-number}}}.`

	financialRequest = `Account holder Sandra M. Patel, SSN 318-44-7762, has a checking account ending in 6204 at First National Bank. Monthly statements are sent to spatelpersonal@financemail.net. Her adviser, Michael Torres, can be reached at mitorres@firstnational.com.`

	financialResponse = `Account holder {{{REDACTED-person-name}}}, SSN {{{REDACTED-ssn}}}, has a checking account ending in 6204 at First National Bank. Monthly statements are sent to {{{REDACTED-email-address}}}. Her adviser, {{{REDACTED-person-name}}}, can be reached at {{{REDACTED-email-address}}}.`

	defaultResponse = `{{{REDACTED-entity}}} was a patient.`
)

func main() {

	fmt.Println("Starting service...")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/filter", filter).Methods("POST")
	router.HandleFunc("/api/explain", explain).Methods("POST")
	router.HandleFunc("/api/status", status).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))

}

func getFilteredText(body string) string {
	switch strings.TrimSpace(body) {
	case medicalRequest:
		return medicalResponse
	case legalRequest:
		return legalResponse
	case financialRequest:
		return financialResponse
	default:
		return defaultResponse
	}
}

func filter(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	w.Header().Add("x-document-id", "asdfghjkl12345678")

	fmt.Fprintln(w, getFilteredText(string(body)))

}

func explain(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	context := params["context"]
	documentId := params["documentId"]

	explain := Explain{
		FilteredText: getFilteredText(string(body)),
		Context:      context,
		DocumentId:   documentId,
	}

	if err := json.NewEncoder(w).Encode(explain); err != nil {
		panic(err)
	}

}

func status(w http.ResponseWriter, r *http.Request) {

	status := Status{Status: "Healthy", Version: "1.0.0"}

	if err := json.NewEncoder(w).Encode(status); err != nil {
		panic(err)
	}

}
