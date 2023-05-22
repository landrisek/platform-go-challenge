#!/bin/bash

# Define the JSON file path
jsonFile="artifacts/asset/create.json"

# Read the JSON content from the file
jsonContent=$(cat "$jsonFile")

# Make the cURL request with the JSON payload
curl -v -X POST -H "Authorization: Bearer XXX" -H "Content-Type: application/json" --data-binary "$jsonContent" http://localhost:8080/create
