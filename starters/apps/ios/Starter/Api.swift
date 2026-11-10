import API
import Connect
import Foundation

struct Api: Sendable {
    let origin: String
    let client: ProtocolClient
    let app: Hgmx_App_V1_AppServiceClient
    let auth: Hgmx_Auth_V1_AuthServiceClient
    private let appID: String
    private let build: String

    init(origin: String, appID: String, build: String) {
        precondition(Api.exact(origin), "HGMXAPIOrigin must be an exact https origin, such as https://example.com")
        self.origin = origin
        self.appID = appID
        self.build = build
        client = ProtocolClient(httpClient: Transport(), config: ProtocolClientConfig(host: origin, codec: ProtoCodec()))
        app = Hgmx_App_V1_AppServiceClient(client: client)
        auth = Hgmx_Auth_V1_AuthServiceClient(client: client)
    }

    func headers(_ token: String?) -> Headers {
        var headers: Headers = ["HGMX-App-ID": [appID], "HGMX-App-Build": [build]]
        if let token {
            headers["Authorization"] = ["Bearer \(token)"]
        }
        return headers
    }

    static func exact(_ origin: String) -> Bool {
        guard let url = URLComponents(string: origin), let host = url.host, !host.isEmpty else {
            return false
        }
        #if DEBUG
        let local = url.scheme == "http" && host == "localhost"
        #else
        let local = false
        #endif
        return (url.scheme == "https" || local) && url.user == nil && url.password == nil
            && url.path.isEmpty && url.query == nil && url.fragment == nil
    }
}

private final class Transport: HTTPClientInterface {
    private let session: URLSession

    init() {
        let config = URLSessionConfiguration.ephemeral
        config.httpCookieAcceptPolicy = .never
        config.httpShouldSetCookies = false
        config.httpCookieStorage = nil
        config.urlCache = nil
        session = URLSession(configuration: config, delegate: NoRedirects(), delegateQueue: nil)
    }

    func unary(
        request: HTTPRequest<Data?>,
        onMetrics: @escaping @Sendable (HTTPMetrics) -> Void,
        onResponse: @escaping @Sendable (HTTPResponse) -> Void
    ) -> Cancelable {
        var urlRequest = URLRequest(url: request.url)
        urlRequest.httpMethod = request.method.rawValue
        urlRequest.httpBody = request.message
        for (name, values) in request.headers {
            urlRequest.setValue(values.joined(separator: ","), forHTTPHeaderField: name)
        }
        let task = session.dataTask(with: urlRequest) { data, response, error in
            guard let http = response as? HTTPURLResponse else {
                let code: Code = (error as? URLError)?.code == .cancelled ? .canceled : .unavailable
                onResponse(HTTPResponse(
                    code: code, headers: [:], message: data, trailers: [:],
                    error: ConnectError(code: code, message: error?.localizedDescription, exception: error),
                    tracingInfo: nil
                ))
                return
            }
            var headers = Headers()
            for (name, value) in http.allHeaderFields {
                guard let name = name as? String else { continue }
                headers[name.lowercased(), default: []] += "\(value)".split(separator: ",").map {
                    $0.trimmingCharacters(in: .whitespaces)
                }
            }
            onResponse(HTTPResponse(
                code: .fromHTTPStatus(http.statusCode), headers: headers, message: data, trailers: [:],
                error: error, tracingInfo: .init(httpStatus: http.statusCode)
            ))
        }
        task.resume()
        return Cancelable { task.cancel() }
    }

    func stream(request: HTTPRequest<Data?>, responseCallbacks: ResponseCallbacks) -> RequestCallbacks<Data> {
        responseCallbacks.receiveClose(
            .unimplemented, [:], ConnectError(code: .unimplemented, message: "This app's transport only makes unary calls.")
        )
        return RequestCallbacks(cancel: {}, sendData: { _ in }, sendClose: {})
    }
}

private final class NoRedirects: NSObject, URLSessionTaskDelegate {
    func urlSession(
        _ session: URLSession, task: URLSessionTask, willPerformHTTPRedirection response: HTTPURLResponse,
        newRequest request: URLRequest
    ) async -> URLRequest? {
        nil
    }
}
