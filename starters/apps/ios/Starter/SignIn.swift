import API
import SwiftUI

@MainActor
@Observable
final class SignInModel {
    var email = "" {
        didSet {
            if email != oldValue {
                sent = false
                code = ""
            }
        }
    }
    var code = ""
    private(set) var sent = false
    let work = Work()

    private let session: Session
    let codes: Bool
    let passkeys: Bool

    init(session: Session) {
        self.session = session
        codes = session.offers("SendCode")
        passkeys = session.offers("BeginPasskeySignin")
    }

    func sendCode() async {
        await work.run {
            var request = Hgmx_Auth_V1_SendCodeRequest()
            request.email = email.trimmingCharacters(in: .whitespaces)
            _ = try await session.call { await self.session.api.auth.sendCode(request: request, headers: $0) }
            sent = true
        }
    }

    func verify() async -> Bool {
        var done = false
        await work.run {
            var request = Hgmx_Auth_V1_VerifyCodeRequest()
            request.email = email.trimmingCharacters(in: .whitespaces)
            request.code = code.trimmingCharacters(in: .whitespaces)
            request.label = session.label
            try await session.signIn({ await self.session.api.auth.verifyCode(request: request, headers: $0) }) { $0.grant }
            done = true
        }
        return done
    }

    func passkey() async -> Bool {
        var done = false
        await work.run {
            let begin = try await session.call { await self.session.api.auth.beginPasskeySignin(request: .init(), headers: $0) }
            var request = Hgmx_Auth_V1_FinishPasskeySigninRequest()
            request.ceremony = begin.options.ceremony
            request.credentialJson = try await Credentials.passkey(begin.options.publicKeyOptionsJson)
            request.label = session.label
            try await session.signIn({ await self.session.api.auth.finishPasskeySignin(request: request, headers: $0) }) { $0.grant }
            done = true
        }
        return done
    }
}

struct SignIn: View {
    let notice: String?
    let onDone: () -> Void
    @State private var model: SignInModel

    init(session: Session, notice: String?, onDone: @escaping () -> Void = {}) {
        self.notice = notice
        self.onDone = onDone
        _model = State(initialValue: SignInModel(session: session))
    }

    var body: some View {
        @Bindable var model = model
        Form {
            Section {
                Text("Sign in").font(.largeTitle)
                if let notice {
                    Text(notice)
                }
                if model.passkeys {
                    Button("Sign in with a passkey") { act { await model.passkey() } }
                }
            }
            if model.codes {
                Section {
                    TextField("Email", text: $model.email)
                        .keyboardType(.emailAddress)
                        .textContentType(.username)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                    if model.sent {
                        TextField("Code from your email", text: $model.code)
                            .keyboardType(.numberPad)
                            .textContentType(.oneTimeCode)
                        Button("Sign in") { act { await model.verify() } }
                            .disabled(model.code.trimmingCharacters(in: .whitespaces).isEmpty)
                    } else {
                        Button("Email me a code") { Task { await model.sendCode() } }
                            .disabled(model.email.trimmingCharacters(in: .whitespaces).isEmpty)
                    }
                }
            }
            if let problem = model.work.problem {
                Section { Text(problem.message).foregroundStyle(.red) }
            }
        }
        .disabled(model.work.busy)
    }

    private func act(_ action: @escaping () async -> Bool) {
        Task {
            if await action() {
                onDone()
            }
        }
    }
}
