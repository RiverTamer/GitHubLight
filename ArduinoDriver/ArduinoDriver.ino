//
//  ArduinoDriver.ino
//  GitHubLight
//
//  Created by Karl Kraft on 12/28/2023
//  Copyright 2023 Karl Kraft. All rights reserved.
//

#include "secrets.h"
#include <WiFiS3.h>
#include <FastLED.h>

//
// WIFI
//
int status = WL_IDLE_STATUS;

char ssid[] = SECRET_SSID;  // your network SSID (name)
char pass[] = SECRET_PASS;  // your network password (use for WPA, or use as key for WEP)
int keyIndex = 0;           // your network key index number (needed only for WEP)

unsigned int localPort = 0xc1ad;

WiFiUDP Udp;


//
// Packet
//
typedef struct {
  int16_t magic;
  int8_t lightStart;
  int8_t lightRange;
  int8_t red;
  int8_t green;
  int8_t blue;
} Packet;

#define PACKET_MAGIC 0xda1c

Packet packetBuffer;

//
//  LEDs
//
#define NUM_LEDS 9
#define LED_PIN 5
CRGB leds[NUM_LEDS];


void setup() {
  // wait for serial port
  Serial.begin(9600);
  while (!Serial) {
    ;
  }

  // Configured LEDs
  FastLED.addLeds<WS2812B, LED_PIN, GRB>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);
  FastLED.setBrightness(25);
  // check for the WiFi module:
  if (WiFi.status() == WL_NO_MODULE) {
    Serial.println("WiFi module failed (WL_NO_MODULE)");
    while (true)
      ;
  }

  String fv = WiFi.firmwareVersion();
  if (fv < WIFI_FIRMWARE_LATEST_VERSION) {
    Serial.println("WiFi module failed (WIFI_FIRMWARE_LATEST_VERSION)");
  }

  // attempt to connect to WiFi network:
  while (status != WL_CONNECTED) {
    Serial.print("Connecting to: ");
    Serial.println(ssid);
    status = WiFi.begin(ssid, pass);
    delay(10000);
  }
  Serial.println("Connected");

  IPAddress ip = WiFi.localIP();
  Serial.print("IP: ");
  Serial.println(ip);

  // print the received signal strength:
  long rssi = WiFi.RSSI();
  Serial.print("RSSI:");
  Serial.print(rssi);
  Serial.println(" dBm");


  Serial.print("\nListening on port");
  Serial.println(localPort);
  Udp.begin(localPort);
}

void loop() {

  // if there's data available, read a packet
  int packetSize = Udp.parsePacket();
  if (packetSize != sizeof(Packet)) {
    Serial.print("Invalid packet size ");
    Serial.print(packetSize);
    Serial.print(" != ");
    Serial.println(sizeof(Packet));
    return;
  }
  Udp.read((char *)&packetBuffer, sizeof(Packet));
  if (packetBuffer.magic != PACKET_MAGIC) {
    Serial.print("Invalid value for magic");
    Serial.print(packetBuffer.magic, HEX);
    Serial.print(" != ");
    Serial.println(PACKET_MAGIC, HEX);
    return;
  }

  for (int x = 0; x < packetBuffer.lightRange; x++) {
    int idx = x + packetBuffer.lightStart;
    if (idx > NUM_LEDS) {
      continue;
    }
    leds[idx] = CRGB(packetBuffer.red, packetBuffer.green, packetBuffer.blue);
  }
  FastLED.show();
  Serial.println("Updated...");
}
