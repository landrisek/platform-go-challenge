#!/bin/bash

curl -X POST -H "Authorization: Bearer XXX" -H "Content-Type: application/json" -d '{
  "charts": [
    {
      "title": "Chart 1",
      "user_id": 1,
      "axes_titles": "X-Axis, Y-Axis",
      "data": "1,2,3,4,5"
    },
    {
      "title": "Chart 2",
      "user_id": 1,
      "axes_titles": "X-Axis, Y-Axis",
      "data": "5,4,3,2,1"
    }
  ],
  "insights": [
    {
      "title": "Insight 1",
      "user_id": 2,
      "text": "This is Insight 1"
    },
    {
      "title": "Insight 2",
      "user_id": 2,
      "text": "This is Insight 2"
    }
  ],
  "audiences": [
    {
      "title": "Audience 1",
      "user_id": 3,
      "characteristics": "Age: 25-35, Gender: Male",
      "description": "This is Audience 1"
    },
    {
      "title": "Audience 2",
      "user_id": 3,
      "characteristics": "Age: 18-24, Gender: Female",
      "description": "This is Audience 2"
    }
  ]
}' http://localhost:8080/create

