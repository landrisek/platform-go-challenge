#!/bin/bash

# Define the JSON file path
jsonFile="artifacts/asset/update.json"

# Read the JSON content from the file
jsonContent=$(cat "$jsonFile")

# Make the cURL request with the JSON payload
curl -v -X PUT -H "Authorization: Bearer XXX" -H "Content-Type: application/json" --data-binary "$jsonContent" http://localhost:8080/update
