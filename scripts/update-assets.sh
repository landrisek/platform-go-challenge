#!/bin/bash

curl -v -X PUT -H "Authorization: Bearer XXX" -H "Content-Type: application/json" -d '[
  {
    "id": 1,
    "charts": [
      {
        "id": 1,
        "title": "Chart 1B",
        "description": "test",
        "data": ""
      }
    ],
    "insights": [
      {
        "title": "Insight 1B"
      }
    ],
    "audiences": [
      {
        "title": "Audience 1B"
      }
    ]
  }
]' http://localhost:8080/update
