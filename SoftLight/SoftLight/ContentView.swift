//
//  ContentView.swift
//  SoftLight
//
//  Created by Karl Kraft on 2/20/24.
//  Copyright 2024 Karl Kraft. All rights reserved.
//

import OpenAPIRuntime
import OpenAPIURLSession
import SwiftUI
import Toml

struct Report: Identifiable {
  let id = UUID()
  let repository: String
  let section: String
  let age: Int
  let url: String
  let notes: String
}

// todo - reset button

struct ContentView: View {
  @State var mergeColor = Color.gray
  @State var reviewColor = Color.gray
  @State var pullColor = Color.gray
  @State var reports = [Report]()

  @State var refreshTime: Date?

  var body: some View {
    VStack {
      List {
        ForEach(reports) { report in
          if report.section == "review" {
            HStack {
              Text("Review:")
              Link(report.notes, destination: URL(string: report.url)!)
            }
          } else if report.section == "merge" {
            HStack {
              Text("Merge:")
              Link(report.notes, destination: URL(string: report.url)!)
            }
          } else if report.section == "pull" {
            Text("Pull:")
          }
        }
      }
      Group {
        Text("REVIEW").foregroundStyle(reviewColor)
        Text("MERGE").foregroundStyle(mergeColor)
        Text("PULL").foregroundStyle(pullColor)
      }.font(.system(size: 48.0)).fontWeight(.heavy)

      Spacer()

      Group {
        HStack {
          Button {
            Task { try? await refreshLights() }
            Task { try? await refreshReports() }
          } label: {
            Text("Refresh")
          }
          if let refreshTime {
            Text(refreshTime, style: .time)
          } else {
            Text("....")
          }
          Spacer()
          Button {
            Task { try? await reset() }
          } label: {
            Text("Reset")
          }
        }
      }
    }
    .padding()
    .onAppear {
      Task { try? await refreshLights() }
      Task { try? await refreshReports() }
    }.frame(minWidth: 320, maxWidth: 400)
  }

  let client: Client
  init() {
    let configPath = NSHomeDirectory() + "/.githubLightBox"
    if FileManager.default.fileExists(atPath: configPath) {
      let config = try! Toml(contentsOfFile: configPath)
      let url = Foundation.URL(string: config.string("LightServer")!)

      client = Client(serverURL: url!, transport: URLSessionTransport())
    } else {
      client = Client(serverURL: try! Servers.server1(), transport: URLSessionTransport())
    }
  }

  func reset() async throws {
    _ = try await client.get_sol_reset()
    try await refreshLights()
    try await refreshReports()
  }

  func refreshLights() async throws {
    let response = try await client.get_sol_lights()
    switch response {
      case let .ok(okResponse):
        switch okResponse.body {
          case let .json(json):
            if json.reviewRGB == "#000000" {
              reviewColor = Color(hex: "#AAAAAA")
            } else {
              reviewColor = Color(hex: json.reviewRGB)
            }
            if json.mergeRGB == "#000000" {
              mergeColor = Color(hex: "#AAAAAA")
            } else {
              mergeColor = Color(hex: json.mergeRGB)
            }
            if json.pullRGB == "#000000" {
              pullColor = Color(hex: "#AAAAAA")
            } else {
              pullColor = Color(hex: json.pullRGB)
            }
            refreshTime = Date()
        }
      case let .default(statusCode: statusCode, _):
        print("Could not get lights \(statusCode)")
    }
  }

  func refreshReports() async throws {
    let response = try await client.get_sol_status()
    switch response {
      case let .ok(okResponse):
        switch okResponse.body {
          case let .json(json):
            reports.removeAll()
            if let responseReports = json.reports {
              for r in responseReports {
                let myReport = Report(repository: r.value1!.repository, section: r.value1!.section.rawValue, age: r.value1!.age, url: r.value1!.url, notes: r.value1!.notes)
                reports.append(myReport)
              }
            }
            refreshTime = Date()
        }
      case let .default(statusCode: statusCode, _):
        print("Could not get get_sol_status \(statusCode)")
    }
  }
}

#Preview {
  ContentView()
}
