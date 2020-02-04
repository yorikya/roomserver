# roomserver
Smart Home Project, based on Arduino/ESP8622 devices (clients)

This project have a few components for Smart Home, based those components you an build smart home with cheap Arduino/ESP8622 
devices.

## Content
* Client - Arduino/ESP8622 device which comunicate with server and introduce a single room
* Server - Server writen in Golang, accept room clients view and control them, rendering html for clients. 

## Instalation
### Golang
```
go get github.com/yorikya/roomserve
```
### Git 
```
git clone https://github.com/yorikya/roomserve.git
```

## Usage
### Server
```
go run .
```

### Client
1. Replace `ssID` and `wifiPass` with your own local network in [file](ardiono/room1/esp8622main.ino)
```
const String ssID     = "YuriIotLocal";         
const String wifiPass = "12345678";
```
2. Upload arduino code to an Arduino device with WiFi module or just use ESP8622 bord with GPIO pins
3. Each client runs local server on port 80:
```
curl http://10.0.0.9/logs 
```
Endpoints:
* `/logs` return last 50 internal logs
* `/action` internal use endpoint for ardiono commands
* `/data` return sensors values
 
## Features
### Server
* Accept rooms/devices via configuration [file](config/rooms.json)
* Web server rendering `room` view
* Web socket using for updating web view
* Comunicate with clients to perform actions for diffrent devices 
* Serving static files 

### Client
* Send DHT Sensor data
* Control Door lock
* Air Conditioner remote control
* Control Room Light
* Control RGB Strip 



