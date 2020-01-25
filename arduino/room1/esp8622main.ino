#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ArduinoJson.h>
#include <ESP8266HTTPClient.h>
#include <ESP8266WebServer.h>
#include <IRremoteESP8266.h>
#include <IRsend.h>

//Client Data
const String roomID = "room1";
const String clientID = String(roomID+ "_main");

//const String ssID     = "Danielle_2.4";         
//const String wifiPass = "0524713014";    
const String ssID     = "YuriIotLocal";         
const String wifiPass = "12345678";    

const String serverIP = "192.168.123.125:3000"; //Server static IP
const String clientInfo = "--- Client Info: --- \nClientID: " + clientID + ", WiFi SSID: " + ssID + ", Main Server IP: " + serverIP;

//Process Sheduling
unsigned long previousMillis = 0;
const long clientDataInterval = 10000;
const long openDoorInterval = 5000;

//Client data
StaticJsonDocument<512> clientData;

WiFiClient client;

HTTPClient http;

ESP8266WebServer server(80); //Server on port 80

const uint16_t kIrLed = 4;  // D4
int khz = 38; // 38kHz carrier frequency for the NEC protocol

IRsend irsend(kIrLed);  // Set the GPIO to be used to sending the message.

//RGBColorStrip
int stripR;
int stripG;
int stripB;
int stripFade;

//RGB Color val
String rgbColorStrip;

//IR Air Code Air Cool
String irACAirCool;

#define RED_LED 14  //D5
#define GREEN_LED 12  //D6
#define BLUE_LED 13  //D7

#include <DHT.h>
#define DHTPIN 5     // D3
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

void handleData() {
  String str = String(clientInfo + "\n\n");
  str += "RGB Color:  "+ rgbColorStrip +"\n";
  str += "IR AC AirCool:  "+ irACAirCool +"\n";
   
  server.send(200, "text/plain", str);
}

int getM(int a, int b){
  if (a >= b) {
    return -1;
  } else {
    return 1;
  }
}

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
   stripFade = fade.toInt();

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
    logPrintln("change rgbstrip to:"+ cmd);
}


void snedIRACAirCool(String cmd) {
  int pos =-1;
  uint16_t data[211] = {};
  uint16_t dataCnt =0;
  String val;
  int indx = cmd.indexOf(',') ; 

  while(indx != -1 ){
     val = cmd.substring(pos + 1, indx);
     pos = indx;
     indx = cmd.indexOf(',', indx +1);
     data[dataCnt] = val.toInt();
     dataCnt++; 
  }
  data[dataCnt] = cmd.substring(pos + 1, cmd.length()).toInt();
  dataCnt++;
  logPrintln("ir_ac_aircool get data length: " + String(dataCnt) + ", ir code: " + cmd);
  
  irsend.sendRaw(data, dataCnt, khz);  // Send a raw data capture at 38kHz.
  delay(3000);
}

#define DOOR_RELAY 16  //D2
volatile byte relayState = LOW;


void handleAction() { 
  ///action?deviceid=rgbstrip&val=LUM_1900&cmd=10,255,0,50
  String id;
  String act;
  String val;
  for (int i = 0; i < server.args(); i++) { 
    if (server.argName(i) == "deviceid") {
      id = server.arg(i);  
    } else if (server.argName(i) == "cmd") {
      act = server.arg(i);
    } else if (server.argName(i) == "val") {
      val = server.arg(i);
    }
    
  } 
  logPrintln("get an action for device id:"+id);
  if (id == "rgbstrip") {
     changeRGBStripColor(act);
     rgbColorStrip = val;
  } else if (id == "ir_ac_aircool") {
     irACAirCool = val;
     snedIRACAirCool(act);
  }  else if (id == "door") {
     if (act == "open") {
       digitalWrite(DOOR_RELAY, HIGH);
       relayState = HIGH;
       logPrintln("open main door");

     } 
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
  dht.begin();

  pinMode(GREEN_LED, OUTPUT);
  pinMode(RED_LED, OUTPUT);
  pinMode(BLUE_LED, OUTPUT);

  pinMode(DOOR_RELAY, OUTPUT);

  irsend.begin();

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
  server.on("/data", handleData);
  
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

    sendSensorData("dht_Humidity", String(dht.readHumidity()).c_str());
    sendSensorData("dht_Temperature", String(dht.readTemperature()).c_str());
    logPrintln("send dht_Humidity, dht_Temperature");
    currentMillis = millis();  
  } 
  
  
  if (relayState == HIGH && (currentMillis - previousMillis >= openDoorInterval)) {
    previousMillis = currentMillis;
    digitalWrite(DOOR_RELAY, LOW);
    logPrintln("release door relay");
    relayState == LOW;
  } 
  
}

// --- Help Functions
void sendSensorData(String device ,String value) {
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