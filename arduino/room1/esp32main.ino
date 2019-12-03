#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>

//Client Data
const String roomID = "room1";
const String clientID = String(roomID + "_main");

const String ssID = "Danielle_2.4";
const String wifiPass = "0524713014";

const String serverIP = "10.0.0.4";
const int serverPort = 3000;
const String clientInfo = "--- Client Info: --- \nClientID: " + clientID + ", WiFi SSID: " + ssID + ", Main Server IP: " + serverIP;

//Process Sheduling
unsigned long previousMillis = 0;
const long clientDataInterval = 5000;

String humidityValStr = "46.9";
String temperatureValStr = "21.0";

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


//==============================================================
//     This rutine is exicuted when you open its IP in browser
//==============================================================
void handleUpdate() {
   String device = server.arg("device"); //this lets you access a query param (http://x.x.x.x/action1?value=1)
   String sensor = server.arg("sensor");
   String value = server.arg("value");
   String clientid = server.arg("clientid");
   Serial.print("device: " + device);
   Serial.print(" sensor: " + sensor);
   Serial.print(" value: " + value);
   Serial.println(" clientid: " + clientid);
   
   server.send(200, "text/plain", "hello from esp8266!");
}

void setup() {
 
  Serial.begin(115200);

  logPrintln("Configuring access point...");
  WiFi.disconnect();
  delay(1);
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssID, wifiPass);
  logPrintln("Connecting to " + ssID);

  int i = 0;
  while (WiFi.localIP().toString() == "0.0.0.0") { 
    delay(1000);
    logPrintln("connecting to SSID:" + ssID + " attempt: " + String(++i) + ", WiFi Status: " + String(WiFi.status()));
  }
  logPrintln("Connection established!\nIP address: " + WiFi.localIP().toString());  
 
  server.on("/update", handleUpdate);      //Which routine to handle at root location
  server.on("/logs", handleLogs);     
 
  server.begin();     
  logPrintln("HTTP server started"); //Start server
}
   
void loop(){
  server.handleClient();    
  unsigned long currentMillis = millis();

  if (currentMillis - previousMillis >= clientDataInterval) {
    previousMillis = currentMillis;

    sendSensorData("hdt", "Humidity", humidityValStr);
    sendSensorData("hdt", "Temperature", temperatureValStr);
        
    logPrintln("send client dht sensor data, move senso");
  }  
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
  int stok = client.connect(serverIP, serverPort);
  if (!stok) {
    
    logPrintln("connection failed status: " + String(stok) +" server: " + serverIP +":"+ String(serverPort)+" was failed");
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