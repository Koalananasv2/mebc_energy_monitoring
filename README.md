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
go to http://localhost:3000/
admin:admin
![image](https://github.com/Koalananasv2/mebc_energy_monitoring/assets/152738791/533ddd72-4efc-4dd3-b2d8-ad4746ec511b)
![image](https://github.com/Koalananasv2/mebc_energy_monitoring/assets/152738791/49878eab-6934-49a9-8bf7-cc4972fa45ad)

### Install GOLANG
https://go.dev/doc/install
download the good version on https://go.dev/dl/
```sh
curl -O -L https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
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
```sh
sudo nano /etc/systemd/system/REST-MEBCV2.service
 ```

adapt the "User" and the "WorkingDirectory" and the "INFLUXDB_TOKEN"
```ini
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
enable and check status
```sh
sudo systemctl enable REST-MEBCV2.service
sudo systemctl start REST-MEBCV2.service
sudo journalctl -xefu REST-MEBCV2.service #to monitore requests
```

### Push a monitoring command using
Using CURL
```sh
curl --location 'http://localhost:30001/monitoringdata/' \
--header 'Content-Type: application/json' \
--data '{
  "generic1_temp1": 25.5,
  "generic1_temp2": 30.2,
  "generic1_temp3": 25.5,
  "generic1_temp4": 30.2,
  "generic1_temp5": 60,
  "generic1_voltage": 56,
  "generic1_current": 200.23,
  "generic1_lat": 48.8566,
  "generic1_lon": 2.3522,
  "team": "gxv3AWTxzmJEHUCaQz5AW3wyhAWsUQ5X"
}'
```

Using Python3
```python
import requests
URL = "http://localhost:30001/monitoringdata/"
PARAMS  = {
  "generic1_temp1": 25.5,
  "generic1_temp2": 30.2,
  "generic1_temp3": 25.5,
  "generic1_temp4": 30.2,
  "generic1_temp5": 60,
  "generic1_voltage": 56,
  "generic1_current": 200.23,
  "generic1_lat": 48.8566,
  "generic1_lon": 2.3522,
  "team": "gxv3AWTxzmJEHUCaQz5AW3wyhAWsUQ5X"
}
r = requests.post(url = URL, json = PARAMS)
```

### Create a visualization
Create a Dashboard
Create a Visualisation
Set the influx command as follow for the temperature1 info:
```
from(bucket: "mybucket")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "monitoring_data")
  |> filter(fn: (r) => r["_field"] == "generic1_temp1")
  |> aggregateWindow(every: v.windowPeriod, fn: mean, createEmpty: false)
  |> yield(name: "mean")
```

## Dashboard
![image](https://github.com/Koalananasv2/mebc_energy_monitoring/assets/152738791/9a236483-e518-4a14-beed-1ce56b7d07fb)


## API architecture
The archicture of the monitoring service is developed as follow:
![APIarchicture](https://github.com/Koalananasv2/mebc_energy_monitoring/blob/master/architecture.png?raw=true)
