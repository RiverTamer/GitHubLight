//
//  ArduinoDriver.ino
//  GitHubLight
//
//  Created by Karl Kraft on 12/28/2023
//  Copyright 2023-2024 Karl Kraft. All rights reserved
//

#include "secrets.h"
#include <SharedSecrets.h>
#include <ESP8266WiFi.h>  // Include the Wi-Fi library
#include <WiFiUdp.h>

// Rearange digital pin order on baord ESP8266, and set to not use bitbanging
// https://github.com/FastLED/FastLED/wiki/ESP8266-notes
#define FASTLED_ESP8266_NODEMCU_PIN_ORDER
#define FASTLED_ALL_PINS_HARDWARE_SPI


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

  int8_t red1;
  int8_t green1;
  int8_t blue1;

  int8_t red2;
  int8_t green2;
  int8_t blue2;

  int8_t red3;
  int8_t green3;
  int8_t blue3;

  int8_t padding;
} Packet;

#define PACKET_MAGIC 0x1234

Packet packetBuffer;

//
//  LEDs
//
#define LEDS_PER_SECTION 4
#define NUM_LEDS LEDS_PER_SECTION * 3
#define LED_PIN 2
CRGB leds[NUM_LEDS];


void setup() {
  // wait for serial port
  Serial.begin(9600);
  while (!Serial) {
    ;
  }
  delay(1000);
  FastLED.addLeds<WS2812B, LED_PIN, GRB>(leds, NUM_LEDS);
  FastLED.setCorrection(TypicalLEDStrip);
  FastLED.setBrightness(100);

  Serial.println("started");

  while (true) {
    Serial.print("Connecting to: ");
    Serial.println(ssid);
    status = WiFi.begin(ssid, pass);
    if (status == WL_CONNECTED) break;
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

  Udp.begin(localPort);
  for (int x = 0; x < NUM_LEDS; x++) {
    leds[x] = CRGB(0, 0, 0);
  }
  FastLED.show();
  FastLED.show();
  delay(1000);
}

void loop() {

  // if there's data available, read a packet
  int packetSize = Udp.parsePacket();
  if (packetSize == 0) {
    delay(500);
    return;
  }
  Udp.read((char *)&packetBuffer, sizeof(Packet));
  if (packetSize != sizeof(Packet)) {
    Serial.print("Invalid packet size ");
    Serial.print(packetSize);
    Serial.print(" != ");
    Serial.println(sizeof(Packet));
    return;
  }
  if (packetBuffer.magic != PACKET_MAGIC) {
    Serial.print("Invalid value for magic ");
    Serial.print(packetBuffer.magic, HEX);
    Serial.print(" != ");
    Serial.println(PACKET_MAGIC, HEX);
    return;
  }

  for (int x = 0; x < LEDS_PER_SECTION; x++) {
    leds[x] = CRGB(packetBuffer.red1, packetBuffer.green1, packetBuffer.blue1);
    leds[x + LEDS_PER_SECTION] = CRGB(packetBuffer.red2, packetBuffer.green2, packetBuffer.blue2);
    leds[x + LEDS_PER_SECTION * 2] = CRGB(packetBuffer.red3, packetBuffer.green3, packetBuffer.blue3);
  }
  FastLED.show();
  FastLED.show();
}
