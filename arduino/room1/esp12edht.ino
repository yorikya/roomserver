#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ArduinoJson.h>
#include <ESP8266HTTPClient.h>
#include <ESP8266WebServer.h>

//Client Data
const String roomID = "room1";
const String clientID = String(roomID+ "_dht");

const String ssID     = "Danielle_2.4";         
const String wifiPass = "0524713014";    

const String serverIP = "10.0.0.4"; 
const String clientInfo = "--- Client Info: --- \nClientID: " + clientID + ", WiFi SSID: " + ssID + ", Main Server IP: " + serverIP;

//Process Sheduling
unsigned long previousMillis = 0;
const long clientDataInterval = 5000;

//Client data
StaticJsonDocument<512> clientData;

WiFiClient client;

HTTPClient http;

ESP8266WebServer server(80); //Server on port 80

#include <DHT.h>
#define DHTPIN 4     // what pin we're connected to
#define DHTTYPE DHT22   // DHT 22  (AM2302)

DHT dht(DHTPIN, DHTTYPE, 11); //// Initialize DHT sensor for normal 16mhz Arduino

// --- Logger 
#define LOGSSIZE 50

struct LoggerBuffer {
  int maxsize;
  int capacity;
  String logs[LOGSSIZE];
};

//Buffer logs 
LoggerBuffer logger = {LOGSSIZE, 0, {}};

void logPrintln(String logstr) {
  char tspref[16];
  sprintf(tspref,"[%u ms] ", millis());

  if (logger.capacity >= logger.maxsize) {
    logger.capacity = 0;
  } 
 
  logger.logs[logger.capacity] = tspref + logstr; 
  Serial.println(logger.logs[logger.capacity]);  
  logger.capacity++;
}

void handleLogs() {
   String str = String(clientInfo + "\n\n");
   str += "Logs: (latst: "+ String(logger.maxsize )+ ")\n";
   for (int i = 0; i <= logger.capacity; i++) {
     str += logger.logs[i] + "\n";
   } 
   server.send(200, "text/plain", str);
}
// --- End Logger

void setup() 
{
  dht.begin();
  Serial.begin(115200);
  Serial.setDebugOutput(true);

  logPrintln("Configuring access point...");
  WiFi.mode(WIFI_STA);  
  WiFi.begin(ssID, wifiPass);             // Connect to the network
  logPrintln("Connecting to " + ssID);

  int i = 0;
  while (WiFi.localIP().toString() == "0.0.0.0") { // Wait for the Wi-Fi to connect
    delay(1000);
    logPrintln("connecting to  SSID:" + ssID + " attempt: " + String(++i) + ", WiFi Status: " + String(WiFi.status()));
  }
  logPrintln("Connection established!\nIP address: " + WiFi.localIP().toString());  

  server.on("/logs", handleLogs);      //Which routine to handle at root location
  
  server.begin();     
  logPrintln(clientInfo); //Start server
  authRequest();
}

void authRequest() {
  String url = "http://" + serverIP + "/auth?clientid=" + clientID;
  String resp;
  int i;
  while (resp == "") {
    logPrintln("try auth attempt: " + String(++i));
    resp = getRequest(url);
    if (resp != "") {
      DeserializationError error = deserializeJson(clientData, resp);
      // Test if parsing succeeds.
      if (error) {
        logPrintln("auth response deserializeJson() failed:" );
        logPrintln(error.c_str());
        resp = "";
        continue;
      }
      
      if (clientData["Success"] == true) {
         logPrintln("authentication success, resp:"+ resp);
         return;
      } 
    }
    
    delay(1000);
  }
}

String humidity;
String temperature;

void loop() 
{
  server.handleClient();
  unsigned long currentMillis = millis();
  
  if (currentMillis - previousMillis >= clientDataInterval) {
    previousMillis = currentMillis;
    
    sendSensorData(clientData["Devices"]["room1_dht_Humidity"]["Name"],String(dht.readHumidity()).c_str());
    sendSensorData(clientData["Devices"]["room1_dht_Temperature"]["Name"],String(dht.readTemperature()).c_str()); 
  } 
  
}

// --- Help Functions
void sendSensorData(String device, String sensor,String value) {
  String url = "http://" + serverIP + "/update";
  url += "?device=" + device;
  url += "&value=" + value;
  url += "&clientid=" + clientID;

  getRequest(url);
}

String getRequest(String url) {
    String payload;
    if (http.begin(client, url)) {  // HTTP
      int httpCode = http.GET();
      // httpCode will be negative on error
      if (httpCode > 0) {
        if (httpCode == HTTP_CODE_OK) {
          payload = http.getString();
        }
      } else {
        logPrintln("[HTTP] GET... failed, url:" + url + " error:" + http.errorToString(httpCode).c_str());
      }
      http.end();
      logPrintln("[HTTP] end connection to url: " + url);
    } else {
      logPrintln("[HTTP] Unable to connect to url: " + url);
    }
    return payload;
}
