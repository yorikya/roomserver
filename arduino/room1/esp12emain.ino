#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ArduinoJson.h>
#include <ESP8266HTTPClient.h>
#include <ESP8266WebServer.h>

//Client Data
const String roomID = "room1";
const String clientID = String(roomID+ "_main");

const String ssID     = "Danielle_2.4";         
const String wifiPass = "0524713014";    

const String serverIP = "10.0.0.4:3000"; 
const String clientInfo = "--- Client Info: --- \nClientID: " + clientID + ", WiFi SSID: " + ssID + ", Main Server IP: " + serverIP;

//Process Sheduling
unsigned long previousMillis = 0;
const long clientDataInterval = 5000;

//Client data
StaticJsonDocument<512> clientData;

WiFiClient client;

HTTPClient http;

ESP8266WebServer server(80); //Server on port 80


#define RED_LED 5
#define GREEN_LED 4
#define BLUE_LED 0

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

int getM(int a, int b){
  if (a >= b) {
    return -1;
  } else {
    return 1;
  }
}

int stripR;
int stripG;
int stripB;
int stripFade;

void changeRGBStripColor(String cmd) {
   int ind1 = cmd.indexOf(',');  //finds location of first ,
   String rcolor = cmd.substring(0, ind1);   //captures first data String
   int ind2 = cmd.indexOf(',', ind1+1 );   //finds location of second ,
   String gcolor = cmd.substring(ind1+1, ind2);   //captures second data String
   int ind3 = cmd.indexOf(',', ind2+1 );
   String bcolor = cmd.substring(ind2+1, ind3);
   int ind4 = cmd.indexOf(',', ind3);
   String fade = cmd.substring(ind3+1); //captur
   int tmpR = stripR;
   int tmpG = stripG;
   int tmpB = stripB;
   int tmpFade = stripFade;
   stripR = rcolor.toInt();
   stripG = gcolor.toInt();
   stripB = bcolor.toInt();
   stripFade = fade;

    int m = getM(tmpR ,stripR);
    int iter = (stripR - tmpR ) * m;
    for (int i =0 ; i < iter; i++){
       analogWrite(RED_LED, tmpR); 
       tmpR += m;
       delay(10);
    }
    
    m = getM(tmpG ,stripG);
    iter = (stripG - tmpG ) * m;
    for (int i =0 ; i < iter; i++){
       analogWrite(GREEN_LED, tmpG); 
       tmpG += m;
       delay(10);
    }

    m = getM(tmpB ,stripB);
    iter = (stripB - tmpB ) * m;
    for (int i =0 ; i < iter; i++){
       analogWrite(BLUE_LED, tmpB); 
       tmpB += m;
       delay(10);
    }
}

void handleAction() {
  String id;
  String act;
  for (int i = 0; i < server.args(); i++) {
    if (server.argName(i) == "deviceid") {
      id = server.arg(i);  
    } else if (server.argName(i) == "cmd") {
      act = server.arg(i);
    }
    
  } 
  if (id == "rgbstrip") {
     changeRGBStripColor(act);
  }
  String message = "action id: " + id + ", cmd: " + act;
  logPrintln(message);

   server.send(200, "text/plain", message);
}
// --- End Logger

void setup() 
{
  Serial.begin(115200);
  Serial.setDebugOutput(true);

  pinMode(GREEN_LED, OUTPUT);
  pinMode(RED_LED, OUTPUT);
  pinMode(BLUE_LED, OUTPUT);
    

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
  server.on("/action", handleAction);
  
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

void loop() 
{
  server.handleClient();
  unsigned long currentMillis = millis();
  
  if (currentMillis - previousMillis >= clientDataInterval) {
    previousMillis = currentMillis;

    logPrintln("5 sec pass...");
//    analogWrite(RED_LED, 255);   // 1900 Luminate6
//    analogWrite(GREEN_LED, 147); 
//    analogWrite(BLUE_LED, 41);
//    
    } 
  
}

// --- Help Functions
void sendSensorData(String device, String sensor,String value) {
  String url = "http://" + serverIP + "/update";
  url += "?device=" + device;
  url += "&sensor=" + sensor;
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