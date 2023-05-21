#!/bin/bash

curl -v -X POST -H "Authorization: Bearer XXX" -H "Content-Type: application/json" -d '[{"id": 1}]' http://localhost:8080/read