import API
import Connect
import Foundation
import Observation
import Security
import SwiftProtobuf
import Testing
@testable import Starter

@MainActor
struct SessionTests {
    let api = Api(origin: "https://example.com", appID: "test", build: "1")
    let tokens = TokenStore(origin: "https://\(UUID().uuidString).example")

    @Test func tokenIsDeviceOnlyKeychainItem() throws {
        try tokens.write("secret")
        #expect(tokens.read() == "secret")
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: tokens.origin,
            kSecUseDataProtectionKeychain as String: true,
            kSecReturnAttributes as String: true,
        ]
        var item: CFTypeRef?
        #expect(SecItemCopyMatching(query as CFDictionary, &item) == errSecSuccess)
        let attributes = try #require(item as? [String: Any])
        #expect(attributes[kSecAttrAccessible as String] as? String == kSecAttrAccessibleWhenUnlockedThisDeviceOnly as String)
        #expect(attributes[kSecAttrSynchronizable as String] as? Bool != true)
        tokens.clear()
        #expect(tokens.read() == nil)
    }

    @Test func grantIsStored() async throws {
        let session = try await signedIn()
        #expect(tokens.read() == "first")
        #expect(session.account?.id == "a")
        tokens.clear()
    }

    @Test func lateGrantAfterSignOutIsDropped() async throws {
        let session = Session(api: api, tokens: tokens)
        let (started, start) = AsyncStream<Void>.makeStream()
        let (release, go) = AsyncStream<Void>.makeStream()
        let late = Task {
            try await session.signIn({ _ in
                start.yield()
                for await _ in release { break }
                return ResponseMessage(result: .success(Self.grant("late")))
            }) { $0.grant }
        }
        for await _ in started { break }
        session.signOut()
        go.yield()
        await #expect(throws: CancellationError.self) { try await late.value }
        #expect(tokens.read() == nil)
        #expect(session.gate == .signedOut(notice: nil))
    }

    @Test func unconfirmedSignOutSaysSo() async throws {
        let offline = Api(origin: "http://localhost:18559", appID: "test", build: "1")
        let session = Session(api: offline, tokens: tokens)
        try await session.signIn({ _ in ResponseMessage(result: .success(Self.grant("first"))) }) { $0.grant }
        session.signOut()
        #expect(session.gate == .signedOut(notice: nil))
        #expect(tokens.read() == nil)
        await withCheckedContinuation { changed in
            withObservationTracking { _ = session.gate } onChange: { changed.resume() }
        }
        guard case .signedOut(let notice?) = session.gate else {
            Issue.record("no notice after an unconfirmed sign-out: \(session.gate)")
            return
        }
        #expect(notice.contains("didn't confirm"))
        #expect(notice.contains("http://localhost:18559/auth/account"))
    }

    @Test func unauthenticatedEndsSession() async throws {
        let session = try await signedIn()
        let error = ConnectError(code: .unauthenticated, message: "Your session ended.")
        await #expect(throws: Problem("Your session ended.")) {
            try await session.call { _ -> ResponseMessage<Hgmx_Auth_V1_GetAccountResponse> in ResponseMessage(result: .failure(error)) }
        }
        #expect(session.gate == .signedOut(notice: "Your session ended."))
        #expect(tokens.read() == nil)
    }

    @Test func typedDetailsSteer() async throws {
        let session = try await signedIn()
        var again = Hgmx_Auth_V1_SignInAgain()
        again.message = "Confirm it's you."
        await #expect(throws: Problem("Confirm it's you.", reauth: true)) {
            try await session.call(Self.failing(.permissionDenied, "hgmx.auth.v1.SignInAgain", again))
        }
        #expect(session.account?.id == "a")

        var update = Hgmx_App_V1_UpdateRequired()
        update.storeURL = "https://apps.apple.com/app/id1"
        update.message = "Update to keep going."
        for type in ["hgmx.app.v1.UpdateRequired", "type.googleapis.com/hgmx.app.v1.UpdateRequired"] {
            await #expect(throws: Problem("Update to keep going.")) {
                try await session.call(Self.failing(.failedPrecondition, type, update))
            }
            #expect(session.gate == .update(store: update.storeURL, message: update.message))
        }
        tokens.clear()
    }

    @Test func cancelledCallLeavesStateAlone() async throws {
        let session = try await signedIn()
        let (started, start) = AsyncStream<Void>.makeStream()
        let (release, go) = AsyncStream<Void>.makeStream()
        let call = Task {
            try await session.call { _ in
                start.yield()
                for await _ in release { break }
                return ResponseMessage<Hgmx_Auth_V1_GetAccountResponse>(result: .failure(ConnectError(code: .unauthenticated, message: "late")))
            }
        }
        for await _ in started { break }
        call.cancel()
        go.yield()
        await #expect(throws: CancellationError.self) { try await call.value }
        #expect(session.account?.id == "a")
        #expect(tokens.read() == "first")
        tokens.clear()
    }

    @Test func unknownEnumsParse() throws {
        let device = try Hgmx_Auth_V1_Session(serializedBytes: Data([0x18, 0x63]))
        #expect(device.kind == .UNRECOGNIZED(99))
    }

    @Test func originMustBeExact() {
        #expect(Api.exact("https://example.com"))
        #expect(Api.exact("https://example.com:8443"))
        #expect(Api.exact("http://localhost:18550"))
        for origin in ["http://example.com", "https://example.com/", "https://example.com/api", "https://u@example.com", "https://example.com?x=1", "https://example.com#x", "example.com"] {
            #expect(!Api.exact(origin), "\(origin)")
        }
    }

    private func signedIn() async throws -> Session {
        let session = Session(api: api, tokens: tokens)
        try await session.signIn({ _ in ResponseMessage(result: .success(Self.grant("first"))) }) { $0.grant }
        return session
    }

    private static func grant(_ token: String) -> Hgmx_Auth_V1_VerifyCodeResponse {
        var response = Hgmx_Auth_V1_VerifyCodeResponse()
        response.grant.token = token
        response.grant.account.id = "a"
        return response
    }

    private static func failing<M: SwiftProtobuf.Message>(
        _ code: Code, _ type: String, _ detail: M
    ) -> (Headers) async -> ResponseMessage<Hgmx_Auth_V1_GetAccountResponse> {
        let payload: Data? = try? detail.serializedBytes()
        let error = ConnectError(code: code, message: "detail", details: [.init(type: type, payload: payload)])
        return { _ in ResponseMessage(result: .failure(error)) }
    }
}
