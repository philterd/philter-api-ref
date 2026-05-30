#!/bin/bash -e

docker build -t philterd/philter-api-ref .
docker push philterd/philter-api-ref
