import API
import Connect
import Foundation
import Observation
import SwiftProtobuf
import UIKit

enum Gate: Equatable {
    case starting
    case unreachable(Problem)
    case update(store: String, message: String)
    case signedOut(notice: String?)
    case signedIn(Hgmx_Auth_V1_Account)
}

struct Problem: Error, Equatable {
    let message: String
    var reauth = false

    init(_ message: String, reauth: Bool = false) {
        self.message = message
        self.reauth = reauth
    }
}

@MainActor
@Observable
final class Session {
    let api: Api
    let label = UIDevice.current.model
    private(set) var gate = Gate.starting
    private(set) var warning: String?

    @ObservationIgnored private let tokens: TokenStore
    @ObservationIgnored private var token: String?
    @ObservationIgnored private var generation = 0
    @ObservationIgnored private var methods: Set<String> = []
    @ObservationIgnored private var checking = false

    init(api: Api, tokens: TokenStore) {
        self.api = api
        self.tokens = tokens
    }

    var account: Hgmx_Auth_V1_Account? {
        if case .signedIn(let account) = gate { account } else { nil }
    }

    func offers(_ method: String) -> Bool {
        methods.contains { $0.hasSuffix("/\(method)") }
    }

    func start() async {
        gate = .starting
        do {
            let config = try await configure()
            if config.updateRequired { return }
            methods = Set(config.methods)
            token = tokens.read()
            advance()
            guard token != nil else {
                gate = .signedOut(notice: nil)
                return
            }
            try await exchange({ await self.api.auth.getAccount(request: .init(), headers: $0) }) {
                self.gate = .signedIn($0.account)
            }
        } catch let problem as Problem {
            if gate == .starting {
                gate = .unreachable(problem)
            }
        } catch {}
    }

    func recheck() async {
        switch gate {
        case .starting, .update: return
        default: break
        }
        guard !checking else { return }
        checking = true
        defer { checking = false }
        _ = try? await configure()
    }

    func call<T>(_ rpc: (Headers) async -> ResponseMessage<T>) async throws -> T {
        try await exchange(rpc)
    }

    func signIn<T>(_ rpc: (Headers) async -> ResponseMessage<T>, grant: (T) -> Hgmx_Auth_V1_SessionGrant) async throws {
        var previous: String?
        try await exchange(rpc) {
            let g = grant($0)
            try self.tokens.write(g.token)
            previous = self.token
            self.token = g.token
            self.advance()
            self.gate = .signedIn(g.account)
        }
        if let previous {
            Task { _ = await revoke(previous) }
        }
    }

    func update<T>(_ rpc: (Headers) async -> ResponseMessage<T>, account: (T) -> Hgmx_Auth_V1_Account) async throws {
        try await exchange(rpc) { self.gate = .signedIn(account($0)) }
    }

    func deleteAccount() async throws {
        var request = Hgmx_Auth_V1_DeleteAccountRequest()
        request.confirm = true
        try await exchange({ await self.api.auth.deleteAccount(request: request, headers: $0) }) { _ in
            self.end(nil)
        }
    }

    func signOut() {
        let old = token
        end(nil)
        guard let old else { return }
        let sent = generation
        Task {
            guard await !revoke(old), sent == generation else { return }
            warning = "You're signed out on this device, but the server didn't confirm it. To end that session, sign in at \(api.origin)/auth/account and sign out \(label) under Devices."
        }
    }

    private func revoke(_ old: String) async -> Bool {
        let response = await api.auth.signOut(request: .init(), headers: api.headers(old))
        if case .failure(let error) = response.result, error.code != .unauthenticated {
            return false
        }
        return true
    }

    private func end(_ notice: String?) {
        tokens.clear()
        token = nil
        advance()
        gate = .signedOut(notice: notice)
    }

    private func advance() {
        generation += 1
        warning = nil
    }

    private func configure() async throws -> Hgmx_App_V1_GetConfigurationResponse {
        try await exchange({ await self.api.app.getConfiguration(request: .init(), headers: $0) }) {
            if $0.updateRequired {
                self.gate = .update(store: $0.storeURL, message: "Update the app to keep going.")
            }
        }
    }

    @discardableResult
    private func exchange<T>(
        _ rpc: (Headers) async -> ResponseMessage<T>,
        commit: (T) throws -> Void = { _ in }
    ) async throws -> T {
        let sent = generation
        let headers = api.headers(token)
        let response = await rpc(headers)
        guard sent == generation, !Task.isCancelled else {
            throw CancellationError()
        }
        switch response.result {
        case .success(let message):
            try commit(message)
            return message
        case .failure(let error):
            throw problem(error, signedIn: headers["Authorization"] != nil)
        }
    }

    private func problem(_ error: ConnectError, signedIn: Bool) -> Problem {
        if let update: Hgmx_App_V1_UpdateRequired = error.detail() {
            gate = .update(store: update.storeURL, message: update.message)
            return Problem(update.message)
        }
        if let again: Hgmx_Auth_V1_SignInAgain = error.detail() {
            return Problem(again.message, reauth: true)
        }
        let message = error.exception != nil
            ? "Couldn't reach \(api.origin). Check your connection."
            : error.message ?? "Something went wrong. Try again."
        if error.code == .unauthenticated && signedIn {
            end(message)
        }
        return Problem(message)
    }
}

extension ConnectError {
    func detail<M: SwiftProtobuf.Message>() -> M? {
        details
            .first { $0.type.split(separator: "/").last.map(String.init) == M.protoMessageName }
            .flatMap(\.payload)
            .flatMap { try? M(serializedBytes: $0) }
    }
}
