#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== Status ==="
curl -s -X GET "${BASE_URL}/api/status"
echo -e "\n"

echo "=== Filter: Medical ==="
curl -s -X POST "${BASE_URL}/api/filter" \
  -H "Content-Type: text/plain" \
  -d "Patient Margaret Collins, born on 04/12/1978, with SSN 523-88-4021 was admitted to the ER at St. Luke's Medical Center. Her primary care physician, Dr. Howard Banks, can be reached at hbanks@stlukesmed.org or (555) 342-9187."
echo -e "\n"

echo "=== Filter: Legal ==="
curl -s -X POST "${BASE_URL}/api/filter" \
  -H "Content-Type: text/plain" \
  -d "This agreement is entered into between Robert T. Harmon (SSN: 412-67-9034) of 1842 Birchwood Drive, Austin, TX 78701, and Meridian Law Group. Mr. Harmon can be contacted at robert.harmon@legalmail.com or by phone at (512) 778-4490."
echo -e "\n"

echo "=== Filter: Financial ==="
curl -s -X POST "${BASE_URL}/api/filter" \
  -H "Content-Type: text/plain" \
  -d "Account holder Sandra M. Patel, SSN 318-44-7762, has a checking account ending in 6204 at First National Bank. Monthly statements are sent to spatelpersonal@financemail.net. Her adviser, Michael Torres, can be reached at mitorres@firstnational.com."
echo -e "\n"

echo "=== Filter: Default ==="
curl -s -X POST "${BASE_URL}/api/filter" \
  -H "Content-Type: text/plain" \
  -d "This request does not match any preset."
echo -e "\n"

echo "=== Explain: Medical ==="
curl -s -X POST "${BASE_URL}/api/explain" \
  -H "Content-Type: text/plain" \
  -d "Patient Margaret Collins, born on 04/12/1978, with SSN 523-88-4021 was admitted to the ER at St. Luke's Medical Center. Her primary care physician, Dr. Howard Banks, can be reached at hbanks@stlukesmed.org or (555) 342-9187."
echo -e "\n"

echo "=== Explain: Legal ==="
curl -s -X POST "${BASE_URL}/api/explain" \
  -H "Content-Type: text/plain" \
  -d "This agreement is entered into between Robert T. Harmon (SSN: 412-67-9034) of 1842 Birchwood Drive, Austin, TX 78701, and Meridian Law Group. Mr. Harmon can be contacted at robert.harmon@legalmail.com or by phone at (512) 778-4490."
echo -e "\n"

echo "=== Explain: Financial ==="
curl -s -X POST "${BASE_URL}/api/explain" \
  -H "Content-Type: text/plain" \
  -d "Account holder Sandra M. Patel, SSN 318-44-7762, has a checking account ending in 6204 at First National Bank. Monthly statements are sent to spatelpersonal@financemail.net. Her adviser, Michael Torres, can be reached at mitorres@firstnational.com."
echo -e "\n"

echo "=== Explain: Default ==="
curl -s -X POST "${BASE_URL}/api/explain" \
  -H "Content-Type: text/plain" \
  -d "This request does not match any preset."
echo -e "\n"
