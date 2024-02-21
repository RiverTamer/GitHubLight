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

struct Report: Identifiable {
  let id = UUID()
  let owner: String
  let repository: String
  let section: String
  let age: Int
  let reference: String
}

// todo -refresh
// todo - fetch data and display
// todo - reset button

struct ContentView: View {
  @State var mergeColor = Color.gray
  @State var reviewColor = Color.gray
  @State var pullColor = Color.gray
  @State var reports = [Report]()

  var body: some View {
    HStack {
      VStack {
        Group {
          Text("REVIEW").foregroundStyle(reviewColor)
          Text("MERGE").foregroundStyle(mergeColor)
          Text("PULL").foregroundStyle(pullColor)
        }.font(.system(size: 48.0)).fontWeight(.heavy)

        Spacer()
        Button {
          Task { try? await refreshLights() }
          Task { try? await refreshReports() }
        } label: {
          Text("Refresh")
        }
      }
      .padding()
      Spacer()
      Table(reports) {
        TableColumn("Owner", value: \.owner)
        TableColumn("Repositry", value: \.repository)
        TableColumn("Section", value: \.section)
        TableColumn("Age") { _ in
          Text("1")
        }
        TableColumn("Ref", value: \.reference)
      }
    }.onAppear {
      Task { try? await refreshLights() }
      Task { try? await refreshReports() }
    }
  }

  let client: Client
  init() {
    client = Client(serverURL: try! Servers.server1(), transport: URLSessionTransport())
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
                let myReport = Report(owner: r.value1!.owner, repository: r.value1!.repository, section: r.value1!.section.rawValue, age: r.value1!.age, reference: r.value1!.reference)
                reports.append(myReport)
              }
            }
        }
      case let .default(statusCode: statusCode, _):
        print("Could not get get_sol_status \(statusCode)")
    }
  }
}

#Preview {
  ContentView()
}
