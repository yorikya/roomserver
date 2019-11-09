#include "EspMQTTClient.h"
#include <ArduinoJson.h>

const char * roomID = "office";
const char * clientID = "officeClient";
const char * inTopic = "officeInTopic";
const char * outTopic = "officeOutTopic";

const char * ssid = "Danielle_2.4";
const char * wifipass = "0524713014";
const char * serverip = "10.0.0.2"; 

char * movementSen = "78.8";
char * tempSen = "26.6";
char * airCond = "on";
char * lightMain = "on";
char * lightSec = "off";

unsigned long previousMillis = 0;    
// constants won't change:
const long clientDataInterval = 5000;   

EspMQTTClient client(
  ssid,
  wifipass,
  serverip,  // MQTT Broker server ip
  "",   // Can be omitted if not needed
  "",   // Can be omitted if not needed
  clientID      // Client name that uniquely identify your device
);

void setup() {
  Serial.begin(115200);
}

void onConnectionEstablished() {
  client.subscribe(inTopic, [] (const String &payload)  {
    DynamicJsonDocument doc(1024);
    deserializeJson(doc, payload);

    const char * b = doc["msg"];
    Serial.println(b);
  });
}

void loop() {
  client.loop();

  unsigned long currentMillis = millis();

  if (currentMillis - previousMillis >= clientDataInterval) {
    previousMillis = currentMillis;
    
    sendClientData();
  } 
}

void sendClientData() {
   StaticJsonDocument<1024> data;
   data["movementSen"] = movementSen;
   data["tempSen"] = tempSen;
   data["airCond"] = airCond;
   data["lightMain"] = lightMain;
   data["lightSec"] = lightSec;
   char buffer[1024];
   serializeJson(data, buffer);
   client.publish(outTopic, buffer);
   Serial.println(buffer);
}
