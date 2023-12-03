# MEBC Energy Class Monitoring API

In this page you will find the information about the API used to monitor the boat of the Energy Boat Challenge 2024.
Each team will connect to the API with a unique token and will provide information about the boat at a frequency between 1Hz and 2Hz.
A link to the dashboard will also be provided to all teams.

## Features

The following information are mandatory:
- Temperature of all batteries
- Current & voltage of primary batteries (if there are two or more primary batteries, provide the downstream value if the electrical architecure allows it)

The following information are nice to have:
- Lattitude & longitude
- Voltage/current or the power at the motor (not yet implemented).
Other information can be added the monitoring system uppon team's request.

"Mandatory" and "Nice to have" information will be share in live during the competition.
But teams can also request to push confidential metrics available only on their own dashboard.

## Interface

The interface with the API is a simple REST API with only the POST method exposed.

Using Curl command:
```sh
curl --location 'http://[IP]:[PORT]/monitoringdata/' \
--header 'Content-Type: application/json' \
--data '{
    "temp1": 39.55,
    "temp2": 34.53,
    "temp3": 34.35,
    "voltage": 54.78,
    "current": 167.16,
    "lat": 7.4423319,
    "lon": 43.7412106,
    "team": "pC9rVUr9F3WV7TX3qF584hUcuzh2WXxA"
}'
```

Using Python3
```python
import requests
URL = "http://[IP]:[PORT]/monitoringdata/"
PARAMS  = {
    "temp1": 39.55,
    "temp2": 34.53,
    "temp3": 34.35,
    "temp4": 54.78,
    "temp5": 60.02,
    "voltage": 54.78,
    "current": 167.16,
    "lat": 7.4423319,
    "lon": 43.7412106,
    "team": "pC9rVUr9F3WV7TX3qF584hUcuzh2WXxA"
}
r = requests.post(url = URL, json = PARAMS)
```

## Dashboard

Dashboard is rendered with Grafana. Credential will be provided to each teams.
_Data below are random_
![Dashboard](https://s5.gifyu.com/images/SiSOF.gif)

## API architecture

The archicture of the monitoring service is developed as follow:
![APIarchicture](https://github.com/Koalananasv2/mebc_energy_monitoring/blob/master/architecture.jpg?raw=true)
