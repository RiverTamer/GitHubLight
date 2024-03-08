// Generated by swift-openapi-generator, do not modify.
@_spi(Generated) import OpenAPIRuntime
#if os(Linux)
  @preconcurrency import struct Foundation.Data
  @preconcurrency import struct Foundation.Date
  @preconcurrency import struct Foundation.URL
#else
  import struct Foundation.Data
  import struct Foundation.Date
  import struct Foundation.URL
#endif
/// A type that performs HTTP operations defined by the OpenAPI document.
protocol APIProtocol: Sendable {
  /// Returns current lights for software implementations
  ///
  /// - Remark: HTTP `GET /lights`.
  /// - Remark: Generated from `#/paths//lights/get`.
  func get_sol_lights(_ input: Operations.get_sol_lights.Input) async throws -> Operations.get_sol_lights.Output
  /// Returns current status
  ///
  /// - Remark: HTTP `GET /status`.
  /// - Remark: Generated from `#/paths//status/get`.
  func get_sol_status(_ input: Operations.get_sol_status.Input) async throws -> Operations.get_sol_status.Output
  /// Resets the system
  ///
  /// - Remark: HTTP `GET /reset`.
  /// - Remark: Generated from `#/paths//reset/get`.
  func get_sol_reset(_ input: Operations.get_sol_reset.Input) async throws -> Operations.get_sol_reset.Output
  /// Reports the current status for monitored tuple
  ///
  /// - Remark: HTTP `POST /report`.
  /// - Remark: Generated from `#/paths//report/post`.
  func post_sol_report(_ input: Operations.post_sol_report.Input) async throws -> Operations.post_sol_report.Output
}

/// Convenience overloads for operation inputs.
extension APIProtocol {
  /// Returns current lights for software implementations
  ///
  /// - Remark: HTTP `GET /lights`.
  /// - Remark: Generated from `#/paths//lights/get`.
  func get_sol_lights(headers: Operations.get_sol_lights.Input.Headers = .init()) async throws -> Operations.get_sol_lights.Output {
    try await get_sol_lights(Operations.get_sol_lights.Input(headers: headers))
  }

  /// Returns current status
  ///
  /// - Remark: HTTP `GET /status`.
  /// - Remark: Generated from `#/paths//status/get`.
  func get_sol_status(headers: Operations.get_sol_status.Input.Headers = .init()) async throws -> Operations.get_sol_status.Output {
    try await get_sol_status(Operations.get_sol_status.Input(headers: headers))
  }

  /// Resets the system
  ///
  /// - Remark: HTTP `GET /reset`.
  /// - Remark: Generated from `#/paths//reset/get`.
  func get_sol_reset(headers: Operations.get_sol_reset.Input.Headers = .init()) async throws -> Operations.get_sol_reset.Output {
    try await get_sol_reset(Operations.get_sol_reset.Input(headers: headers))
  }

  /// Reports the current status for monitored tuple
  ///
  /// - Remark: HTTP `POST /report`.
  /// - Remark: Generated from `#/paths//report/post`.
  func post_sol_report(
    headers: Operations.post_sol_report.Input.Headers = .init(),
    body: Operations.post_sol_report.Input.Body) async throws -> Operations.post_sol_report.Output
  {
    try await post_sol_report(Operations.post_sol_report.Input(
      headers: headers,
      body: body))
  }
}

/// Server URLs defined in the OpenAPI document.
enum Servers {
  /// Local Machine
  static func server1() throws -> Foundation.URL {
    try Foundation.URL(
      validatingOpenAPIServerURL: "http://localhost:8080/v1",
      variables: [])
  }

  /// Production Server
  static func server2() throws -> Foundation.URL {
    try Foundation.URL(
      validatingOpenAPIServerURL: "https://ghls.karlkraft.com/v1",
      variables: [])
  }
}

/// Types generated from the components section of the OpenAPI document.
enum Components {
  /// Types generated from the `#/components/schemas` section of the OpenAPI document.
  enum Schemas {
    /// - Remark: Generated from `#/components/schemas/Status`.
    struct Status: Codable, Hashable, Sendable {
      /// - Remark: Generated from `#/components/schemas/Status/reports`.
      var reports: Components.Schemas.Reports?
      /// Creates a new `Status`.
      ///
      /// - Parameters:
      ///   - reports:
      init(reports: Components.Schemas.Reports? = nil) {
        self.reports = reports
      }

      enum CodingKeys: String, CodingKey {
        case reports
      }
    }

    /// - Remark: Generated from `#/components/schemas/LightColor`.
    struct LightColor: Codable, Hashable, Sendable {
      /// - Remark: Generated from `#/components/schemas/LightColor/reviewRGB`.
      var reviewRGB: Swift.String
      /// - Remark: Generated from `#/components/schemas/LightColor/mergeRGB`.
      var mergeRGB: Swift.String
      /// - Remark: Generated from `#/components/schemas/LightColor/pullRGB`.
      var pullRGB: Swift.String
      /// Creates a new `LightColor`.
      ///
      /// - Parameters:
      ///   - reviewRGB:
      ///   - mergeRGB:
      ///   - pullRGB:
      init(
        reviewRGB: Swift.String,
        mergeRGB: Swift.String,
        pullRGB: Swift.String)
      {
        self.reviewRGB = reviewRGB
        self.mergeRGB = mergeRGB
        self.pullRGB = pullRGB
      }

      enum CodingKeys: String, CodingKey {
        case reviewRGB
        case mergeRGB
        case pullRGB
      }
    }

    /// - Remark: Generated from `#/components/schemas/ReportTuple`.
    struct ReportTuple: Codable, Hashable, Sendable {
      /// - Remark: Generated from `#/components/schemas/ReportTuple/repository`.
      var repository: Swift.String
      /// - Remark: Generated from `#/components/schemas/ReportTuple/section`.
      @frozen enum sectionPayload: String, Codable, Hashable, Sendable {
        case review
        case merge
        case pull
      }

      /// - Remark: Generated from `#/components/schemas/ReportTuple/section`.
      var section: Components.Schemas.ReportTuple.sectionPayload
      /// - Remark: Generated from `#/components/schemas/ReportTuple/age`.
      var age: Swift.Int
      /// - Remark: Generated from `#/components/schemas/ReportTuple/url`.
      var url: Swift.String
      /// - Remark: Generated from `#/components/schemas/ReportTuple/notes`.
      var notes: Swift.String
      /// Creates a new `ReportTuple`.
      ///
      /// - Parameters:
      ///   - repository:
      ///   - section:
      ///   - age:
      ///   - url:
      ///   - notes:
      init(
        repository: Swift.String,
        section: Components.Schemas.ReportTuple.sectionPayload,
        age: Swift.Int,
        url: Swift.String,
        notes: Swift.String)
      {
        self.repository = repository
        self.section = section
        self.age = age
        self.url = url
        self.notes = notes
      }

      enum CodingKeys: String, CodingKey {
        case repository
        case section
        case age
        case url
        case notes
      }
    }

    /// - Remark: Generated from `#/components/schemas/Reports`.
    struct ReportsPayload: Codable, Hashable, Sendable {
      /// - Remark: Generated from `#/components/schemas/Reports/value1`.
      var value1: Components.Schemas.ReportTuple?
      /// Creates a new `ReportsPayload`.
      ///
      /// - Parameters:
      ///   - value1:
      init(value1: Components.Schemas.ReportTuple? = nil) {
        self.value1 = value1
      }

      init(from decoder: any Decoder) throws {
        var errors: [any Error] = []
        do {
          value1 = try .init(from: decoder)
        } catch {
          errors.append(error)
        }
        try Swift.DecodingError.verifyAtLeastOneSchemaIsNotNil(
          [
            value1,
          ],
          type: Self.self,
          codingPath: decoder.codingPath,
          errors: errors)
      }

      func encode(to encoder: any Encoder) throws {
        try value1?.encode(to: encoder)
      }
    }

    /// - Remark: Generated from `#/components/schemas/Reports`.
    typealias Reports = [Components.Schemas.ReportsPayload]
    /// - Remark: Generated from `#/components/schemas/ClientReport`.
    struct ClientReport: Codable, Hashable, Sendable {
      /// - Remark: Generated from `#/components/schemas/ClientReport/clientid`.
      var clientid: Swift.String
      /// - Remark: Generated from `#/components/schemas/ClientReport/reports`.
      var reports: Components.Schemas.Reports
      /// Creates a new `ClientReport`.
      ///
      /// - Parameters:
      ///   - clientid:
      ///   - reports:
      init(
        clientid: Swift.String,
        reports: Components.Schemas.Reports)
      {
        self.clientid = clientid
        self.reports = reports
      }

      enum CodingKeys: String, CodingKey {
        case clientid
        case reports
      }
    }

    /// - Remark: Generated from `#/components/schemas/Result`.
    struct Result: Codable, Hashable, Sendable {
      /// - Remark: Generated from `#/components/schemas/Result/summary`.
      var summary: Swift.String?
      /// Creates a new `Result`.
      ///
      /// - Parameters:
      ///   - summary:
      init(summary: Swift.String? = nil) {
        self.summary = summary
      }

      enum CodingKeys: String, CodingKey {
        case summary
      }
    }

    /// - Remark: Generated from `#/components/schemas/Error`.
    struct _Error: Codable, Hashable, Sendable {
      /// - Remark: Generated from `#/components/schemas/Error/summary`.
      var summary: Swift.String?
      /// Creates a new `_Error`.
      ///
      /// - Parameters:
      ///   - summary:
      init(summary: Swift.String? = nil) {
        self.summary = summary
      }

      enum CodingKeys: String, CodingKey {
        case summary
      }
    }
  }

  /// Types generated from the `#/components/parameters` section of the OpenAPI document.
  enum Parameters {}
  /// Types generated from the `#/components/requestBodies` section of the OpenAPI document.
  enum RequestBodies {}
  /// Types generated from the `#/components/responses` section of the OpenAPI document.
  enum Responses {}
  /// Types generated from the `#/components/headers` section of the OpenAPI document.
  enum Headers {}
}

/// API operations, with input and output types, generated from `#/paths` in the OpenAPI document.
enum Operations {
  /// Returns current lights for software implementations
  ///
  /// - Remark: HTTP `GET /lights`.
  /// - Remark: Generated from `#/paths//lights/get`.
  enum get_sol_lights {
    static let id: Swift.String = "get/lights"
    struct Input: Sendable, Hashable {
      /// - Remark: Generated from `#/paths/lights/GET/header`.
      struct Headers: Sendable, Hashable {
        var accept: [OpenAPIRuntime.AcceptHeaderContentType<Operations.get_sol_lights.AcceptableContentType>]
        /// Creates a new `Headers`.
        ///
        /// - Parameters:
        ///   - accept:
        init(accept: [OpenAPIRuntime.AcceptHeaderContentType<Operations.get_sol_lights.AcceptableContentType>] = .defaultValues()) {
          self.accept = accept
        }
      }

      var headers: Operations.get_sol_lights.Input.Headers
      /// Creates a new `Input`.
      ///
      /// - Parameters:
      ///   - headers:
      init(headers: Operations.get_sol_lights.Input.Headers = .init()) {
        self.headers = headers
      }
    }

    @frozen enum Output: Sendable, Hashable {
      struct Ok: Sendable, Hashable {
        /// - Remark: Generated from `#/paths/lights/GET/responses/200/content`.
        @frozen enum Body: Sendable, Hashable {
          /// - Remark: Generated from `#/paths/lights/GET/responses/200/content/application\/json`.
          case json(Components.Schemas.LightColor)
          /// The associated value of the enum case if `self` is `.json`.
          ///
          /// - Throws: An error if `self` is not `.json`.
          /// - SeeAlso: `.json`.
          var json: Components.Schemas.LightColor {
            get throws {
              switch self {
                case let .json(body):
                  body
              }
            }
          }
        }

        /// Received HTTP response body
        var body: Operations.get_sol_lights.Output.Ok.Body
        /// Creates a new `Ok`.
        ///
        /// - Parameters:
        ///   - body: Received HTTP response body
        init(body: Operations.get_sol_lights.Output.Ok.Body) {
          self.body = body
        }
      }

      /// Report
      ///
      /// - Remark: Generated from `#/paths//lights/get/responses/200`.
      ///
      /// HTTP response code: `200 ok`.
      case ok(Operations.get_sol_lights.Output.Ok)
      /// The associated value of the enum case if `self` is `.ok`.
      ///
      /// - Throws: An error if `self` is not `.ok`.
      /// - SeeAlso: `.ok`.
      var ok: Operations.get_sol_lights.Output.Ok {
        get throws {
          switch self {
            case let .ok(response):
              response
            default:
              try throwUnexpectedResponseStatus(
                expectedStatus: "ok",
                response: self)
          }
        }
      }

      struct Default: Sendable, Hashable {
        /// - Remark: Generated from `#/paths/lights/GET/responses/default/content`.
        @frozen enum Body: Sendable, Hashable {
          /// - Remark: Generated from `#/paths/lights/GET/responses/default/content/application\/json`.
          case json(Components.Schemas._Error)
          /// The associated value of the enum case if `self` is `.json`.
          ///
          /// - Throws: An error if `self` is not `.json`.
          /// - SeeAlso: `.json`.
          var json: Components.Schemas._Error {
            get throws {
              switch self {
                case let .json(body):
                  body
              }
            }
          }
        }

        /// Received HTTP response body
        var body: Operations.get_sol_lights.Output.Default.Body
        /// Creates a new `Default`.
        ///
        /// - Parameters:
        ///   - body: Received HTTP response body
        init(body: Operations.get_sol_lights.Output.Default.Body) {
          self.body = body
        }
      }

      /// Error
      ///
      /// - Remark: Generated from `#/paths//lights/get/responses/default`.
      ///
      /// HTTP response code: `default`.
      case `default`(statusCode: Swift.Int, Operations.get_sol_lights.Output.Default)
      /// The associated value of the enum case if `self` is `.`default``.
      ///
      /// - Throws: An error if `self` is not `.`default``.
      /// - SeeAlso: `.`default``.
      var `default`: Operations.get_sol_lights.Output.Default {
        get throws {
          switch self {
            case let .default(_, response):
              response
            default:
              try throwUnexpectedResponseStatus(
                expectedStatus: "default",
                response: self)
          }
        }
      }
    }

    @frozen enum AcceptableContentType: AcceptableProtocol {
      case json
      case other(Swift.String)
      init?(rawValue: Swift.String) {
        switch rawValue.lowercased() {
          case "application/json":
            self = .json
          default:
            self = .other(rawValue)
        }
      }

      var rawValue: Swift.String {
        switch self {
          case let .other(string):
            string
          case .json:
            "application/json"
        }
      }

      static var allCases: [Self] {
        [
          .json,
        ]
      }
    }
  }

  /// Returns current status
  ///
  /// - Remark: HTTP `GET /status`.
  /// - Remark: Generated from `#/paths//status/get`.
  enum get_sol_status {
    static let id: Swift.String = "get/status"
    struct Input: Sendable, Hashable {
      /// - Remark: Generated from `#/paths/status/GET/header`.
      struct Headers: Sendable, Hashable {
        var accept: [OpenAPIRuntime.AcceptHeaderContentType<Operations.get_sol_status.AcceptableContentType>]
        /// Creates a new `Headers`.
        ///
        /// - Parameters:
        ///   - accept:
        init(accept: [OpenAPIRuntime.AcceptHeaderContentType<Operations.get_sol_status.AcceptableContentType>] = .defaultValues()) {
          self.accept = accept
        }
      }

      var headers: Operations.get_sol_status.Input.Headers
      /// Creates a new `Input`.
      ///
      /// - Parameters:
      ///   - headers:
      init(headers: Operations.get_sol_status.Input.Headers = .init()) {
        self.headers = headers
      }
    }

    @frozen enum Output: Sendable, Hashable {
      struct Ok: Sendable, Hashable {
        /// - Remark: Generated from `#/paths/status/GET/responses/200/content`.
        @frozen enum Body: Sendable, Hashable {
          /// - Remark: Generated from `#/paths/status/GET/responses/200/content/application\/json`.
          case json(Components.Schemas.Status)
          /// The associated value of the enum case if `self` is `.json`.
          ///
          /// - Throws: An error if `self` is not `.json`.
          /// - SeeAlso: `.json`.
          var json: Components.Schemas.Status {
            get throws {
              switch self {
                case let .json(body):
                  body
              }
            }
          }
        }

        /// Received HTTP response body
        var body: Operations.get_sol_status.Output.Ok.Body
        /// Creates a new `Ok`.
        ///
        /// - Parameters:
        ///   - body: Received HTTP response body
        init(body: Operations.get_sol_status.Output.Ok.Body) {
          self.body = body
        }
      }

      /// Report
      ///
      /// - Remark: Generated from `#/paths//status/get/responses/200`.
      ///
      /// HTTP response code: `200 ok`.
      case ok(Operations.get_sol_status.Output.Ok)
      /// The associated value of the enum case if `self` is `.ok`.
      ///
      /// - Throws: An error if `self` is not `.ok`.
      /// - SeeAlso: `.ok`.
      var ok: Operations.get_sol_status.Output.Ok {
        get throws {
          switch self {
            case let .ok(response):
              response
            default:
              try throwUnexpectedResponseStatus(
                expectedStatus: "ok",
                response: self)
          }
        }
      }

      struct Default: Sendable, Hashable {
        /// - Remark: Generated from `#/paths/status/GET/responses/default/content`.
        @frozen enum Body: Sendable, Hashable {
          /// - Remark: Generated from `#/paths/status/GET/responses/default/content/application\/json`.
          case json(Components.Schemas._Error)
          /// The associated value of the enum case if `self` is `.json`.
          ///
          /// - Throws: An error if `self` is not `.json`.
          /// - SeeAlso: `.json`.
          var json: Components.Schemas._Error {
            get throws {
              switch self {
                case let .json(body):
                  body
              }
            }
          }
        }

        /// Received HTTP response body
        var body: Operations.get_sol_status.Output.Default.Body
        /// Creates a new `Default`.
        ///
        /// - Parameters:
        ///   - body: Received HTTP response body
        init(body: Operations.get_sol_status.Output.Default.Body) {
          self.body = body
        }
      }

      /// Accepted
      ///
      /// - Remark: Generated from `#/paths//status/get/responses/default`.
      ///
      /// HTTP response code: `default`.
      case `default`(statusCode: Swift.Int, Operations.get_sol_status.Output.Default)
      /// The associated value of the enum case if `self` is `.`default``.
      ///
      /// - Throws: An error if `self` is not `.`default``.
      /// - SeeAlso: `.`default``.
      var `default`: Operations.get_sol_status.Output.Default {
        get throws {
          switch self {
            case let .default(_, response):
              response
            default:
              try throwUnexpectedResponseStatus(
                expectedStatus: "default",
                response: self)
          }
        }
      }
    }

    @frozen enum AcceptableContentType: AcceptableProtocol {
      case json
      case other(Swift.String)
      init?(rawValue: Swift.String) {
        switch rawValue.lowercased() {
          case "application/json":
            self = .json
          default:
            self = .other(rawValue)
        }
      }

      var rawValue: Swift.String {
        switch self {
          case let .other(string):
            string
          case .json:
            "application/json"
        }
      }

      static var allCases: [Self] {
        [
          .json,
        ]
      }
    }
  }

  /// Resets the system
  ///
  /// - Remark: HTTP `GET /reset`.
  /// - Remark: Generated from `#/paths//reset/get`.
  enum get_sol_reset {
    static let id: Swift.String = "get/reset"
    struct Input: Sendable, Hashable {
      /// - Remark: Generated from `#/paths/reset/GET/header`.
      struct Headers: Sendable, Hashable {
        var accept: [OpenAPIRuntime.AcceptHeaderContentType<Operations.get_sol_reset.AcceptableContentType>]
        /// Creates a new `Headers`.
        ///
        /// - Parameters:
        ///   - accept:
        init(accept: [OpenAPIRuntime.AcceptHeaderContentType<Operations.get_sol_reset.AcceptableContentType>] = .defaultValues()) {
          self.accept = accept
        }
      }

      var headers: Operations.get_sol_reset.Input.Headers
      /// Creates a new `Input`.
      ///
      /// - Parameters:
      ///   - headers:
      init(headers: Operations.get_sol_reset.Input.Headers = .init()) {
        self.headers = headers
      }
    }

    @frozen enum Output: Sendable, Hashable {
      struct Ok: Sendable, Hashable {
        /// - Remark: Generated from `#/paths/reset/GET/responses/200/content`.
        @frozen enum Body: Sendable, Hashable {
          /// - Remark: Generated from `#/paths/reset/GET/responses/200/content/application\/json`.
          case json(Components.Schemas.Result)
          /// The associated value of the enum case if `self` is `.json`.
          ///
          /// - Throws: An error if `self` is not `.json`.
          /// - SeeAlso: `.json`.
          var json: Components.Schemas.Result {
            get throws {
              switch self {
                case let .json(body):
                  body
              }
            }
          }
        }

        /// Received HTTP response body
        var body: Operations.get_sol_reset.Output.Ok.Body
        /// Creates a new `Ok`.
        ///
        /// - Parameters:
        ///   - body: Received HTTP response body
        init(body: Operations.get_sol_reset.Output.Ok.Body) {
          self.body = body
        }
      }

      /// Report
      ///
      /// - Remark: Generated from `#/paths//reset/get/responses/200`.
      ///
      /// HTTP response code: `200 ok`.
      case ok(Operations.get_sol_reset.Output.Ok)
      /// The associated value of the enum case if `self` is `.ok`.
      ///
      /// - Throws: An error if `self` is not `.ok`.
      /// - SeeAlso: `.ok`.
      var ok: Operations.get_sol_reset.Output.Ok {
        get throws {
          switch self {
            case let .ok(response):
              response
            default:
              try throwUnexpectedResponseStatus(
                expectedStatus: "ok",
                response: self)
          }
        }
      }

      struct Default: Sendable, Hashable {
        /// - Remark: Generated from `#/paths/reset/GET/responses/default/content`.
        @frozen enum Body: Sendable, Hashable {
          /// - Remark: Generated from `#/paths/reset/GET/responses/default/content/application\/json`.
          case json(Components.Schemas._Error)
          /// The associated value of the enum case if `self` is `.json`.
          ///
          /// - Throws: An error if `self` is not `.json`.
          /// - SeeAlso: `.json`.
          var json: Components.Schemas._Error {
            get throws {
              switch self {
                case let .json(body):
                  body
              }
            }
          }
        }

        /// Received HTTP response body
        var body: Operations.get_sol_reset.Output.Default.Body
        /// Creates a new `Default`.
        ///
        /// - Parameters:
        ///   - body: Received HTTP response body
        init(body: Operations.get_sol_reset.Output.Default.Body) {
          self.body = body
        }
      }

      /// Accepted
      ///
      /// - Remark: Generated from `#/paths//reset/get/responses/default`.
      ///
      /// HTTP response code: `default`.
      case `default`(statusCode: Swift.Int, Operations.get_sol_reset.Output.Default)
      /// The associated value of the enum case if `self` is `.`default``.
      ///
      /// - Throws: An error if `self` is not `.`default``.
      /// - SeeAlso: `.`default``.
      var `default`: Operations.get_sol_reset.Output.Default {
        get throws {
          switch self {
            case let .default(_, response):
              response
            default:
              try throwUnexpectedResponseStatus(
                expectedStatus: "default",
                response: self)
          }
        }
      }
    }

    @frozen enum AcceptableContentType: AcceptableProtocol {
      case json
      case other(Swift.String)
      init?(rawValue: Swift.String) {
        switch rawValue.lowercased() {
          case "application/json":
            self = .json
          default:
            self = .other(rawValue)
        }
      }

      var rawValue: Swift.String {
        switch self {
          case let .other(string):
            string
          case .json:
            "application/json"
        }
      }

      static var allCases: [Self] {
        [
          .json,
        ]
      }
    }
  }

  /// Reports the current status for monitored tuple
  ///
  /// - Remark: HTTP `POST /report`.
  /// - Remark: Generated from `#/paths//report/post`.
  enum post_sol_report {
    static let id: Swift.String = "post/report"
    struct Input: Sendable, Hashable {
      /// - Remark: Generated from `#/paths/report/POST/header`.
      struct Headers: Sendable, Hashable {
        var accept: [OpenAPIRuntime.AcceptHeaderContentType<Operations.post_sol_report.AcceptableContentType>]
        /// Creates a new `Headers`.
        ///
        /// - Parameters:
        ///   - accept:
        init(accept: [OpenAPIRuntime.AcceptHeaderContentType<Operations.post_sol_report.AcceptableContentType>] = .defaultValues()) {
          self.accept = accept
        }
      }

      var headers: Operations.post_sol_report.Input.Headers
      /// - Remark: Generated from `#/paths/report/POST/requestBody`.
      @frozen enum Body: Sendable, Hashable {
        /// - Remark: Generated from `#/paths/report/POST/requestBody/content/application\/json`.
        case json(Components.Schemas.ClientReport)
      }

      var body: Operations.post_sol_report.Input.Body
      /// Creates a new `Input`.
      ///
      /// - Parameters:
      ///   - headers:
      ///   - body:
      init(
        headers: Operations.post_sol_report.Input.Headers = .init(),
        body: Operations.post_sol_report.Input.Body)
      {
        self.headers = headers
        self.body = body
      }
    }

    @frozen enum Output: Sendable, Hashable {
      struct Ok: Sendable, Hashable {
        /// - Remark: Generated from `#/paths/report/POST/responses/200/content`.
        @frozen enum Body: Sendable, Hashable {
          /// - Remark: Generated from `#/paths/report/POST/responses/200/content/application\/json`.
          case json(Components.Schemas.Result)
          /// The associated value of the enum case if `self` is `.json`.
          ///
          /// - Throws: An error if `self` is not `.json`.
          /// - SeeAlso: `.json`.
          var json: Components.Schemas.Result {
            get throws {
              switch self {
                case let .json(body):
                  body
              }
            }
          }
        }

        /// Received HTTP response body
        var body: Operations.post_sol_report.Output.Ok.Body
        /// Creates a new `Ok`.
        ///
        /// - Parameters:
        ///   - body: Received HTTP response body
        init(body: Operations.post_sol_report.Output.Ok.Body) {
          self.body = body
        }
      }

      /// Accepted
      ///
      /// - Remark: Generated from `#/paths//report/post/responses/200`.
      ///
      /// HTTP response code: `200 ok`.
      case ok(Operations.post_sol_report.Output.Ok)
      /// The associated value of the enum case if `self` is `.ok`.
      ///
      /// - Throws: An error if `self` is not `.ok`.
      /// - SeeAlso: `.ok`.
      var ok: Operations.post_sol_report.Output.Ok {
        get throws {
          switch self {
            case let .ok(response):
              response
            default:
              try throwUnexpectedResponseStatus(
                expectedStatus: "ok",
                response: self)
          }
        }
      }

      struct Default: Sendable, Hashable {
        /// - Remark: Generated from `#/paths/report/POST/responses/default/content`.
        @frozen enum Body: Sendable, Hashable {
          /// - Remark: Generated from `#/paths/report/POST/responses/default/content/application\/json`.
          case json(Components.Schemas._Error)
          /// The associated value of the enum case if `self` is `.json`.
          ///
          /// - Throws: An error if `self` is not `.json`.
          /// - SeeAlso: `.json`.
          var json: Components.Schemas._Error {
            get throws {
              switch self {
                case let .json(body):
                  body
              }
            }
          }
        }

        /// Received HTTP response body
        var body: Operations.post_sol_report.Output.Default.Body
        /// Creates a new `Default`.
        ///
        /// - Parameters:
        ///   - body: Received HTTP response body
        init(body: Operations.post_sol_report.Output.Default.Body) {
          self.body = body
        }
      }

      /// Accepted
      ///
      /// - Remark: Generated from `#/paths//report/post/responses/default`.
      ///
      /// HTTP response code: `default`.
      case `default`(statusCode: Swift.Int, Operations.post_sol_report.Output.Default)
      /// The associated value of the enum case if `self` is `.`default``.
      ///
      /// - Throws: An error if `self` is not `.`default``.
      /// - SeeAlso: `.`default``.
      var `default`: Operations.post_sol_report.Output.Default {
        get throws {
          switch self {
            case let .default(_, response):
              response
            default:
              try throwUnexpectedResponseStatus(
                expectedStatus: "default",
                response: self)
          }
        }
      }
    }

    @frozen enum AcceptableContentType: AcceptableProtocol {
      case json
      case other(Swift.String)
      init?(rawValue: Swift.String) {
        switch rawValue.lowercased() {
          case "application/json":
            self = .json
          default:
            self = .other(rawValue)
        }
      }

      var rawValue: Swift.String {
        switch self {
          case let .other(string):
            string
          case .json:
            "application/json"
        }
      }

      static var allCases: [Self] {
        [
          .json,
        ]
      }
    }
  }
}
