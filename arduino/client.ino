#include "EspMQTTClient.h"

//Client Data
const char * roomID = "room1";
const char * clientID = "room1Client";
//Topic names
const char * inTopic = "room1InTopic";
const char * outTopic = "room1OutTopic";

//Network data
const char * ssid = "Danielle_2.4";
const char * wifipass = "0524713014";
const char * serverip = "10.0.0.2";

//Movement Sensor
int movesensor = 5;              // the pin that the sensor is atteched to


//Libraries
#include <DHT.h>;

//Constants
#define DHTPIN 4     // what pin we're connected to
#define DHTTYPE DHT22   // DHT 22  (AM2302)
DHT dht(DHTPIN, DHTTYPE, 11); //// Initialize DHT sensor for normal 16mhz Arduino

//Process Sheduling
unsigned long previousMillis = 0;
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
  dht.begin();
  pinMode(movesensor, INPUT);    // initialize sensor as an input
  Serial.begin(115200);
}


void onConnectionEstablished() {
  client.subscribe(inTopic, [] (const String & payload)  {
    
    Serial.print("get data in-topic: ");
    Serial.println(payload);
  });


}

    
void loop() {
  client.loop();
  unsigned long currentMillis = millis();


  if (currentMillis - previousMillis >= clientDataInterval) {
    previousMillis = currentMillis;
    String deviceid = "/dht/"; 

    client.publish(outTopic, deviceid + "Humidity/" + String(dht.readHumidity()).c_str());
    client.publish(outTopic, deviceid + "Temperature/" + String(dht.readTemperature()).c_str());

    client.publish(outTopic, "/movesensor/state/" + String(digitalRead(movesensor)));
    Serial.println(digitalRead(movesensor));
    
    Serial.println("send client dht sensor data, move senso");
  }
  
}
