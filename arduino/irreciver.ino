#include <ArduinoJson.h>
#include "EspMQTTClient.h"

//Client Data
const char * roomID = "office";
const char * clientID = "officeClient";
//Topic names
const char * inTopic = "officeInTopic";
const char * outTopic = "officeOutTopic";
const char * updateTopic = "officeUpdateTopic";
const char * updateTopicIn = "officeUpdateTopicIn";

//Network data
const char * ssid = "Danielle_2.4";
const char * wifipass = "0524713014";
const char * serverip = "10.0.0.2";

//Sensors
char * movementSen = "78.8";
char * tempSen = "26.6";
char * airCond = "on";
char * lightMain = "on";
char * lightSec = "off";

StaticJsonDocument<1024> client_data;

//Process Sheduling
unsigned long previousMillis = 0;
const long clientDataInterval = 60000;

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
  client_data["updated"] = false;
}


void onConnectionEstablished() {
  client.subscribe(updateTopicIn, [] (const String & payload)  {
    DynamicJsonDocument doc(1024);
    deserializeJson(doc, payload);

    const char * action = doc["action"];
    const char * deviceid  = doc["deviceid"];
    if (strcmp(action, "update") == 0) {
      client_data["updated"] = true;
      Serial.println("data was updated!!!!!");
    }
    Serial.print(updateTopicIn);
    Serial.print(" get action ");
    Serial.print(action);
     Serial.print(" device-id ");
    Serial.print(deviceid);
    Serial.println(payload);
  });

   client.subscribe(inTopic, [] (const String & payload)  {
    DynamicJsonDocument doc(1024);
    deserializeJson(doc, payload);

    const char * action = doc["action"];
    const char * deviceid  = doc["deviceid"];
    Serial.print(inTopic);
    Serial.print(" get action ");
    Serial.print(action);
     Serial.print(" device-id ");
    Serial.print(deviceid);
    Serial.println(payload);
  });

}

    
void loop() {
  client.loop();
  unsigned long currentMillis = millis();

  if (client_data["updated"] == false && currentMillis - previousMillis >= 1000) {
    previousMillis = currentMillis;
    sendRequestData();
  }

  if (currentMillis - previousMillis >= clientDataInterval) {
    previousMillis = currentMillis;
    sendClientData();
  }

  
}

void sendClientData() {
  StaticJsonDocument<1024> msg;
  msg["movementSen"] = movementSen;
  msg["tempSen"] = tempSen;
  msg["airCond"] = airCond;
  msg["lightMain"] = lightMain;
  msg["lightSec"] = lightSec;
  msg["action"] = "update";
  char buffer[1024];
  serializeJson(msg, buffer);
  client.publish(outTopic, buffer);
  Serial.print("send client update data");
  Serial.println(buffer);
}

void sendRequestData() {
    StaticJsonDocument<1024> msg;
    msg["action"] = "update";
    msg["in"] =  random(0, 10);
    char buffer[1024];
    serializeJson(msg, buffer);
    client.publish(updateTopic, buffer);
    Serial.print("send client request data");
    Serial.println(buffer);  
}