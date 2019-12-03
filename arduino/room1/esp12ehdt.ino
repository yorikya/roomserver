#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>

//Client Data
const String roomID = "room1";
const String clientID = String(roomID+ "_hdt");

const String ssID     = "Danielle_2.4";         
const String wifiPass = "0524713014";    

const String serverIP = "10.0.0.4"; 
const int serverPort = 3000;
const String clientInfo = "--- Client Info: --- \nClientID: " + clientID + ", WiFi SSID: " + ssID + ", Main Server IP: " + serverIP+ ":" + serverPort;

//Process Sheduling
unsigned long previousMillis = 0;
const long clientDataInterval = 5000;

#include <DHT.h>
#define DHTPIN 4     // what pin we're connected to
#define DHTTYPE DHT22   // DHT 22  (AM2302)

String humidityValStr;
String temperatureValStr;

DHT dht(DHTPIN, DHTTYPE, 11); //// Initialize DHT sensor for normal 16mhz Arduino

WiFiClient client;

ESP8266WebServer server(80); //Server on port 80

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

  logPrintln("Configuring access point...");
  WiFi.disconnect();
  delay(1);
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
  logPrintln("HTTP server started"); //Start server
}

 
void loop() 
{
  server.handleClient();
  unsigned long currentMillis = millis();
  
  if (currentMillis - previousMillis >= clientDataInterval) {
    previousMillis = currentMillis;

    sendSensorData("hdt", "Humidity", humidityValStr);
    sendSensorData("hdt", "Temperature", temperatureValStr);
  } 
  humidityValStr =  String(dht.readHumidity()).c_str();
  temperatureValStr = String(dht.readTemperature()).c_str();
}

// --- Help Functions
void sendSensorData(String device, String sensor,String value) {
  String url = "/update";
  url += "?device=" + device;
  url += "&sensor=" + sensor;
  url += "&value=" + value;
  url += "&clientid=" + clientID;

  getRequest(url);
}
void getRequest(String url) {
  if (!client.connect(serverIP, serverPort)) {
    logPrintln("connection to: " + serverIP +":"+ String(serverPort)+" was failed");
    return;
  }
  
  client.print(String("GET ") + url + " HTTP/1.1\r\n" + "Host: " + serverIP + "\r\n" + "Connection: close\r\n\r\n");
  unsigned long timeout = millis();
  while (client.available() == 0) {
    if (millis() - timeout > 5000){ 
      logPrintln("connect to: " + serverIP +":"+ String(serverPort)+ ", url: "+url +" >>> Client Timeout !");
      client.stop(); 
      return; 
    } 
  } // Read all the lines of the reply from server and print them to Serial
  
  while (client.available()){ 
    String line = client.readStringUntil('\r'); 
  }
  logPrintln("GET Request, http://" + serverIP +":"+ String(serverPort)+ url + ", connection closed");
}

  