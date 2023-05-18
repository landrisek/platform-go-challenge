#!/bin/bash

curl -X POST -H "Authorization: Bearer XXX" -H "Content-Type: application/json" -d '{
  "users": {
    "1": {
      "charts": [
        {
          "title": "Chart 1",
          "axes_titles": "X-Axis, Y-Axis",
          "data": "1,2,3,4,5"
        },
        {
          "title": "Chart 2",
          "axes_titles": "X-Axis, Y-Axis",
          "data": "5,4,3,2,1"
        }
      ],
      "insights": [
        {
          "title": "Insight 1",
          "text": "This is Insight 1"
        },
        {
          "title": "Insight 2",
          "text": "This is Insight 2"
        }
      ],
      "audiences": [
        {
          "title": "Audience 1",
          "characteristics": "Age: 25-35, Gender: Male",
          "description": "This is Audience 1"
        },
        {
          "title": "Audience 2",
          "characteristics": "Age: 18-24, Gender: Female",
          "description": "This is Audience 2"
        }
      ]
    },
    "4": {
      "charts": [
        {
          "title": "Chart 3",
          "axes_titles": "X-Axis, Y-Axis",
          "data": "3,2,1"
        }
      ],
      "insights": [
        {
          "title": "Insight 3",
          "text": "This is Insight 3"
        }
      ],
      "audiences": [
        {
          "title": "Audience 3",
          "characteristics": "Age: 40-50, Gender: Male",
          "description": "This is Audience 3"
        }
      ]
    }
  }
}' http://localhost:8080/create
