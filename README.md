# goserve

goserve is an application for visualizing CGM data. 

I built it as a sandbox for experimenting with GoLang and MySQL as the backend with a ReactJS application consuming data from its API endpoint.

## Tech Stack
- Go Programming Language - for the API server
- ReactJS - Single Page Application
- D3.js - visualization
- gorilla/mux - router for request handling  and dispatch of incoming request

## Frontend Usage

Set date range in the URL as follows
http://localhost:3000/#?fromDate=2024-10-01&toDate=2024-10-13

## Backend Usage

Raw Data
http://localhost:3000/api/raw/2024-10-01/2024-10-13
```json
[
    {
        "timestamp": "2024-10-01T00:02:00-07:00",
        "glucose": 172
    },
    {
        "timestamp": "2024-10-01T00:07:00-07:00",
        "glucose": 168
    },
    {
        "timestamp": "2024-10-01T00:12:00-07:00",
        "glucose": 164
    },
    ...
]
```


Summarized Glucose by Time of Day
http://localhost:3000/api/glucose/2024-10-01/2024-10-13
```json
[
    {
        "timeslot": "00:00:00",
        "readings": null,
        "pct_05": 111,
        "pct_95": 215,
        "average": 157.27272727272728
    },
    {
        "timeslot": "00:05:00",
        "readings": null,
        "pct_05": 117,
        "pct_95": 223,
        "average": 158.36363636363637
    },
    {
        "timeslot": "00:10:00",
        "readings": null,
        "pct_05": 117,
        "pct_95": 227,
        "average": 158.63636363636363
    },
    ...
]    
```