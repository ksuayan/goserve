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

## Setup Docker
```bash
docker network create mysql-network

docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=admin -p 3306:3306 -d mysql:latest

docker exec -it mysql-container mysql -u root -p

docker cp glucose_10-12-2024.csv mysql-container:/var/lib/mysql-files/
```


## Load Data to MySQL
```sql

CREATE DATABASE freestyle;

CREATE TABLE glucose (
    id INT PRIMARY KEY AUTO_INCREMENT,
    device VARCHAR(100),
    serial_number VARCHAR(50),
    device_timestamp DATETIME,
    record_type INT NULL,
    historic_glucose_mg_dl INT NULL,
    scan_glucose_mg_dl INT NULL,
    non_numeric_rapid_acting_insulin VARCHAR(50),
    rapid_acting_insulin_units DECIMAL(5,2) NULL,
    non_numeric_food VARCHAR(50),
    carbohydrates_grams INT NULL,
    carbohydrates_servings DECIMAL(5,2) NULL,
    non_numeric_long_acting_insulin VARCHAR(50),
    long_acting_insulin_units DECIMAL(5,2) NULL,
    notes TEXT,
    strip_glucose_mg_dl INT NULL,
    ketone_mmol_l DECIMAL(5,2) NULL,
    meal_insulin_units DECIMAL(5,2) NULL,
    correction_insulin_units DECIMAL(5,2) NULL,
    user_change_insulin_units DECIMAL(5,2) NULL
);


LOAD DATA INFILE '/var/lib/mysql-files/glucose_10-12-2024.csv'  
INTO TABLE glucose  
FIELDS TERMINATED BY ','  
OPTIONALLY ENCLOSED BY '"'  
LINES TERMINATED BY '\n'  
IGNORE 2 LINES  (device, serial_number, @device_timestamp, record_type,   
				 @historic_glucose_mg_dl, @scan_glucose_mg_dl,   
				 non_numeric_rapid_acting_insulin, @rapid_acting_insulin_units,   
				 non_numeric_food, @carbohydrates_grams, 
				 @carbohydrates_servings,   
				 non_numeric_long_acting_insulin, 
				 @long_acting_insulin_units, notes,   
				 @strip_glucose_mg_dl, @ketone_mmol_l, 
				 @meal_insulin_units, 
				 @correction_insulin_units,
				 @user_change_insulin_units)  
SET id = NULL, 
device_timestamp = STR_TO_DATE(@device_timestamp, '%m-%d-%Y %h:%i %p'),
historic_glucose_mg_dl = COALESCE(NULLIF(@historic_glucose_mg_dl, ''), 0);

SELECT device_timestamp, record_type, historic_glucose_mg_dl
FROM glucose
WHERE device_timestamp BETWEEN '2024-08-01 00:00:00' AND '2024-08-30 23:59:59'
AND record_type = 0;

```

