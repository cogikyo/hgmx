import AuthenticationServices
import Foundation
import UIKit

@MainActor
enum Credentials {
    static func passkey(_ options: Data) async throws -> Data {
        let o = try JSONDecoder().decode(Options.self, from: options)
        let provider = ASAuthorizationPlatformPublicKeyCredentialProvider(relyingPartyIdentifier: o.rpId ?? "")
        let request = provider.createCredentialAssertionRequest(challenge: try o.challenge.base64URL())
        request.allowedCredentials = try (o.allowCredentials ?? []).map {
            ASAuthorizationPlatformPublicKeyCredentialDescriptor(credentialID: try $0.id.base64URL())
        }
        request.userVerificationPreference = .init(rawValue: o.userVerification ?? "preferred")
        guard let assertion = try await Ceremony().perform(request) as? ASAuthorizationPlatformPublicKeyCredentialAssertion else {
            throw Problem("The passkey didn't return a sign-in.")
        }
        return try JSONSerialization.data(withJSONObject: [
            "id": assertion.credentialID.base64URL,
            "rawId": assertion.credentialID.base64URL,
            "type": "public-key",
            "authenticatorAttachment": "platform",
            "clientExtensionResults": [String: Any](),
            "response": [
                "clientDataJSON": assertion.rawClientDataJSON.base64URL,
                "authenticatorData": assertion.rawAuthenticatorData.base64URL,
                "signature": assertion.signature.base64URL,
                "userHandle": assertion.userID.base64URL,
            ],
        ])
    }

    static func createPasskey(_ options: Data) async throws -> Data {
        let o = try JSONDecoder().decode(Options.self, from: options)
        guard let rp = o.rp?.id, let user = o.user else {
            throw Problem("The server sent incomplete passkey options.")
        }
        let provider = ASAuthorizationPlatformPublicKeyCredentialProvider(relyingPartyIdentifier: rp)
        let request = provider.createCredentialRegistrationRequest(
            challenge: try o.challenge.base64URL(), name: user.name, userID: try user.id.base64URL()
        )
        request.userVerificationPreference = .init(rawValue: o.authenticatorSelection?.userVerification ?? "preferred")
        guard let registration = try await Ceremony().perform(request) as? ASAuthorizationPlatformPublicKeyCredentialRegistration,
              let attestation = registration.rawAttestationObject
        else {
            throw Problem("The passkey wasn't created.")
        }
        return try JSONSerialization.data(withJSONObject: [
            "id": registration.credentialID.base64URL,
            "rawId": registration.credentialID.base64URL,
            "type": "public-key",
            "authenticatorAttachment": "platform",
            "clientExtensionResults": [String: Any](),
            "response": [
                "clientDataJSON": registration.rawClientDataJSON.base64URL,
                "attestationObject": attestation.base64URL,
                "transports": ["internal", "hybrid"],
            ],
        ])
    }
}

private struct Options: Decodable {
    struct Entity: Decodable {
        let id: String?
    }

    struct User: Decodable {
        let id: String
        let name: String
    }

    struct Credential: Decodable {
        let id: String
    }

    struct Selection: Decodable {
        let userVerification: String?
    }

    let challenge: String
    let rpId: String?
    let rp: Entity?
    let user: User?
    let allowCredentials: [Credential]?
    let userVerification: String?
    let authenticatorSelection: Selection?
}

@MainActor
private final class Ceremony: NSObject, ASAuthorizationControllerDelegate, ASAuthorizationControllerPresentationContextProviding {
    private var continuation: CheckedContinuation<ASAuthorizationCredential, any Error>?

    func perform(_ request: ASAuthorizationRequest) async throws -> ASAuthorizationCredential {
        let controller = ASAuthorizationController(authorizationRequests: [request])
        controller.delegate = self
        controller.presentationContextProvider = self
        return try await withCheckedThrowingContinuation {
            continuation = $0
            controller.performRequests()
        }
    }

    func authorizationController(controller: ASAuthorizationController, didCompleteWithAuthorization authorization: ASAuthorization) {
        continuation?.resume(returning: authorization.credential)
        continuation = nil
    }

    func authorizationController(controller: ASAuthorizationController, didCompleteWithError error: any Error) {
        if (error as? ASAuthorizationError)?.code == .canceled {
            continuation?.resume(throwing: CancellationError())
        } else {
            continuation?.resume(throwing: Problem(error.localizedDescription))
        }
        continuation = nil
    }

    func presentationAnchor(for controller: ASAuthorizationController) -> ASPresentationAnchor {
        UIApplication.shared.connectedScenes
            .compactMap { ($0 as? UIWindowScene)?.keyWindow }
            .first ?? ASPresentationAnchor()
    }
}

private extension String {
    func base64URL() throws -> Data {
        var s = replacingOccurrences(of: "-", with: "+").replacingOccurrences(of: "_", with: "/")
        s += String(repeating: "=", count: (4 - s.count % 4) % 4)
        guard let data = Data(base64Encoded: s) else {
            throw Problem("The server sent malformed passkey options.")
        }
        return data
    }
}

private extension Data {
    var base64URL: String {
        base64EncodedString()
            .replacingOccurrences(of: "+", with: "-")
            .replacingOccurrences(of: "/", with: "_")
            .replacingOccurrences(of: "=", with: "")
    }
}
