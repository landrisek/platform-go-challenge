#!/bin/bash

curl -v -X DELETE -H "Authorization: Bearer XXX" -H "Content-Type: application/json" -d '[
  {
    "id": 1,
    "charts": [
      {
        "id": 1
      }
    ],
    "insights": [
      {
        "id": 1
      }
    ],
    "audiences": [
      {
        "id": 1
      }
    ]
  }
]' http://localhost:8080/delete
