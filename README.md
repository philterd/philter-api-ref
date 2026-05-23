# Philter API Reference Implementation

[Philter](https://www.philterd.ai/philter/) identifies and manipulates sensitive information in text. This project is a reference implementation of Philter's [Filtering API](https://philterd.github.io/philter/api_and_sdks/api/filtering_api.html). This project can be used for testing Philter integrations without requiring an instance of Philter.

**This project is currently a work in progress.**

## Build and Run

**Local:**
```
make build
./philter-api-ref
```

**Docker:**
```
docker-compose up --build
```

## Sending Requests

Filter text:
```
curl -s -X POST http://localhost:8080/api/filter \
  -H "Content-Type: text/plain" \
  -d "Patient Margaret Collins, born on 04/12/1978, with SSN 523-88-4021 was admitted to the ER at St. Luke's Medical Center. Her primary care physician, Dr. Howard Banks, can be reached at hbanks@stlukesmed.org or (555) 342-9187."
```

Sample response:
```
Patient {{{REDACTED-person-name}}}, born on 04/12/1978, with SSN {{{REDACTED-ssn}}} was admitted to the ER at St. Luke's Medical Center. Her primary care physician, Dr. {{{REDACTED-person-name}}}, can be reached at {{{REDACTED-email-address}}} or {{{REDACTED-phone-number}}}.
```

Check status:
```
curl -s http://localhost:8080/api/status
```

To exercise all preset requests at once, run `./test.sh` with the service already running.

## License

This project is licensed under the Apache License, version 2.0.

Copyright 2023-2026 Philterd, LLC.
Copyright 2022 Mountain Fog, Inc.

Philter is a registered trademark of Philterd, LLC.