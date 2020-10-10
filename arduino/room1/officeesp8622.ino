#include <DHT.h>
#include <ESP8266WiFi.h>
#include <PubSubClient.h>

#define DHTPIN 5     // D3
#define DHTTYPE DHT22   // DHT 22  (AM2302)
int DOOR_RELAY = 16; //D2

const char* ssid = "HomeExtender";
const char* password = "0524713014";
const char* mqtt_server = "192.168.123.100";


WiFiClient espClient;
PubSubClient client(espClient);
unsigned long lastMsg = 0;
#define MSG_BUFFER_SIZE  (50)
char msg[MSG_BUFFER_SIZE];
DHT dht(DHTPIN, DHTTYPE, 11); //// Initialize DHT sensor for normal 16mhz Arduino

void setup_wifi() {

  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  randomSeed(micros());

  Serial.println("");
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
}

void callback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Message arrived [");
  Serial.print(topic);
  Serial.print("] ");
  for (int i = 0; i < length; i++) {
    Serial.print((char)payload[i]);
  }
  Serial.println();

  // Switch on the LED if an 1 was received as first character
  if ((char)payload[0] == '1') {
    digitalWrite(DOOR_RELAY, HIGH);
    Serial.println("open door");
    delay(1500);
    digitalWrite(DOOR_RELAY, LOW);
    Serial.println("release door");
  } else {
    Serial.print("unknown door action");
    Serial.println((char)payload[0]);
  }

}

void reconnect() {
  // Loop until we're reconnected
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Create a random client ID
    String clientId = "ESP8266Client-";
    clientId += String(random(0xffff), HEX);
    // Attempt to connect
    if (client.connect(clientId.c_str())) {
      Serial.println("connected");
      // Once connected, publish an announcement...
      client.publish("deviceMQTT", "hello world");
      // ... and resubscribe
      client.subscribe("home/office/door");
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      // Wait 5 seconds before retrying
      delay(5000);
    }
  }
}

void setup() {
  Serial.begin(115200);
  Serial.setDebugOutput(true);
  setup_wifi();
  client.setServer(mqtt_server, 1883);
  client.setCallback(callback);

  dht.begin();
  pinMode(DOOR_RELAY, OUTPUT);
}

void loop() {

  if (!client.connected()) {
    reconnect();
  }
  client.loop();

  unsigned long now = millis();
  if (now - lastMsg > 60000) {
    lastMsg = now;
    snprintf (msg, MSG_BUFFER_SIZE, "%s", String(dht.readHumidity()).c_str());
    Serial.print("Publish message: ");
    Serial.println(msg);
    client.publish("office/Humidity", msg);
    snprintf (msg, MSG_BUFFER_SIZE, "%s", String(dht.readTemperature()).c_str());
    Serial.print("Publish message Temperature: ");
    Serial.println(msg);
    client.publish("office/Temperature", msg);
  }

  
}