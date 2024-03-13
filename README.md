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

## Installation (on ubuntu)
### Install influxDB
https://docs.influxdata.com/influxdb/v2/install/
```sh
sudo apt install curl
curl -O https://dl.influxdata.com/influxdb/releases/influxdb2_2.7.5-1_amd64.deb
sudo dpkg -i influxdb2_2.7.5-1_amd64.deb
sudo rm influxdb2_2.7.5-1_amd64.deb
sudo service influxdb start
```
go to https://localhost:8086 (open firewall if needed)
![image](https://github.com/Koalananasv2/mebc_energy_monitoring/assets/152738791/594abd1e-8fa3-44b1-948b-e3cb999441f3)
Save the API token

### Install Grafana
https://grafana.com/docs/grafana/latest/setup-grafana/installation/debian/
```sh
sudo apt-get install -y apt-transport-https software-properties-common wget
sudo mkdir -p /etc/apt/keyrings/
wget -q -O - https://apt.grafana.com/gpg.key | gpg --dearmor | sudo tee /etc/apt/keyrings/grafana.gpg > /dev/null
echo "deb [signed-by=/etc/apt/keyrings/grafana.gpg] https://apt.grafana.com stable main" | sudo tee -a /etc/apt/sources.list.d/grafana.list
sudo apt-get update
sudo apt-get install grafana
sudo systemctl daemon-reload
sudo systemctl start grafana-server
sudo systemctl enable grafana-server.service
```
admin:admin
![image](https://github.com/Koalananasv2/mebc_energy_monitoring/assets/152738791/533ddd72-4efc-4dd3-b2d8-ad4746ec511b)
![image](https://github.com/Koalananasv2/mebc_energy_monitoring/assets/152738791/49878eab-6934-49a9-8bf7-cc4972fa45ad)

### Install GOLANG
https://go.dev/doc/install
download the good version on https://go.dev/dl/
```sh
cd [download directory]
sudo tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
rm go1.22.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Install GIT
```sh
sudo apt get git
```

### Install the monitoring API
```sh
cd 
git clone https://github.com/Koalananasv2/mebc_energy_monitoring.git
cd mebc_energy_monitoring/
```

### Install service to run the api on startup

Put the following service in the new file /etc/systemd/system/REST-MEBCV2.service
adapt the "User" and the "WorkingDirectory" and the "INFLUXDB_TOKEN"
```sh
[Unit]
Description=REST API for MEBC monitoring
#After=influxd.service

[Service]
User=me
WorkingDirectory=/home/me/mebc_energy_monitoring/

Environment=INFLUXDB_TOKEN=OJp9NRWe8NrXQ7Y8YznO70OKhAc5m4hojEV4ygYqKZWyfSSirQyrRhuU55pmKRQX51LjwFBeiChOXUg-HWjbjA==
Environment=INFLUXDB_PORT=8086
Environment=INFLUXDB_ORG=myorg
Environment=INFLUXDB_BCK=mybucket
Environment=PATH=$PATH:/usr/local/go/bin

ExecStart=/bin/bash -c "go run ."
# optional items below
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```




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
